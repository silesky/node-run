package main

import (
	"flag"
	"fmt"
	"os"
)

type flags struct {
	Help    bool
	Version bool
}

func getFlags() flags {
	help := flag.Bool("help", false, "Show help text")
	version := flag.Bool("version", false, "Show version")

	// custom help function
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\nUsage of ntk \n")
		flag.VisitAll(func(f *flag.Flag) {
			// override to show e.g (--flag) instead of -flag
			fmt.Fprintf(os.Stderr, "  --%s\n\t%s\n", f.Name, f.Usage)
		})

	}
	flag.Parse()
	return flags{
		Help:    *help,
		Version: *version,
	}
}
