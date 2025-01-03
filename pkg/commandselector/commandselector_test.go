package commandselector

import (
	helpers "node-task-runner/pkg/testhelpers"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandSelector(t *testing.T) {
	runner := InteractivePackageCommandRunner{
		command: "echo Hello, World!",
	}

	// Capture stdout
	output := helpers.CaptureOutput(func() {
		runner.runCommand()
	})

	// Assert on stdout
	assert.Contains(t, output, "Hello, World!")
}
