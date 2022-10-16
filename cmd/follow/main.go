package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/Soulsbane/follow/internal/fileutils"
	"github.com/alexflint/go-arg"
	"github.com/jedib0t/go-pretty/v6/table"
)

type Link struct {
	name   string
	path   string
	exists bool
}

func handleFileName(fileName string, ugly bool) {
	info, err := os.Lstat(fileName)

	if err != nil {
		fmt.Println("File doesn't exist!")
	} else {
		if fileutils.IsLink(info) {
			linkPath, pathExists := fileutils.GetLinkPath(info)
			results := []Link{}
			currentLink := Link{name: fileName, path: linkPath, exists: pathExists}

			outputResults(append(results, currentLink), ugly)
		} else {
			fmt.Printf("%s is not a link!\n", fileName)
		}
	}
}

func listLinks(ugly bool, showHidden bool) {
	links := fileutils.GetListOfLinks(showHidden)
	results := []Link{}

	for _, link := range links {
		linkPath, exists := fileutils.GetLinkPath(link)
		currentLink := Link{name: link.Name(), path: linkPath, exists: exists}

		results = append(results, currentLink)
	}

	if len(results) > 0 {
		outputResults(results, ugly)
	} else {
		fmt.Println("No links found!")
	}
}

// TODO: Add color to broken links
func outputResults(results []Link, ugly bool) {
	writer := tabwriter.NewWriter(os.Stdout, 1, 4, 1, ' ', 0)
	outputTable := table.NewWriter()

	outputTable.SetOutputMirror(os.Stdout)
	outputTable.AppendHeader(table.Row{"Name", "Destination", "Exists"})

	for _, link := range results {
		if ugly {
			if link.exists {
				fmt.Fprintf(writer, "%s\t => %s\n", link.name, link.path)
			} else {
				fmt.Fprintf(writer, "%s\t => %s (BROKEN)\n", link.name, link.path)
			}
		} else {
			outputTable.AppendRow(table.Row{link.name, link.path, link.exists})
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

	if len(args.FileName) > 0 {
		if args.FileName[0] != "" {
			handleFileName(args.FileName[0], args.Ugly)
		}
	} else {
		listLinks(args.Ugly, args.Hidden)
	}
}
