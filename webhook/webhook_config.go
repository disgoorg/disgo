package webhook

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/log"
)

//goland:noinspection GoUnusedGlobalVariable
var DefaultConfig = Config{
	RestClientConfig:       &rest.DefaultConfig,
	DefaultAllowedMentions: &DefaultAllowedMentions,
}

type Config struct {
	Logger                 log.Logger
	RestClient             rest.Client
	RestClientConfig       *rest.Config
	WebhookService         rest.WebhookService
	EntityBuilder          EntityBuilder
	DefaultAllowedMentions *discord.AllowedMentions
}

type ConfigOpt func(config *Config)

func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithLogger(logger log.Logger) ConfigOpt {
	return func(config *Config) {
		config.Logger = logger
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithRestClient(restClient rest.Client) ConfigOpt {
	return func(config *Config) {
		config.RestClient = restClient
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithRestClientConfig(restConfig rest.Config) ConfigOpt {
	return func(config *Config) {
		config.RestClientConfig = &restConfig
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithRestClientConfigOpts(opts ...rest.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		if config.RestClientConfig == nil {
			config.RestClientConfig = &rest.DefaultConfig
		}
		config.RestClientConfig.Apply(opts)
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithWebhookService(webhookService rest.WebhookService) ConfigOpt {
	return func(config *Config) {
		config.WebhookService = webhookService
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithEntityBuilder(entityBuilder EntityBuilder) ConfigOpt {
	return func(config *Config) {
		config.EntityBuilder = entityBuilder
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithDefaultAllowedMentions(allowedMentions discord.AllowedMentions) ConfigOpt {
	return func(config *Config) {
		config.DefaultAllowedMentions = &allowedMentions
	}
}
