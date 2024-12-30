package app

import (
	"context"
	"fmt"
	"node-task-runner/pkg/fuzzsearch"
	"node-task-runner/pkg/logger"
	"os"
	"path/filepath"
	"strings"
)

// expands/normalize the ~ to the user's home directory
func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(homeDir, path[1:])
	}
	return path, nil
}

func findAllPackageJSONs(startDir string) ([]string, error) {
	startDir, err := expandPath(startDir)
	if err != nil {
		return nil, err
	}

	var packageJSONPaths []string
	currentDir := startDir

	// keep iterating until you find all of the package.jsons
	for {
		packageJSONPath := filepath.Join(currentDir, "package.json")
		if _, err := os.Stat(packageJSONPath); err == nil {
			packageJSONPaths = append(packageJSONPaths, packageJSONPath)
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			// Reached the root directory -- if its a parent of itself, we're at the root
			break
		}

		currentDir = parentDir
	}

	if len(packageJSONPaths) == 0 {
		return nil, fmt.Errorf("no package.json found")
	}

	return packageJSONPaths, nil
}

func Run(ctx context.Context) {
	settings := FromSettingsContext(ctx)
	logger.SetDebug(settings.Debug)
	currentDirectory := settings.Cwd
	if currentDirectory == "" {
		var err error
		currentDirectory, err = os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		}
	}
	logger.Debugf("Current directory: %s", currentDirectory)
	packages, err := findAllPackageJSONs(currentDirectory)

	// make an anonymous struct
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error looking for packages: %v\n", err)
	} else {
		fmt.Printf("found packages: %v", packages)
	}
	// TODO: add packages
	fuzzsearch.GetCommandsFromPaths(packages)
}
