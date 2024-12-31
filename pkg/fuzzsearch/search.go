package fuzzsearch

import (
	"encoding/json"
	"fmt"
	"log"
	"node-task-runner/pkg/logger"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Package struct {
	Path   string
	IsRoot bool
	Json   PkgJson
}

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
		return nil, fmt.Errorf("could not get packages at %v: %v", cwd, err)
	}
	logger.Debugf("Found packages: %#v\n", packages)
	commands := parseAllCommands(packages)
	selectedCommand, err := displayCommandSelector(commands)
	if err != nil {
		return nil, err
	}
	logger.Debugf("Selected Command: %+v", selectedCommand)
	return selectedCommand, nil
}

// Display command selector menu (returns user input)
func displayCommandSelector(commands []Command) (*Command, error) {
	selectedCmd, err := DisplayCommandSelector(commands)
	return selectedCmd, err
}

// Read and parse JSON files from the provided paths
func parseAllCommands(Packages []Package) []Command {
	var packages []PkgJson
	for _, p := range Packages {
		path := p.Path
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
func findAllPackageJSONs(startDir string) ([]Package, error) {

	monorepoRoot, err := findMonorepoRoot(startDir)
	if err != nil {
		return nil, err
	}

	// get the workspaces array in the root package.json
	packageJSONPaths, err := parseWorkspacesArray(monorepoRoot)

	// append the root package.json
	packageJSONPaths = append(packageJSONPaths, filepath.Join(monorepoRoot, "package.json"))

	if err != nil {
		return nil, err
	}

	if len(packageJSONPaths) == 0 {
		return nil, fmt.Errorf("no package.json found")
	}

	var packages []Package
	for i, path := range packageJSONPaths {
		isRoot := i == 0
		pkg, err := CreatePackageFromPath(path, isRoot)
		if err != nil {
			return nil, err
		}
		packages = append(packages, *pkg)
	}
	return packages, nil
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
			return "", fmt.Errorf("could not find the root of the monorepo: %s", startDir)
		}
		startDir = parentDir
	}
}

// type Package struct {
// 	Name    string
// 	Scripts map[string]string
// }

func CreatePackageFromPath(path string, isRoot bool) (*Package, error) {
	pkgJson, err := parsePkgJsonFile(path)
	if err != nil {
		return nil, err
	}

	return &Package{
		Path:   path,
		IsRoot: isRoot,
		Json:   *pkgJson,
	}, nil
}

// GetPackages gets all packages in the current monorepo, regardless of the cwd
func GetPackages(cwd string) ([]Package, error) {
	println("cwd: ", cwd)

	// save the original cwd
	ogCwd := cwd
	// find the root of the monorepo
	for {
		if _, err := os.Stat(filepath.Join(cwd, "package.json")); err == nil {
			break
		}
		parentDir := filepath.Dir(cwd)
		if parentDir == cwd {
			return nil, fmt.Errorf("could not find the root of the monorepo: %s", ogCwd)
		}
		cwd = parentDir
	}

	if pkgs, err := findAllPackageJSONs(cwd); err != nil {
		return nil, err
	} else {
		return sortByClosedToCwd(pkgs, ogCwd), nil
	}
}

func normalizePath(path string) string {
	return strings.TrimSuffix(path, string(filepath.Separator))
}

// _isSubdirectory checks if subdir is a subdirectory of parent.
func isSubdirectory(parent, subdir string) bool {
	parent = normalizePath(parent)
	subdir = normalizePath(subdir)
	res := strings.Contains(subdir, parent)
	return res
}

func sortByClosedToCwd(packages []Package, cwd string) []Package {
	sort.Slice(packages, func(aIdx, bIdx int) bool {
		a := strings.Replace(packages[aIdx].Path, "package.json", "", 1)
		b := strings.Replace(packages[bIdx].Path, "package.json", "", 1)

		aIsSubdir := isSubdirectory(a, cwd)
		bIsSubdir := isSubdirectory(b, cwd)

		if aIsSubdir && !bIsSubdir {
			return true
		}
		if !aIsSubdir && bIsSubdir {
			return false
		}

		// if both are subdirectories (e.g the monorepo root, the shortest path should be first)
		return len(a) > len(b)
	})
	return packages
}
