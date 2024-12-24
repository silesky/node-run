package app

import (
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

	settings, _ := NewSettings(WithCwd(tmpDir), WithDebug(true))

	Run(settings)
	// Add assertions as needed to verify the behavior of Run
}
