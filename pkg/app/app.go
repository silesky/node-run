package app

import (
	"context"
	"fmt"
	"node-task-runner/pkg/fuzzsearch"
	"node-task-runner/pkg/logger"
	"os"
)

func Run(ctx context.Context) {
	settings := FromSettingsContext(ctx)
	logger.SetDebug(settings.Debug)
	cwd := settings.Cwd
	if cwd == "" {
		var err error
		cwd, err = os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
			return
		}
	}
	logger.Debugf("Current directory: %s", cwd)

	selectedCommand, err := fuzzsearch.GetCommandsFromPaths(cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting commands: %v\n", err)
		return
	} else {
		fmt.Printf("Executing: %v", selectedCommand)
	}
}
