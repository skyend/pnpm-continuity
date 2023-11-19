package lib

import (
	"fmt"
	"os"
	"os/exec"
)

func PublishErrorHandleENEEDAUTH(tarball PackageTarball) error {
	err := DecompressTarballPackage(tarball, tarball.TarballName)
	if err != nil {
		panic(err)
	}

	// edit package.json
	editPackageJsonCmd := exec.Command(
		"node",
		"./scripts/packageJsonEditor.js",
		"-t",
		tarball.TarballName,
	)
	out, err := editPackageJsonCmd.CombinedOutput()
	if err != nil {
		WriteAppend(
			"go.log",
			fmt.Sprintf("Failed to edit package.json %s\n%s", tarball.TarballName, string(out)),
		)
		return err
	}

	defer os.RemoveAll(tarball.TarballName)
	return compressTarballPackage(tarball, tarball.TarballName)
}
