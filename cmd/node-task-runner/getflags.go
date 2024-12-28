package main

import (
	"flag"
)

type flags struct {
	Help    bool
	Version bool
}

func getFlags() flags {
	help := flag.Bool("help", false, "Show help text")
	version := flag.Bool("version", false, "Show version")

	flag.Parse()
	return flags{
		Help:    *help,
		Version: *version,
	}
}
