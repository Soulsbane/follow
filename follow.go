package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"text/tabwriter"

	"github.com/alexflint/go-arg"
	"github.com/fatih/color"
)

func getFileName(info os.FileInfo, colorize bool) string {
	if colorize == true {
		if isFileExecutable(info) {
			return color.HiRedString(info.Name())
		}

		return color.WhiteString(info.Name())
	}

	return info.Name()
}

func isFileHidden(info os.FileInfo) bool {
	if runtime.GOOS != "windows" {
		if info.Name()[0:1] == "." {
			return true
		}

		return false
	}
	// FIXME: Can't seem to find documentation for properly handling this on windows.
	/*else {
		//if runtime.GOOS == "windows" {
		pointer, err := syscall.UTF16PtrFromString(info.Name())

		if err != nil {
			return false
		}

		attributes, err := syscall.GetFileAttributes(pointer)

		if err != nil {
			return false
		}

		return attributes&syscall.FILE_ATTRIBUTE_HIDDEN != 0
	}*/
	return false
}

func getLinkPath(info os.FileInfo) string {
	mode := info.Mode()
	link := mode & os.ModeSymlink

	if link != 0 {
		linkPath, _ := filepath.EvalSymlinks(info.Name())
		return linkPath
	}

	return ""
}

func isFileExecutable(info os.FileInfo) bool {
	mode := info.Mode()
	exec := mode & 0111

	if exec != 0 {
		return true
	}

	return false
}

func listFiles(ugly bool, showHidden bool) {
	var filteredFiles []os.FileInfo
	files, err := ioutil.ReadDir(".")

	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if !f.IsDir() {
			if isFileHidden(f) {
				if showHidden {
					filteredFiles = append(filteredFiles, f)
				}
			} else {
				filteredFiles = append(filteredFiles, f)
			}
		}
	}

	outputResults(filteredFiles, ugly)
}

func outputResults(files []os.FileInfo, ugly bool) {
	writer := tabwriter.NewWriter(os.Stdout, 1, 4, 1, ' ', 0)

	for _, f := range files {
		if ugly {
			fmt.Fprintf(writer, "%s\n", getFileName(f, false))
		} else {
			fmt.Fprintf(writer, "%s\n", getFileName(f, true))
		}
	}

	writer.Flush()
}

// If passed a file name it will show the linked path. If no arguments it will scan directory for links and display their paths.
func main() {
	var args struct {
		Ugly   bool `arg:"-u" default:"false" help:"Remove colorized output. Yes it's ugly."`
		Hidden bool `arg:"-i" default:"false" help:"Show hidden files."`
	}

	arg.MustParse(&args)
	listFiles(args.Ugly, args.Hidden)
}
