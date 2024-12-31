package app

import (
	"context"
)

type Settings struct {
	Cwd   string
	Debug bool
	// Add more fields as needed
}

// Option is a function that sets an option on the Settings struct
type Option func(*Settings)

// NewSettings creates a new Settings struct with the provided options
func NewSettings(opts ...Option) Settings {
	settings := &Settings{}
	for _, opt := range opts {
		opt(settings)
	}
	return *settings
}

// WithCwd sets the Cwd field on the Settings struct
func WithCwd(cwd string) Option {
	return func(s *Settings) {
		// check if cwd is a valid directory

	}
}

// WithDebug sets the Debug field on the Settings struct
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
	return settings
}
