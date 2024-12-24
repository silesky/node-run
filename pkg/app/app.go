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

func Run() {
	currentDirectory, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
	}
	packages, err := findAllPackageJSONs(currentDirectory)
	if err != nil {
		fmt.Println("Error looking for packages", err)
	}
	fmt.Println(packages)
	source := []string{"apple", "banana", "cherry"}
	substring := "ban"
	item, found := fuzzsearch.Search(source, substring)
	if found {
		fmt.Printf("Match found: %s\n", item)
	} else {
		fmt.Println("No match found")
	}
}
