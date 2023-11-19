package lib

import (
	"fmt"
	"os"
	"os/exec"
)

func DecompressTarballPackage(tarball PackageTarball, destination string) error {
	os.Mkdir(destination, os.ModePerm)

	cmd := exec.Command(
		"tar",
		"-xzf",
		tarball.WorkingDirRelativePath,
		"-C",
		destination,
		"--strip-components", // 압축내용의 루트 디렉토리 (package) 를 제외한다
		"1",
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		WriteAppend(
			"./go.log",
			fmt.Sprintf(":failed to decompress tar \nerr:%s\n%s\n\n", err.Error(), string(out)),
		)
		return err
	} else {
		WriteAppend(
			"./go.log",
			fmt.Sprintf(":decompressed\n%s\n\n", string(out)),
		)
	}

	return nil
}

func compressTarballPackage(tarball PackageTarball, targetDirectory string) error {

	cmd := exec.Command(
		"tar",
		"-czf",
		tarball.WorkingDirRelativePath,
		targetDirectory,
	)

	out, err := cmd.CombinedOutput()

	if err != nil {
		WriteAppend(
			"./go.log",
			fmt.Sprintf(":failed to compress tar \nerr:%s\n%s\n\n", err.Error(), string(out)),
		)
		return err
	} else {
		WriteAppend(
			"./go.log",
			fmt.Sprintf(":compressed\n%s\n\n", string(out)),
		)
	}

	return nil
}
