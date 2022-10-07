package fileutils

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

func IsFileHidden(info os.FileInfo) bool {
	if runtime.GOOS != "windows" {
		return info.Name()[0:1] == "."
	}

	return false
}

func FileOrPathExists(fileName string) bool {
	if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		return true
	}
}

func GetLinkPath(info os.FileInfo) string {
	mode := info.Mode()
	link := mode & os.ModeSymlink

	if link != 0 {
		linkPath, err := filepath.EvalSymlinks(info.Name())

		if err != nil {
			return ""
		}

		return linkPath
	}

	return ""
}
