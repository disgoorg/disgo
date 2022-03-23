package oauth2

import (
	"github.com/DisgoOrg/log"
	"github.com/disgoorg/disgo/rest"
)

// DefaultConfig is the configuration which is used by default
var DefaultConfig = Config{
	RestClientConfig: &rest.DefaultConfig,
}

// Config is the configuration for the OAuth2 client
type Config struct {
	Logger                log.Logger
	RestClient            rest.Client
	RestClientConfig      *rest.Config
	OAuth2Service         rest.OAuth2Service
	SessionController     SessionController
	StateControllerConfig *StateControllerConfig
	StateController       StateController
}

// ConfigOpt can be used to supply optional parameters to New
type ConfigOpt func(config *Config)

// Apply applies the given ConfigOpt(s) to the Config
func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

// WithLogger applies a custom logger to the OAuth2 client
//goland:noinspection GoUnusedExportedFunction
func WithLogger(logger log.Logger) ConfigOpt {
	return func(config *Config) {
		config.Logger = logger
	}
}

// WithRestClient applies a custom rest.Client to the OAuth2 client
//goland:noinspection GoUnusedExportedFunction
func WithRestClient(restClient rest.Client) ConfigOpt {
	return func(config *Config) {
		config.RestClient = restClient
	}
}

// WithRestClientConfig applies a custom rest.Config to the OAuth2 client
//goland:noinspection GoUnusedExportedFunction
func WithRestClientConfig(restConfig rest.Config) ConfigOpt {
	return func(config *Config) {
		config.RestClientConfig = &restConfig
	}
}

// WithRestClientConfigOpts applies rest.ConfigOpt for the rest.Client to the OAuth2 client
//goland:noinspection GoUnusedExportedFunction
func WithRestClientConfigOpts(opts ...rest.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		if config.RestClientConfig == nil {
			config.RestClientConfig = &rest.DefaultConfig
		}
		config.RestClientConfig.Apply(opts)
	}
}

// WithOAuth2Service applies a custom rest.OAuth2Service to the OAuth2 client
//goland:noinspection GoUnusedExportedFunction
func WithOAuth2Service(oauth2service rest.OAuth2Service) ConfigOpt {
	return func(config *Config) {
		config.OAuth2Service = oauth2service
	}
}

// WithSessionController applies a custom SessionController to the OAuth2 client
//goland:noinspection GoUnusedExportedFunction
func WithSessionController(sessionController SessionController) ConfigOpt {
	return func(config *Config) {
		config.SessionController = sessionController
	}
}

// WithStateController applies a custom StateController to the OAuth2 client
//goland:noinspection GoUnusedExportedFunction
func WithStateController(stateController StateController) ConfigOpt {
	return func(config *Config) {
		config.StateController = stateController
	}
}

// WithStateControllerOpts applies all StateControllerConfigOpt(s) to the StateController
//goland:noinspection GoUnusedExportedFunction
func WithStateControllerOpts(opts ...StateControllerConfigOpt) ConfigOpt {
	return func(config *Config) {
		if config.StateControllerConfig == nil {
			config.StateControllerConfig = &DefaultStateControllerConfig
		}
		config.StateControllerConfig.Apply(opts)
	}
}
