package main

import (
	"flag"
)

type flags struct {
	Help    bool
	Version bool
	Cwd     string
	Debug   bool
}

func getFlags() flags {
	help := flag.Bool("help", false, "Show help text")
	version := flag.Bool("version", false, "Show version")
	cwd := flag.String("cwd", "", "Specify current working directory")
	debug := flag.Bool("debug", false, "Turn debug logging on")
	flag.Parse()
	return flags{
		Help:    *help,
		Version: *version,
		Cwd:     *cwd,
		Debug:   *debug,
	}
}
