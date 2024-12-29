package fuzzsearch

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/ktr0731/go-fuzzyfinder"
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
func GetCommandsFromPaths(pkgJsonPaths []string) (*Command, error) {
	commands := parseCommandsFromFiles(pkgJsonPaths)
	selectedIdx, err := displayCommandSelector(commands)
	if err != nil {
		return nil, err
	}
	fmt.Printf("selected: %v\n", selectedIdx)
	return &commands[selectedIdx], nil
}

func displayCommandSelector(commands []Command) (int, error) {
	idx, err := fuzzyfinder.Find(
		commands,
		func(idx int) string {
			return fmt.Sprintf("[%s] %s - %s", commands[idx].PackageName, commands[idx].Name, commands[idx].Command)
		},
	)
	return idx, err
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
