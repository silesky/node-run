package app

type Settings struct {
	Cwd string
	// Add more fields as needed
}

// Option is a function that sets an option on the Settings struct
type Option func(*Settings)

// NewSettings creates a new Settings struct with the provided options
func NewSettings(opts ...Option) *Settings {
	settings := &Settings{}
	for _, opt := range opts {
		opt(settings)
	}
	return settings
}

// WithCwd sets the Cwd field on the Settings struct
func WithCwd(cwd string) Option {
	return func(s *Settings) {
		s.Cwd = cwd
	}
}
