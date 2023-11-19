package lib

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

type PackageTarball struct {
	Filename               string
	WorkingDirRelativePath string
	TarballName            string
}
type PublishResult struct {
	Tarball      PackageTarball
	Success      bool
	Message      string
	Error        error
	NpmErrorCode string
}

var MatchErrorCode, _ = regexp.Compile("^npm ERR! code ([^\\s\\n\\t]+)")

func PublishPackedTarball(tarball PackageTarball, retry int) PublishResult {
	ret := PublishResult{
		Tarball: tarball,
	}

	cmd := exec.Command(
		"npm",
		"publish",
		fmt.Sprintf("./%s", tarball.WorkingDirRelativePath),
	)

	out, err := cmd.CombinedOutput()
	ret.Message = string(out)

	//fmt.Println(ret.Message)
	if err != nil {
		outLines := strings.Split(ret.Message, "\n")
		for _, outLine := range outLines {
			// npm ERR! code
			trimmedLine := strings.TrimSpace(outLine)
			if MatchErrorCode.Match([]byte(trimmedLine)) {
				matches := MatchErrorCode.FindSubmatch([]byte(trimmedLine))
				if len(matches) > 1 {
					ret.NpmErrorCode = string(matches[1])
				}
			}
		}

		// npm.registry 권한 문제
		if ret.NpmErrorCode == "ENEEDAUTH" {
			err := PublishErrorHandleENEEDAUTH(tarball)

			if err == nil {
				return PublishPackedTarball(tarball, retry+1)
			}
		}

		//fmt.Printf(err.Error())
		ret.Success = false
		ret.Error = err
		return ret
	}

	ret.Success = true
	return ret
}
