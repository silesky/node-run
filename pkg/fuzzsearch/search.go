package fuzzsearch

import (
	"encoding/json"
	"fmt"
	"log"
	"node-task-runner/pkg/logger"
	"os"
)

type PkgJson struct {
	Name    string            `json:"name"`
	Scripts map[string]string `json:"scripts"`
}

type Command struct {
	PackageName string
	Name        string
	Command     string
}

// Get commands from the scripts key and return them
func GetCommandsFromPaths(cwd string) (*Command, error) {
	packages, err := GetPackages(cwd)
	if err != nil {
		return nil, fmt.Errorf("could not get packages: %v", err)
	}
	logger.Debugf("Found packages: %#v\n", packages)
	commands := parseCommandsFromFiles(packages)
	selectedCommand, err := displayCommandSelector(commands)
	if err != nil {
		return nil, err
	}
	logger.Debugf("Selected Command: %+v", selectedCommand)
	return selectedCommand, nil
}

// Display command selector menu (returns user input)
func displayCommandSelector(commands []Command) (*Command, error) {
	return DisplayCommandSelector(commands)
}

// Read and parse JSON files from the provided paths
func parseCommandsFromFiles(pkgJsonPaths []string) []Command {
	var packages []PkgJson
	for _, path := range pkgJsonPaths {
		packageJSON, err := parsePkgJsonFile(path)
		if err != nil {
			log.Printf("Error parsing JSON file %s: %v", path, err)
			continue
		}
		packages = append(packages, *packageJSON)
	}
	var commands []Command
	for _, pkg := range packages {
		commands = append(commands, parseCommands(pkg)...)
	}
	return commands
}

// Parse commands list from a package json
func parseCommands(packageJson PkgJson) []Command {
	var commands []Command
	for key, value := range packageJson.Scripts {
		commands = append(commands, Command{
			PackageName: packageJson.Name,
			Name:        key,
			Command:     value,
		})
	}
	return commands
}

// Parse package json
func parsePkgJsonFile(path string) (*PkgJson, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var packageJSON PkgJson
	err = json.Unmarshal(file, &packageJSON)
	if err != nil {
		return nil, err
	}

	return &packageJSON, nil
}
