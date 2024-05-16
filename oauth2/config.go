package oauth2

import (
	"log/slog"

	"github.com/disgoorg/disgo/rest"
)

// DefaultConfig is the configuration which is used by default
func DefaultConfig() *Config {
	return &Config{
		Logger: slog.Default(),
	}
}

// Config is the configuration for the OAuth2 client
type Config struct {
	Logger                    *slog.Logger
	RestClient                rest.Client
	RestClientConfigOpts      []rest.ConfigOpt
	OAuth2                    rest.OAuth2
	StateController           StateController
	StateControllerConfigOpts []StateControllerConfigOpt
}

// ConfigOpt can be used to supply optional parameters to New
type ConfigOpt func(config *Config)

// Apply applies the given ConfigOpt(s) to the Config
func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
	if c.RestClient == nil {
		c.RestClient = rest.NewClient("", append([]rest.ConfigOpt{rest.WithLogger(c.Logger)}, c.RestClientConfigOpts...)...)
	}
	if c.OAuth2 == nil {
		c.OAuth2 = rest.NewOAuth2(c.RestClient)
	}
	if c.StateController == nil {
		c.StateController = NewStateController(append([]StateControllerConfigOpt{WithStateControllerLogger(c.Logger)}, c.StateControllerConfigOpts...)...)
	}
}

// WithLogger applies a custom logger to the OAuth2 client
func WithLogger(logger *slog.Logger) ConfigOpt {
	return func(config *Config) {
		config.Logger = logger
	}
}

// WithRestClient applies a custom rest.Client to the OAuth2 client
func WithRestClient(restClient rest.Client) ConfigOpt {
	return func(config *Config) {
		config.RestClient = restClient
	}
}

// WithRestClientConfigOpts applies rest.ConfigOpt for the rest.Client to the OAuth2 client
func WithRestClientConfigOpts(opts ...rest.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		config.RestClientConfigOpts = append(config.RestClientConfigOpts, opts...)
	}
}

// WithOAuth2 applies a custom rest.OAuth2 to the OAuth2 client
func WithOAuth2(oauth2 rest.OAuth2) ConfigOpt {
	return func(config *Config) {
		config.OAuth2 = oauth2
	}
}

// WithStateController applies a custom StateController to the OAuth2 client
func WithStateController(stateController StateController) ConfigOpt {
	return func(config *Config) {
		config.StateController = stateController
	}
}

// WithStateControllerOpts applies all StateControllerConfigOpt(s) to the StateController
func WithStateControllerOpts(opts ...StateControllerConfigOpt) ConfigOpt {
	return func(config *Config) {
		config.StateControllerConfigOpts = append(config.StateControllerConfigOpts, opts...)
	}
}
