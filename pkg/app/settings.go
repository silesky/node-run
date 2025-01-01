package app

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
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

type Settings struct {
	Cwd   string
	Debug bool
}

type Option func(*Settings)

func NewSettings(opts ...Option) (Settings, error) {
	settings := &Settings{}
	for _, opt := range opts {
		opt(settings)
	}
	err := settings.Validate()
	return *settings, err
}

func WithCwd(cwd string) Option {
	return func(s *Settings) {
		expanded, err := expandPath(cwd)
		if err != nil {
			log.Fatalf("Failed to expand path: %v", err)
		}
		s.Cwd = expanded
	}
}

func WithDebug(debug bool) Option {
	return func(s *Settings) {
		s.Debug = debug
	}
}

func (settings *Settings) Validate() error {
	if settings.Cwd != "" {
		if _, err := os.Stat(settings.Cwd); err != nil {
			return fmt.Errorf("--cwd is invalid: %v", settings.Cwd)
		}
	}
	return nil
}
