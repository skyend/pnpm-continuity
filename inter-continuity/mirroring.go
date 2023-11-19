package main

import (
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"os"
	"os/exec"
	"path"
	"pnpm-inter-continuity/inter-continuity/lib"
	"sync"
)

var resultsGatheringMutex = sync.Mutex{}

// PnpmInstallCmd
// Pnpm install command for multiplatform
// --force for include optional packages
var PnpmInstallCmd = exec.Command("pnpm", "i", "--force")

func AppendResultThreadSafe(results *[]lib.PackResult, packResult lib.PackResult) {
	resultsGatheringMutex.Lock()
	*results = append(*results, packResult)
	resultsGatheringMutex.Unlock()
}

func main() {
	readyToMirror := make(chan bool)
	packageTaskQueue := make(chan lib.NpmPackage)

	claimingCount := 0
	results := []lib.PackResult{}

	fmt.Println("Pnpm Installing....")
	PnpmInstallCmd.Run()
	fmt.Println("Done Pnpm Install")

	// Destination Directory 생성
	os.Mkdir(lib.PackDestination, os.ModePerm)

	go func() {
		packages := lib.GatheringNodeModules("./node_modules")
		distinctList := map[string]lib.NpmPackage{}
		distinctCount := 0
		for _, packageInfo := range packages {
			fullName := fmt.Sprintf("%s-%s", packageInfo.PackageName, packageInfo.PackageVersion)
			if _, ok := distinctList[fullName]; !ok {
				distinctList[fullName] = packageInfo
				distinctCount++
			}
		}
		claimingCount = distinctCount

		fmt.Println("Detected packages:", len(packages))
		fmt.Println("Distinct packages:", distinctCount)

		readyToMirror <- true

		for _, packageInfo := range distinctList {
			packageTaskQueue <- packageInfo
			//close(packageTaskQueue) // 한개만 테스트
			//return
		}
		close(packageTaskQueue)
	}()
	<-readyToMirror

	bar := pb.StartNew(claimingCount)

	workerCount := 10

	workerWait := sync.WaitGroup{}
	workerWait.Add(workerCount)

	for i := 0; i < workerCount; i++ {
		go func() {
			defer workerWait.Done()

			for {
				packageInfo, more := <-packageTaskQueue
				if !more {
					return
				}
				packageSpecName := packageInfo.SpecName()

				packResult := lib.PackResult{
					Err:                 nil,
					Success:             false,
					InputPackageName:    packageInfo.PackageName,
					InputPackageVersion: packageInfo.PackageVersion,
				}
				cmd := exec.Command("npm", "pack", packageSpecName)
				out, err := cmd.CombinedOutput()
				if err != nil {
					packResult.Err = err
					AppendResultThreadSafe(&results, packResult)
					bar.Increment()
					continue
				}

				output, err := lib.ParseNpmPackCmdOut(string(out))
				if err != nil {
					packResult.Err = err
					packResult.Success = false

					AppendResultThreadSafe(&results, packResult)
					bar.Increment()
					continue
				}

				packResult.Success = true
				packResult.Filename = output.Filename
				packResult.ShaSum = output.ShaSum
				packResult.PackageSize = output.PackageSize
				packResult.UnpackedSize = output.UnpackedSize
				packResult.Integrity = output.Integrity
				packResult.TotalFiles = output.TotalFiles

				AppendResultThreadSafe(&results, packResult)

				os.Rename(packResult.Filename, path.Join(lib.PackDestination, packResult.Filename))
				bar.Increment()
			}
		}()
	}

	workerWait.Wait()

	bar.Finish()

	agg := lib.AggregatePackResults(results)

	resultDoc := ""
	resultDoc += agg.Print()
	if agg.Failed > 0 {
		resultDoc += lib.PrintFailedPacks(results)
	}

	os.WriteFile("mirror-log.txt", []byte(resultDoc), os.ModePerm)
}
