package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
)

type NpmPackage struct {
	InstalledPath   string
	PackageName     string
	PackageVersion  string
	PackageJsonPath string
	//
}

func (p NpmPackage) SpecName() string {
	return fmt.Sprintf(
		"%s@%s",
		p.PackageName,
		p.PackageVersion,
	)
}

type PackageJson struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func GatheringNodeModules(targetDir string) []NpmPackage {

	packages := []NpmPackage{}

	ReadDirRecursively(targetDir, func(dir string) {
		packagePath := path.Join(dir, "package.json")
		fileBytes, err := ioutil.ReadFile(packagePath)
		if err != nil {
			// package.json 이 없으면 아무것도안함
			return
		}
		packageJson := PackageJson{}
		err = json.Unmarshal(fileBytes, &packageJson)
		if err != nil {
			panic(err)
		}
		if packageJson.Name != "" && packageJson.Version != "" {

			packages = append(packages, NpmPackage{
				InstalledPath:   dir,
				PackageName:     packageJson.Name,
				PackageVersion:  packageJson.Version,
				PackageJsonPath: packagePath,
			})
		}

	})

	return packages
}
