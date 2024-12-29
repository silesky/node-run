package main

import (
	"flag"
)

type flags struct {
	Help    bool
	Version bool
	Cwd     string
}

func getFlags() flags {
	help := flag.Bool("help", false, "Show help text")
	version := flag.Bool("version", false, "Show version")
	cwd := flag.String("cwd", ".", "Specify current working directory")
	flag.Parse()
	return flags{
		Help:    *help,
		Version: *version,
		Cwd:     *cwd,
	}
}
