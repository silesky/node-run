package app

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestExpandPath(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}

	tests := []struct {
		input    string
		expected string
	}{
		{"~/test", filepath.Join(homeDir, "test")},
		{"/absolute/path", "/absolute/path"},
	}

	for _, test := range tests {
		result, err := expandPath(test.input)
		if err != nil {
			t.Fatalf("Failed to expand path: %v", err)
		}
		if result != test.expected {
			t.Errorf("Expected %q, but got %q", test.expected, result)
		}
	}
}

func TestHasWorkspacesArray(t *testing.T) {
	// Create a temporary package.json file
	packageJSON := `{
		"workspaces": ["packages/*"]
	}`
	tmpFile, err := os.CreateTemp("", "package.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(packageJSON)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	if !hasWorkspacesArray(tmpFile.Name()) {
		t.Errorf("Expected workspaces array to be found")
	}
}

func TestParseWorkspacesArray(t *testing.T) {
	// Create a temporary package.json file
	packageJSON := `{
		"workspaces": ["packages/*"]
	}`
	tmpDir := t.TempDir()
	packageJSONPath := filepath.Join(tmpDir, "package.json")
	if err := os.WriteFile(packageJSONPath, []byte(packageJSON), 0644); err != nil {
		t.Fatalf("Failed to write package.json: %v", err)
	}

	// Create dummy package.json files in workspaces
	os.MkdirAll(filepath.Join(tmpDir, "packages", "pkg1"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "packages", "pkg1", "package.json"), []byte(`{}`), 0644)
	os.MkdirAll(filepath.Join(tmpDir, "packages", "pkg2"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "packages", "pkg2", "package.json"), []byte(`{}`), 0644)

	paths, err := parseWorkspacesArray(tmpDir)
	if err != nil {
		t.Fatalf("Failed to parse workspaces array: %v", err)
	}

	expectedPaths := []string{
		filepath.Join(tmpDir, "packages", "pkg1", "package.json"),
		filepath.Join(tmpDir, "packages", "pkg2", "package.json"),
	}

	for _, expected := range expectedPaths {
		found := false
		for _, path := range paths {
			if path == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected path %q not found", expected)
		}
	}
}

func TestFindMonorepoRoot(t *testing.T) {
	// Create a temporary directory structure
	tmpDir := t.TempDir()
	os.MkdirAll(filepath.Join(tmpDir, "subdir", "subsubdir"), 0755)
	packageJSONPath := filepath.Join(tmpDir, "package.json")
	if err := os.WriteFile(packageJSONPath, []byte(`{"workspaces": ["packages/*"]}`), 0644); err != nil {
		t.Fatalf("Failed to write package.json: %v", err)
	}

	root, err := findMonorepoRoot(filepath.Join(tmpDir, "subdir", "subsubdir"))
	if err != nil {
		t.Fatalf("Failed to find monorepo root: %v", err)
	}
	if root != tmpDir {
		t.Errorf("Expected root %q, but got %q", tmpDir, root)
	}
}

func TestGetPackages(t *testing.T) {
	// Create a temporary directory structure
	tmpDir := t.TempDir()
	os.MkdirAll(filepath.Join(tmpDir, "packages", "pkg1"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "packages", "pkg1", "package.json"), []byte(`{}`), 0644)
	os.MkdirAll(filepath.Join(tmpDir, "packages", "pkg2"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "packages", "pkg2", "package.json"), []byte(`{}`), 0644)
	packageJSONPath := filepath.Join(tmpDir, "package.json")
	if err := os.WriteFile(packageJSONPath, []byte(`{"workspaces": ["packages/*"]}`), 0644); err != nil {
		t.Fatalf("Failed to write package.json: %v", err)
	}

	packages, err := GetPackages(tmpDir)
	if err != nil {
		t.Fatalf("Failed to get packages: %v", err)
	}

	expectedPackages := []string{
		filepath.Join(tmpDir, "packages", "pkg1", "package.json"),
		filepath.Join(tmpDir, "packages", "pkg2", "package.json"),
	}

	for _, expected := range expectedPackages {
		found := false
		for _, pkg := range packages {
			if pkg == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected package %q not found", expected)
		}
	}
}

func TestRun(t *testing.T) {
	// Create a temporary directory structure
	tmpDir := t.TempDir()
	os.MkdirAll(filepath.Join(tmpDir, "packages", "pkg1"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "packages", "pkg1", "package.json"), []byte(`{}`), 0644)
	os.MkdirAll(filepath.Join(tmpDir, "packages", "pkg2"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "packages", "pkg2", "package.json"), []byte(`{}`), 0644)
	packageJSONPath := filepath.Join(tmpDir, "package.json")
	if err := os.WriteFile(packageJSONPath, []byte(`{"workspaces": ["packages/*"]}`), 0644); err != nil {
		t.Fatalf("Failed to write package.json: %v", err)
	}

	settings := NewSettings(WithCwd(tmpDir), WithDebug(true))
	ctx := NewSettingsContext(context.Background(), settings)

	Run(ctx)
	// Add assertions as needed to verify the behavior of Run
}
