package main

import "flag"

type Flags struct {
	Help    bool
	Version bool
}

func getFlags() Flags {
	help := flag.Bool("help", false, "Print help")
	version := flag.Bool("version", false, "Print version")
	flag.Parse()
	return Flags{
		Help:    *help,
		Version: *version,
	}
}
