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
  --name      Specify a name to greet

Examples:
  ntk --name=YourName
  ntk version

Use "ntk [command] --help" for more information about a command.
`
	fmt.Println(helpText)
}

func main() {
	flag.Parse()
	if len(os.Args) > 1 {

		switch os.Args[1] {
		case "help":
			printHelp()
		case "version":
			fmt.Printf("Node Task Runner CLI: %v\n", "1.0.0")
		default:
			fmt.Println("Unrecognized argument")
		}
	} else {
		app.Run()
	}
}
