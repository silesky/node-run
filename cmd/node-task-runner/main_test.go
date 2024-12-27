// Access to all the unexported functions/variable in anything in the 'package main'
package main

import (
	"os/exec"
	"strings"
	"testing"
)

func setup(t *testing.T) {
	// No need to build the binary, as we will use `go run`
}

const (
	sourceDir = "."
)

func TestVersionCommand(t *testing.T) {
	setup(t)

	cmd := exec.Command("go", "run", sourceDir, "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run command: %v", err)
	}
	substr := "Node Task Runner"
	if !strings.Contains(string(output), substr) {
		t.Errorf("Expected %q to be within %q", substr, output)
	}
}

func TestUnrecognizedCommand(t *testing.T) {
	setup(t)

	cmd := exec.Command("go", "run", sourceDir, "unknown")
	_, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("Expected error for unrecognized command, but got none")
	}
}
