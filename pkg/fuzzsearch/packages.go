package fuzzsearch

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
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

	// append the root package.json
	packageJSONPaths = append(packageJSONPaths, filepath.Join(monorepoRoot, "package.json"))

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

type Package struct {
	Name    string
	Scripts map[string]string
}

func CreatePackageFromPath(path string) (*Package, error) {
	pkgJson, err := parsePkgJsonFile(path)
	if err != nil {
		return nil, err
	}

	return &Package{
		Name:    pkgJson.Name,
		Scripts: pkgJson.Scripts,
	}, nil
}

// GetPackages gets all packages in the current monorepo, regardless of the cwd
func GetPackages(cwd string) ([]string, error) {
	cwd, err := expandPath(cwd)

	// save the original cwd
	ogCwd := cwd

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

func sortByClosedToCwd(packagePaths []string, cwd string) []string {
	sort.Slice(packagePaths, func(aIdx, bIdx int) bool {
		a := strings.Replace(packagePaths[aIdx], "package.json", "", 1)
		b := strings.Replace(packagePaths[bIdx], "package.json", "", 1)

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
	return packagePaths
}
