package app

import (
	"context"
	"fmt"
	"log"
	"node-task-runner/pkg/logger"
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

type contextKey string

const settingsKey contextKey = "settings"

func NewSettingsContext(ctx context.Context, settings Settings) context.Context {
	return context.WithValue(ctx, settingsKey, settings)
}

func FromSettingsContext(ctx context.Context) Settings {
	settings, ok := ctx.Value(settingsKey).(Settings)
	if !ok {
		panic("invariant: settings does not exist")
	}
	logger.Debugf("Settings: %v", settings)
	return settings
}

func (settings *Settings) Validate() error {
	if settings.Cwd != "" {
		if _, err := os.Stat(settings.Cwd); err != nil {
			return fmt.Errorf("--cwd is invalid: %v", settings.Cwd)
		}
	}
	return nil
}
