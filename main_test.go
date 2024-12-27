package main

import (
	"os/exec"
	"testing"
)

func setup(t *testing.T) {
	// No need to build the binary, as we will use `go run`
}

const (
	sourceDir = "./cmd/node-task-runner"
)

func TestMain(t *testing.T) {
	setup(t)

	// Test the version command
	cmd := exec.Command("go", "run", sourceDir, "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run command: %v", err)
	}
	expected := "Node Task Runner CLI v1.0.0\n"
	if string(output) != expected {
		t.Errorf("Expected %q, but got %q", expected, string(output))
	}
}
