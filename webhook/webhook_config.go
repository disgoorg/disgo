package webhook

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/log"
)

// DefaultConfig is the default configuration for the webhook client
//goland:noinspection GoUnusedGlobalVariable
var DefaultConfig = Config{
	RestClientConfig:       &rest.DefaultConfig,
	DefaultAllowedMentions: &discord.DefaultAllowedMentions,
}

// Config is the configuration for the webhook client
type Config struct {
	Logger                 log.Logger
	RestClient             rest.Client
	RestClientConfig       *rest.Config
	WebhookService         rest.WebhookService
	EntityBuilder          EntityBuilder
	DefaultAllowedMentions *discord.AllowedMentions
}

// ConfigOpt is used to provide optional parameters to the webhook client
type ConfigOpt func(config *Config)

// Apply applies all options to the config
func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

// WithLogger sets the logger for the webhook client
//goland:noinspection GoUnusedExportedFunction
func WithLogger(logger log.Logger) ConfigOpt {
	return func(config *Config) {
		config.Logger = logger
		if config.RestClientConfig == nil {
			config.RestClientConfig = &rest.DefaultConfig
		}
		config.RestClientConfig.Logger = logger
	}
}

// WithRestClient sets the rest client for the webhook client
//goland:noinspection GoUnusedExportedFunction
func WithRestClient(restClient rest.Client) ConfigOpt {
	return func(config *Config) {
		config.RestClient = restClient
	}
}

// WithRestClientConfig sets the rest client configuration for the webhook client
//goland:noinspection GoUnusedExportedFunction
func WithRestClientConfig(restConfig rest.Config) ConfigOpt {
	return func(config *Config) {
		config.RestClientConfig = &restConfig
	}
}

// WithRestClientConfigOpts sets the rest client configuration for the webhook client
//goland:noinspection GoUnusedExportedFunction
func WithRestClientConfigOpts(opts ...rest.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		if config.RestClientConfig == nil {
			config.RestClientConfig = &rest.DefaultConfig
		}
		config.RestClientConfig.Apply(opts)
	}
}

// WithWebhookService sets the webhook service for the webhook client
//goland:noinspection GoUnusedExportedFunction
func WithWebhookService(webhookService rest.WebhookService) ConfigOpt {
	return func(config *Config) {
		config.WebhookService = webhookService
	}
}

// WithEntityBuilder sets the entity builder for the webhook client
//goland:noinspection GoUnusedExportedFunction
func WithEntityBuilder(entityBuilder EntityBuilder) ConfigOpt {
	return func(config *Config) {
		config.EntityBuilder = entityBuilder
	}
}

// WithDefaultAllowedMentions sets the default allowed mentions for the webhook client
//goland:noinspection GoUnusedExportedFunction
func WithDefaultAllowedMentions(allowedMentions discord.AllowedMentions) ConfigOpt {
	return func(config *Config) {
		config.DefaultAllowedMentions = &allowedMentions
	}
}
