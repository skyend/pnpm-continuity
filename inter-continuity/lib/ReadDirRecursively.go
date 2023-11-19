package lib

import (
	"io/ioutil"
	"path"
)

func ReadDirRecursively(dir string, visitor func(dir string)) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		relativePath := path.Join(dir, file.Name())
		if file.IsDir() {
			ReadDirRecursively(relativePath, visitor)
		}
		visitor(relativePath)
	}
}
