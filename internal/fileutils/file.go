package fileutils

import (
	"errors"
	"os"
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

func IsLink(info os.FileInfo) bool {
	return info.Mode()&os.ModeSymlink != 0
}

// GetLinkPath returns the path of the link and a boolean indicating if the link destination path exists
func GetLinkPath(info os.FileInfo) (string, bool) {
	realPath, err := os.Readlink(info.Name())

	if err != nil {
		return "", FileOrPathExists(realPath)
	}

	return realPath, FileOrPathExists(realPath)
}
