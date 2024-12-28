package main

import "flag"

type flags struct {
	Help    bool
	Version bool
}

func getFlags() flags {
	help := flag.Bool("help", false, "Print help")
	version := flag.Bool("version", false, "Print version")
	flag.Parse()
	return flags{
		Help:    *help,
		Version: *version,
	}
}
