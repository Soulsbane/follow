package main

type ProgramArgs struct {
	FileName []string `arg:"positional" help:"The name of the link(s) to follow"`
	Ugly     bool     `arg:"-u" default:"false" help:"Remove colorized output. Yes it's ugly."`
	Hidden   bool     `arg:"-i" default:"false" help:"Show hidden files."`
}

func (args ProgramArgs) Description() string {
	return "Follows all links to there destination and prints a nice readout of the paths."
}
