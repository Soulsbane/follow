package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/Soulsbane/follow/internal/fileutils"
	"github.com/alexflint/go-arg"
	"github.com/jedib0t/go-pretty/v6/table"
)

func listLinks(ugly bool, showHidden bool) {
	filteredLinks := filterValidLinks(fileutils.GetListOfFiles(showHidden))

	if len(filteredLinks) > 0 {
		outputResults(filteredLinks, ugly)
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
			fmt.Printf("%s\n", fileutils.GetLinkPath(info))
		} else {
			fmt.Printf("%s\n", fileutils.GetLinkPath(info))
		}
	}
}

func filterValidLinks(files []os.FileInfo) map[string]string {
	filteredFiles := make(map[string]string)

	for _, f := range files {
		linkPath := fileutils.GetLinkPath(f)

		if len(linkPath) > 0 {
			filteredFiles[f.Name()] = linkPath
		}
	}

	return filteredFiles
}

func outputResults(files map[string]string, ugly bool) {
	writer := tabwriter.NewWriter(os.Stdout, 1, 4, 1, ' ', 0)
	outputTable := table.NewWriter()

	outputTable.SetOutputMirror(os.Stdout)
	outputTable.AppendHeader(table.Row{"Name", "Destination"})

	for fileName, linkPath := range files {
		if ugly {
			fmt.Fprintf(writer, "%s\t => %s\n", fileName, linkPath)
		} else {
			outputTable.AppendRow(table.Row{fileName, linkPath})
		}
	}

	if ugly {
		writer.Flush()
	} else {
		outputTable.SetStyle(table.StyleLight)
		outputTable.Render()
	}
}

// If passed a file name it will show the linked path. If no arguments it will scan directory for links and display their paths.
func main() {
	var args ProgramArgs

	arg.MustParse(&args)

	if args.FileName != "" {
		handleFileName(args.FileName, args.Ugly)
	} else {
		listLinks(args.Ugly, args.Hidden)
	}
}
