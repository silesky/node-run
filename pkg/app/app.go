package app

import (
	"context"
	"encoding/json"
	"fmt"
	"node-task-runner/pkg/fuzzsearch"
	"node-task-runner/pkg/logger"
	"os"
	"path/filepath"
	"strings"
)

// expandPath expands/normalize the ~ to the user's home directory
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

// hasWorkspacesArray checks if the package.json contains a workspaces array
func hasWorkspacesArray(packageJSONPath string) bool {
	file, err := os.ReadFile(packageJSONPath)
	if err != nil {
		return false
	}

	var packageJSON map[string]interface{}
	err = json.Unmarshal(file, &packageJSON)
	if err != nil {
		return false
	}

	_, ok := packageJSON["workspaces"]
	return ok
}

// parseWorkspacesArray parses the workspaces array in the root package.json and finds all package.json files
func parseWorkspacesArray(monorepoRoot string) ([]string, error) {
	packageJSONPath := filepath.Join(monorepoRoot, "package.json")
	file, err := os.ReadFile(packageJSONPath)
	if err != nil {
		return nil, err
	}

	var packageJSON map[string]interface{}
	err = json.Unmarshal(file, &packageJSON)
	if err != nil {
		return nil, err
	}

	workspaces, ok := packageJSON["workspaces"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("workspaces array not found in package.json")
	}

	var packageJSONPaths []string
	for _, ws := range workspaces {
		wsPattern, ok := ws.(string)
		if !ok {
			continue
		}
		matches, err := filepath.Glob(filepath.Join(monorepoRoot, wsPattern, "package.json"))
		if err != nil {
			continue
		}
		packageJSONPaths = append(packageJSONPaths, matches...)
	}

	return packageJSONPaths, nil
}

// findAllPackageJSONs finds all package.json files starting from the given directory
func findAllPackageJSONs(startDir string) ([]string, error) {
	startDir, err := expandPath(startDir)
	if err != nil {
		return nil, err
	}

	monorepoRoot, err := findMonorepoRoot(startDir)
	if err != nil {
		return nil, err
	}

	// get the workspaces array in the root package.json
	packageJSONPaths, err := parseWorkspacesArray(monorepoRoot)
	if err != nil {
		return nil, err
	}

	if len(packageJSONPaths) == 0 {
		return nil, fmt.Errorf("no package.json found")
	}

	return packageJSONPaths, nil
}

// findMonorepoRoot finds the root of the monorepo by looking for a package.json with a workspaces array
func findMonorepoRoot(startDir string) (string, error) {
	for {
		packageJSONPath := filepath.Join(startDir, "package.json")
		if _, err := os.Stat(packageJSONPath); err == nil {
			// check if the package.json contains a workspaces array
			if hasWorkspacesArray(packageJSONPath) {
				return startDir, nil
			}
		}
		parentDir := filepath.Dir(startDir)
		if parentDir == startDir {
			return "", fmt.Errorf("could not find the root of the monorepo")
		}
		startDir = parentDir
	}
}

// GetPackages gets all packages in the current monorepo, regardless of the cwd
func GetPackages(cwd string) ([]string, error) {
	cwd, err := expandPath(cwd)
	if err != nil {
		return nil, err
	}

	// find the root of the monorepo
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
