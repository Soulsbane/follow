package main

type ProgramArgs struct {
	FileName []string `arg:"positional"`
	Ugly     bool     `arg:"-u" default:"false" help:"Remove colorized output. Yes it's ugly."`
	Hidden   bool     `arg:"-i" default:"false" help:"Show hidden files."`
}
