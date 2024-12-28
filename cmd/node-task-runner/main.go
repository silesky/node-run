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
		default:
			fmt.Printf("Unrecognized command: %s\n", os.Args[1])
			os.Exit(1)
		}
	} else {
		app.Run()
	}
}
