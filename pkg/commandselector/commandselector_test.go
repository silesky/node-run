package commandselector

import (
	helpers "node-task-runner/pkg/testhelpers"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunCommand(t *testing.T) {
	runner := InteractivePackageCommandRunner{
		command: "echo Hello, World!",
	}

	output := helpers.CaptureOutput(func() {
		runner.runCommand()
	})

	assert.Contains(t, output, "Hello, World!")
}
