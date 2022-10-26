package fileutils

import (
	"log"
	"os"
)

func GetListOfLinks(includeHidden bool) []os.DirEntry {
	var files []os.DirEntry
	dirList, err := os.ReadDir(".")

	if err != nil {
		log.Fatal(err)
	}

	for _, f := range dirList {
		if IsLink(f.Type()) {
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
