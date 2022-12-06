package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/Soulsbane/follow/internal/fileutils"
	"github.com/alexflint/go-arg"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jwalton/gchalk"
)

type Link struct {
	name   string
	path   string
	exists bool
}

var ErrorColor = gchalk.WithBold().Red

func handleFileName(fileNames []string, ugly bool) {
	results := []Link{}

	for _, fileName := range fileNames {
		info, err := os.Lstat(fileName)

		if err != nil {
			fmt.Printf("%s could not be found!\n", fileName)
		} else {
			if fileutils.IsLink(info.Mode()) {
				linkPath, pathExists := fileutils.GetLinkPath(fileName)
				currentLink := Link{name: fileName, path: linkPath, exists: pathExists}
				results = append(results, currentLink)
			} else {
				fmt.Printf("%s is not a link!\n", fileName)
			}
		}
	}

	outputResults(results, ugly)
}

func listLinks(ugly bool, showHidden bool) {
	links := fileutils.GetListOfLinks(showHidden)
	results := []Link{}

	for _, link := range links {
		linkPath, exists := fileutils.GetLinkPath(link.Name())
		currentLink := Link{name: link.Name(), path: linkPath, exists: exists}

		results = append(results, currentLink)
	}

	if len(results) > 0 {
		outputResults(results, ugly)
	} else {
		fmt.Println("No links found!")
	}
}

func toYesNo(value bool) string {
	if value {
		return "Yes"
	}

	return "No"
}

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
				fmt.Fprintf(writer, "%s\t => %s %s\n", link.name, link.path, ErrorColor("(BROKEN)"))
			}
		} else {
			if link.exists {
				outputTable.AppendRow(table.Row{link.name, link.path, toYesNo(link.exists)})
			} else {
				outputTable.AppendRow(table.Row{link.name, ErrorColor(link.path), ErrorColor(toYesNo(link.exists))})
			}
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
		handleFileName(args.FileName, args.Ugly)
	} else {
		listLinks(args.Ugly, args.Hidden)
	}
}
