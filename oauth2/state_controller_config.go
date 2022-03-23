package oauth2

import (
	"time"

	"github.com/disgoorg/disgo/internal/insecurerandstr"
)

// DefaultStateControllerConfig is the default configuration for the StateController
func DefaultStateControllerConfig() *StateControllerConfig {
	return &StateControllerConfig{
		States:       map[string]string{},
		NewStateFunc: func() string { return insecurerandstr.RandStr(32) },
		MaxTTL:       time.Hour,
	}
}

// StateControllerConfig is the configuration for the StateController
type StateControllerConfig struct {
	States       map[string]string
	NewStateFunc func() string
	MaxTTL       time.Duration
}

// StateControllerConfigOpt is used to pass optional parameters to NewStateController
type StateControllerConfigOpt func(config *StateControllerConfig)

// Apply applies the given StateControllerConfigOpt(s) to the StateControllerConfig
func (c *StateControllerConfig) Apply(opts []StateControllerConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

// WithStates loads states from an existing map
//goland:noinspection GoUnusedExportedFunction
func WithStates(states map[string]string) StateControllerConfigOpt {
	return func(config *StateControllerConfig) {
		config.States = states
	}
}

// WithNewStateFunc sets the function which is used to generate a new random state
//goland:noinspection GoUnusedExportedFunction
func WithNewStateFunc(newStateFunc func() string) StateControllerConfigOpt {
	return func(config *StateControllerConfig) {
		config.NewStateFunc = newStateFunc
	}
}

// WithMaxTTL sets the maximum time to live for a state
//goland:noinspection GoUnusedExportedFunction
func WithMaxTTL(maxTTL time.Duration) StateControllerConfigOpt {
	return func(config *StateControllerConfig) {
		config.MaxTTL = maxTTL
	}
}
