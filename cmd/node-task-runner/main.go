package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
    // Define flags
    name := flag.String("name", "World", "a name to say hello to")
    flag.Parse()

    // Print the greeting
    fmt.Printf("Hello, %s!\n", *name)

    // Handle additional commands
    if len(os.Args) > 1 {
        switch os.Args[1] {
        case "version":
            fmt.Println("Node Task Runner CLI v1.0.0")
        case "fzf":
            runFzf()
        default:
            fmt.Printf("Unknown command: %s\n", os.Args[1])
        }
    }
}

func runFzf() {
    // Example fzf usage
    opts := fzf.DefaultOptions()
    opts.Prompt = "Select an option> "
    opts.Items = []string{"Option 1", "Option 2", "Option 3"}
    selected, err := fzf.Run(opts)
    if err != nil {
        fmt.Println("Error running fzf:", err)
        return
    }
    fmt.Println("You selected:", selected)
}
