package httpinteraction

import (
	"log/slog"
)

func defaultConfig() config {
	return config{
		Logger:   slog.Default(),
		Endpoint: "/interactions/callback",
		Verifier: DefaultVerifier{},
	}
}

type config struct {
	Logger          *slog.Logger
	Endpoint        string
	Verifier        KeyVerifier
	EnableRawEvents bool
}

// ConfigOpt is a type alias for a function that takes a config and is used to configure your Server.
type ConfigOpt func(config *config)

func (c *config) apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
	c.Logger = c.Logger.With(slog.String("name", "httpserver"))
}

func WithDefault() ConfigOpt {
	return func(config *config) {}
}

// WithLogger sets the Logger of the config.
func WithLogger(logger *slog.Logger) ConfigOpt {
	return func(config *config) {
		config.Logger = logger
	}
}

// WithEndpoint sets the endpoint discord will send interactions to.
func WithEndpoint(endpoint string) ConfigOpt {
	return func(config *config) {
		config.Endpoint = endpoint
	}
}

// WithVerifier sets the KeyVerifier of the config.
func WithVerifier(verifier KeyVerifier) ConfigOpt {
	return func(config *config) {
		config.Verifier = verifier
	}
}
