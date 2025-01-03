package app

import (
	"errors"
	"fmt"
	"node-task-runner/pkg/commandselector"
	"node-task-runner/pkg/logger"
	"os"
)

func Run(settings Settings) {
	logger.SetDebug(settings.Debug)
	cwd := settings.Cwd
	if cwd == "" {
		var err error
		cwd, err = os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
			return
		}
	}
	logger.Debugf("Current directory: %s", cwd)

	selectedCommand, project, err := commandselector.RunCommandSelectorPrompt(cwd)
	if err != nil {
		if errors.Is(err, commandselector.ErrUserAbort) {
			return
		}
		fmt.Fprintf(os.Stderr, "Error getting commands: %v\n", err)
		return
	} else {
		commandselector.Executor(selectedCommand, project)
	}
}
