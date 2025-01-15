package main

import (
	"flag"
	"fmt"
	"node-task-runner/pkg/app"
	"os"
)

const VERSION = "1.0.0"

func main() {
	flags := getFlags()
	flag.Usage = usage
	if flags.Version {
		fmt.Printf("Version: %s\n", VERSION)
		return
	} else if flags.Help {
		flag.Usage()
		return
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "help":
			flag.Usage()
			return
		}
	}
	settings, err := app.NewSettings(app.WithCwd(flags.Cwd), app.WithDebug(flags.Debug))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating settings: %v\n", err)
		return
	}
	app.Run(settings)
}

func usage() {
	fmt.Fprintf(os.Stderr, "\nUsage of nt \n")
	flag.VisitAll(func(f *flag.Flag) {
		// override to show e.g (--flag) instead of -flag
		fmt.Fprintf(os.Stderr, "  --%s\n\t%s\n", f.Name, f.Usage)
	})
}
