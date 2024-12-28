package fuzzsearch

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/ktr0731/go-fuzzyfinder"
)

type PackageJSON struct {
	Name    string            `json:"name"`
	Scripts map[string]string `json:"scripts"`
}

type Command struct {
	PackageName string
	Name        string
	Command     string
}

func parseCommands(packageJSON PackageJSON) []Command {
	var commands []Command
	for key, value := range packageJSON.Scripts {
		commands = append(commands, Command{
			PackageName: packageJSON.Name,
			Name:        key,
			Command:     value,
		})
	}
	return commands
}

func Search(paths []string) {
	var packages []PackageJSON

	// Read and parse JSON files from the provided paths
	for _, path := range paths {
		packageJSON, err := parseJSONFile(path)
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

	idx, err := fuzzyfinder.Find(
		commands,
		func(idx int) string {
			return fmt.Sprintf("[%s] %s - %s", commands[idx].PackageName, commands[idx].Name, commands[idx].Command)
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("selected: %v\n", idx)
}

func parseJSONFile(path string) (*PackageJSON, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var packageJSON PackageJSON
	err = json.Unmarshal(file, &packageJSON)
	if err != nil {
		return nil, err
	}

	return &packageJSON, nil
}
