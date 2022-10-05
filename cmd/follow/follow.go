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
	"github.com/brettski/go-termtables"
)

func isFileHidden(info os.FileInfo) bool {
	if runtime.GOOS != "windows" {
		return info.Name()[0:1] == "."
	}

	return false
}

func fileOrPathExists(fileName string) bool {
	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			return false
		}

		return false
	}

	return true
}

func getLinkPath(info os.FileInfo, colorize bool) string {
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

func listFiles(ugly bool, showHidden bool) {
	var files []os.FileInfo
	dirList, err := ioutil.ReadDir(".")

	if err != nil {
		log.Fatal(err)
	}

	for _, f := range dirList {
		if !f.IsDir() {
			if isFileHidden(f) {
				if showHidden {
					files = append(files, f)
				}
			} else {
				files = append(files, f)
			}
		}
	}

	filteredFiles := filterValidLinks(files)

	if len(filteredFiles) > 0 {
		outputResults(filteredFiles, ugly)
	} else {
		fmt.Println("No links found!")
	}
}

func handleFileName(fileName string, ugly bool) {
	info, err := os.Lstat(fileName)

	if err != nil {
		fmt.Println("File doesn't exist!")
	} else {
		if ugly {
			fmt.Printf("%s\n", getLinkPath(info, false))
		} else {
			fmt.Printf("%s\n", getLinkPath(info, true))
		}
	}
}

func filterValidLinks(files []os.FileInfo) map[string]string {
	filteredFiles := make(map[string]string)

	for _, f := range files {
		linkPath := getLinkPath(f, true)

		if len(linkPath) > 0 {
			filteredFiles[f.Name()] = linkPath
		}
	}

	return filteredFiles
}

func outputResults(files map[string]string, ugly bool) {
	writer := tabwriter.NewWriter(os.Stdout, 1, 4, 1, ' ', 0)
	table := termtables.CreateTable()
	table.AddHeaders("Name", "Destination")

	for fileName, linkPath := range files {
		if ugly {
			fmt.Fprintf(writer, "%s\t => %s\n", fileName, linkPath)
		} else {
			table.AddRow(fileName, linkPath)
		}
	}

	if ugly {
		writer.Flush()
	} else {
		fmt.Println(table.Render())
	}
}

// If passed a file name it will show the linked path. If no arguments it will scan directory for links and display their paths.
func main() {
	var args ProgramArgs

	arg.MustParse(&args)

	if args.FileName != "" {
		handleFileName(args.FileName, args.Ugly)
	} else {
		listFiles(args.Ugly, args.Hidden)
	}
}
