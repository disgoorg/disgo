package oauth2

import (
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/log"
)

//goland:noinspection GoUnusedGlobalVariable
var DefaultConfig = Config{
	RestClientConfig: &rest.DefaultConfig,
}

type Config struct {
	Logger            log.Logger
	RestClient        rest.Client
	RestClientConfig  *rest.Config
	OAuth2Service     rest.OAuth2Service
	SessionController SessionController
	StateController   StateController
}

type ConfigOpt func(config *Config)

func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

func WithLogger(logger log.Logger) ConfigOpt {
	return func(config *Config) {
		config.Logger = logger
	}
}

func WithRestClient(restClient rest.Client) ConfigOpt {
	return func(config *Config) {
		config.RestClient = restClient
	}
}

func WithRestClientConfig(restConfig rest.Config) ConfigOpt {
	return func(config *Config) {
		config.RestClientConfig = &restConfig
	}
}

func WithRestClientConfigOpts(opts ...rest.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		if config.RestClientConfig == nil {
			config.RestClientConfig = &rest.DefaultConfig
		}
		config.RestClientConfig.Apply(opts)
	}
}

func WithOAuth2Service(oauth2service rest.OAuth2Service) ConfigOpt {
	return func(config *Config) {
		config.OAuth2Service = oauth2service
	}
}

func WithSessionController(sessionController SessionController) ConfigOpt {
	return func(config *Config) {
		config.SessionController = sessionController
	}
}

func WithStateController(stateController StateController) ConfigOpt {
	return func(config *Config) {
		config.StateController = stateController
	}
}
