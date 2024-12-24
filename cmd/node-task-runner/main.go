package nodetaskrunner

import (
	"flag"
	"fmt"
	"node-task-runner/pkg/app"
	"os"
)

func main() {
    name := flag.String("name", "World", "a name to say hello to")
    flag.Parse()

    fmt.Printf("Hello, %s!\n", *name)

    // Handle additional commands
    if len(os.Args) > 1 {
        switch os.Args[1] {
        case "version":
            fmt.Println("Node Task Runner CLI v1.0.0")
        default:
            app.Run()
        }
    }
}


   
