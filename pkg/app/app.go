package app

import (
	"fmt"
	"node-task-runner/pkg/fuzzsearch"
	"os"
	"path/filepath"
)

func findAllPackageJSONs(startDir string) ([]string, error) {

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

// TODO: be able to pass the path to the package.json as an argument in flags.
// TODO: be able to specifiy the cwd with cwd
func Run() {
	currentDirectory, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
	}
	packages, err := findAllPackageJSONs(currentDirectory)

	// make an anonymous struct
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error looking for packages: %v\n", err)
	} else {
		fmt.Printf("found packages: %v", packages)
	}
	// TODO: add packages
	fuzzsearch.Search([]string{"/Users/seth.silesky/projects/node-raw-socket/package.json"})
}
