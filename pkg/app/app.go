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
	currentDirectory := settings.Cwd
	if currentDirectory == "" {
		var err error
		currentDirectory, err = os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
			return
		}
	}
	logger.Debugf("Current directory: %s", currentDirectory)
	packages, err := GetPackages(currentDirectory)

	// make an anonymous struct
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error looking for packages: %v\n", err)
		return
	}
	logger.Debugf("Found packages: %#v\n", packages)
	fuzzsearch.GetCommandsFromPaths(packages)
}
