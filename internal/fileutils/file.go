package fileutils

import (
	"errors"
	"os"
	"runtime"
)

func IsFileHidden(info os.DirEntry) bool {
	if runtime.GOOS != "windows" {
		return info.Name()[0:1] == "."
	}

	return false
}

func FileOrPathExists(fileName string) bool {
	if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}

func IsLink(info os.FileMode) bool {
	return info&os.ModeSymlink != 0
}

// GetLinkPath returns the path of the link and a boolean indicating if the link destination path exists
func GetLinkPath(name string) (string, bool) {
	realPath, err := os.Readlink(name)

	if err != nil {
		return "", FileOrPathExists(realPath)
	}

	return realPath, FileOrPathExists(realPath)
}
