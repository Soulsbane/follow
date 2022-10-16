package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/Soulsbane/follow/internal/fileutils"
	"github.com/alexflint/go-arg"
	"github.com/jedib0t/go-pretty/v6/table"
)

func handleFileName(fileName string, ugly bool) {
	info, err := os.Lstat(fileName)

	if err != nil {
		fmt.Println("File doesn't exist!")
	} else {
		linkPath, pathExists := fileutils.GetLinkPath(info)

		if ugly {
			if pathExists {
				fmt.Printf("%s => %s\n", fileName, linkPath)
			} else {
				fmt.Printf("%s => %s (broken)\n", fileName, linkPath)
			}
		} else {
			outputTable := table.NewWriter()

			outputTable.SetOutputMirror(os.Stdout)
			outputTable.AppendHeader(table.Row{"Name", "Destination", "Exists"})
			outputTable.AppendRow(table.Row{fileName, linkPath, pathExists})

			outputTable.SetStyle(table.StyleRounded)
			outputTable.Render()
		}
	}
}

func listLinks(ugly bool, showHidden bool) {
	linkResults := make(map[string]string)
	links := fileutils.GetListOfLinks(showHidden)

	for _, link := range links {
		// TODO: Handle broken links
		linkPath, _ := fileutils.GetLinkPath(link)
		linkResults[link.Name()] = linkPath
	}

	if len(linkResults) > 0 {
		outputResults(linkResults, ugly)
	} else {
		fmt.Println("No links found!")
	}
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
		outputTable.SetStyle(table.StyleRounded)
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
