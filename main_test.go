package main

import (
    "os"
    "os/exec"
    "testing"
)

func TestMain(t *testing.T) {
    // Test the default greeting
    cmd := exec.Command("./node-task-runner")
    output, err := cmd.CombinedOutput()
    if err != nil {
        t.Fatalf("Failed to run command: %v", err)
    }
    expected := "Hello, World!\n"
    if string(output) != expected {
        t.Errorf("Expected %q, but got %q", expected, string(output))
    }

    // Test the --name flag
    cmd = exec.Command("./node-task-runner", "--name=Test")
    output, err = cmd.CombinedOutput()
    if err != nil {
        t.Fatalf("Failed to run command: %v", err)
    }
    expected = "Hello, Test!\n"
    if string(output) != expected {
        t.Errorf("Expected %q, but got %q", expected, string(output))
    }

    // Test the version command
    cmd = exec.Command("./node-task-runner", "version")
    output, err = cmd.CombinedOutput()
    if err != nil {
        t.Fatalf("Failed to run command: %v", err)
    }
    expected = "Node Task Runner CLI v1.0.0\n"
    if string(output) != expected {
        t.Errorf("Expected %q, but got %q", expected, string(output))
    }

    // Test the fzf command
    cmd = exec.Command("./node-task-runner", "fzf")
    cmd.Stdin = os.Stdin // Simulate user input
    output, err = cmd.CombinedOutput()
    if err != nil {
        t.Fatalf("Failed to run command: %v", err)
    }
    // Note: The expected output for fzf will depend on the simulated user input
}
