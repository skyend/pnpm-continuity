package lib

import (
	"strconv"
	"strings"
)

type PackStdOutput struct {
	Filename     string
	ShaSum       string
	PackageSize  string
	UnpackedSize string
	Integrity    string
	TotalFiles   int
}

func ParseNpmPackCmdOut(result string) (PackStdOutput, error) {
	packStdOutput := PackStdOutput{}
	outTextLines := strings.Split(result, "\n")
	for _, outLine := range outTextLines {
		nameAndValuePair := strings.Split(outLine, ":")
		if len(nameAndValuePair) == 2 {
			name := strings.TrimSpace(nameAndValuePair[0])
			value := strings.TrimSpace(nameAndValuePair[1])

			switch name {
			case "npm notice filename":
				packStdOutput.Filename = value
				break
			case "npm notice package size":
				packStdOutput.PackageSize = value
				break
			case "npm notice unpacked size":
				packStdOutput.UnpackedSize = value
				break
			case "npm notice shasum":
				packStdOutput.ShaSum = value
				break
			case "npm notice integrity":
				packStdOutput.Integrity = value
				break
			case "npm notice total files":
				totalFiles, _ := strconv.ParseInt(value, 10, 32)
				packStdOutput.TotalFiles = int(totalFiles)
				break
			}
		}

	}

	return packStdOutput, nil
}
