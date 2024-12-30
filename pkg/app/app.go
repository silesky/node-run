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

// findAllPackageJSONs finds all package.json files starting from the given directory
func findAllPackageJSONs(startDir string) ([]string, error) {
	startDir, err := expandPath(startDir)
	if err != nil {
		return nil, err
	}

	var packageJSONPaths []string
	err = filepath.Walk(startDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == "node_modules" {
			return filepath.SkipDir
		}
		if info.Name() == "package.json" {
			packageJSONPaths = append(packageJSONPaths, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if len(packageJSONPaths) == 0 {
		return nil, fmt.Errorf("no package.json found")
	}

	return packageJSONPaths, nil
}

// GetPackages gets all packages in the current monorepo, regardless of the cwd
func GetPackages(cwd string) ([]string, error) {
	cwd, err := expandPath(cwd)
	if err != nil {
		return nil, err
	}

	// Traverse up to find the root of the monorepo
	for {
		if _, err := os.Stat(filepath.Join(cwd, "package.json")); err == nil {
			break
		}
		parentDir := filepath.Dir(cwd)
		if parentDir == cwd {
			return nil, fmt.Errorf("could not find the root of the monorepo")
		}
		cwd = parentDir
	}

	return findAllPackageJSONs(cwd)
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
	logger.Debugf("Found packages: %v\n", packages)
	fuzzsearch.GetCommandsFromPaths(packages)
}
