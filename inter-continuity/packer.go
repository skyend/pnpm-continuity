package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"pnpm-inter-continuity/inter-continuity/lib"
	"regexp"
	"strings"
)

type JsonMap interface {
}

type PackageJson struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type PackInfo struct {
	PackageJson

	RelativeDirPath string
	PackageJsonPath string
	Err             error
}

var orgReplacer = regexp.MustCompile(`^(@)`)

func RemoveOrgMark(name string) string {
	return string(orgReplacer.ReplaceAll([]byte(name), []byte("")))
}

func main() {
	fmt.Println("Hello")

	packTargetBuffer := make(chan PackInfo, 10000)
	failedPacks := []PackInfo{}
	done := make(chan bool)
	totalCount := 0
	packedCount := 0
	//var wg sync.WaitGroup

	go func() {
		packages := lib.GatheringNodeModules("./node_modules")
		fmt.Println("packages:", len(packages))
		done <- true
	}()
	<-done
	fmt.Println("Done travel")

	packDir := "./packs2"
	os.Mkdir(packDir, os.ModePerm)

	workerCount := 40
	packCounter := make(chan int, workerCount)
	waitGo := make(chan bool, workerCount)
	for i := 0; i < workerCount; i++ {
		go func() {
			for {
				packInfo, more := <-packTargetBuffer
				if !more {
					waitGo <- true
					return
				}
				//fmt.Println("Packing...", packInfo.Name, packInfo.Version)
				outputName := strings.ReplaceAll(RemoveOrgMark(packInfo.Name), "/", "-") + "-" + packInfo.Version + ".tgz"

				if _, err := os.Stat(path.Join(packDir, outputName)); !errors.Is(err, os.ErrNotExist) {
					// path/to/whatever does not exist!
					//fmt.Println("//Skip - Already packed ", outputName, packInfo.Name, packInfo.Version)
					packCounter <- 1
					continue
				}

				//cmd := exec.Command("npm", "pack", packInfo.RelativeDirPath, "--pack-destination", packDir)
				cmd := exec.Command("tar", "-czf", outputName, "-C", packInfo.RelativeDirPath, ".")
				if err := cmd.Run(); err != nil {
					//_, err := cmd.CombinedOutput()
					fmt.Println("Pack fail", packInfo.Name, packInfo.Version, err.Error())
					//wg.Add(1)
					//packInfo.Err = err
					//failedPacks = append(failedPacks, packInfo)
					//wg.Done()
				} else {

					//fmt.Println("Pack done", packInfo.Name, packInfo.Version)
					packCounter <- 1
				}

				cmd = exec.Command("mv", outputName, path.Join(packDir, outputName))
				cmd.Run()
			}
		}()
	}

	go func() {
		// discounter
		for {
			packedCount += <-packCounter
			//fmt.Printf("Progress %d/%d \n", packedCount, totalCount)
		}
	}()

	for i := 0; i < workerCount; i++ {
		<-waitGo
	}
	fmt.Println("Processes")

	fmt.Println("Failed")
	for _, failedPackage := range failedPacks {
		fmt.Printf("%s / %s\n", failedPackage.Name, failedPackage.Version)
		fmt.Printf("- %s\n", failedPackage.Err.Error())
	}

	fmt.Printf("\n")
	fmt.Printf("\n")
	fmt.Printf("\n")
	fmt.Printf("Total:%d, Done: %d \n", totalCount, packedCount)
}
