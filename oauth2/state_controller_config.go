package oauth2

import (
	"log/slog"
	"time"

	"github.com/disgoorg/disgo/internal/insecurerandstr"
)

func defaultStateControllerConfig() stateControllerConfig {
	return stateControllerConfig{
		Logger:       slog.Default(),
		States:       map[string]string{},
		NewStateFunc: func() string { return insecurerandstr.RandStr(32) },
		MaxTTL:       time.Hour,
	}
}

type stateControllerConfig struct {
	Logger       *slog.Logger
	States       map[string]string
	NewStateFunc func() string
	MaxTTL       time.Duration
}

// StateControllerConfigOpt is used to pass optional parameters to NewStateController
type StateControllerConfigOpt func(config *stateControllerConfig)

func (c *stateControllerConfig) apply(opts []StateControllerConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
	c.Logger = c.Logger.With(slog.String("name", "oauth2_state_controller"))
}

// WithStateControllerLogger sets the logger for the StateController
func WithStateControllerLogger(logger *slog.Logger) StateControllerConfigOpt {
	return func(config *stateControllerConfig) {
		config.Logger = logger
	}
}

// WithStates loads states from an existing map
func WithStates(states map[string]string) StateControllerConfigOpt {
	return func(config *stateControllerConfig) {
		config.States = states
	}
}

// WithNewStateFunc sets the function which is used to generate a new random state
func WithNewStateFunc(newStateFunc func() string) StateControllerConfigOpt {
	return func(config *stateControllerConfig) {
		config.NewStateFunc = newStateFunc
	}
}

// WithMaxTTL sets the maximum time to live for a state
func WithMaxTTL(maxTTL time.Duration) StateControllerConfigOpt {
	return func(config *stateControllerConfig) {
		config.MaxTTL = maxTTL
	}
}
