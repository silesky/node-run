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
  ntk [command]

Available Commands:
  help        Show this help message
  version     Print the version number

Flags:
  --help show this help message

Examples:
  ntk --name=YourName
  ntk version

Use "ntk [command] --help" for more information about a command.
`
	fmt.Println(helpText)
}

const VERSION = "1.0.0"

func main() {
	flag.Parse()

	help := *flag.Bool("help", false, "Print help")
	version := *flag.Bool("version", false, "Print version")

	if version {
		fmt.Printf("Version: %s\n", VERSION)
	}
	if help {
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
