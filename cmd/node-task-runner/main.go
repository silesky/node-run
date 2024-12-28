package main

import (
	"flag"
	"fmt"
	"node-task-runner/pkg/app"
	"os"
)

func printHelp() {
	helpText := `
Node Task Runner CLI

Usage:
  ntk

Flags:
  --help - show this help message
  --version - show version

Examples:
  ntk 
`
	fmt.Println(helpText)
}

const VERSION = "1.0.0"

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

func main() {
	flags := getFlags()
	if flags.Version {
		fmt.Printf("Version: %s\n", VERSION)
		return
	} else if flags.Help {
		printHelp()
		return
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "help":
			printHelp()
		default:
			fmt.Printf("Unrecognized command: %s\n", os.Args[1])
			os.Exit(1)
		}
	} else {
		app.Run()
	}
}
