package fileutils

import (
	"io/ioutil"
	"log"
	"os"
)

func GetListOfFiles(includeHidden bool) []os.FileInfo {
	var files []os.FileInfo
	dirList, err := ioutil.ReadDir(".")

	if err != nil {
		log.Fatal(err)
	}

	for _, f := range dirList {
		if !f.IsDir() {
			if IsFileHidden(f) {
				if includeHidden {
					files = append(files, f)
				}
			} else {
				files = append(files, f)
			}
		}
	}

	return files
}
