package main

import (
	"flag"
	"fmt"
	"node-task-runner/pkg/app"
	"os"
)

var VERSION = "dev" // This is overridden during build time using ldflags during release

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

	// eg: "nrun --debug test" -> "test" = initialInut
	var initialInput string
	if len(flag.Args()) > 0 {
		initialInput = flag.Args()[0]
	}

	settings, err := app.NewSettings(app.WithCwd(flags.Cwd), app.WithDebug(flags.Debug), app.WithInitialInput(initialInput))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating settings: %v\n", err)
		return
	}
	app.Run(settings)
}

func usage() {
	fmt.Print("Usage: nrun [options] [initial search input]" + "\n")
	fmt.Print("Example:" + "\n" + "- nrun" + "\n" + "- nrun test" + "\n")
	fmt.Print("Options:" + "\n")
	flag.VisitAll(func(f *flag.Flag) {
		// override to show e.g (--flag) instead of -flag
		fmt.Fprintf(os.Stderr, "  --%s\n\t%s\n", f.Name, f.Usage)
	})
}
