package webhook

import (
	"github.com/DisgoOrg/log"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
)

// DefaultConfig is the default configuration for the webhook client
//goland:noinspection GoUnusedGlobalVariable
func DefaultConfig() *Config {
	return &Config{
		Logger:                 log.Default(),
		DefaultAllowedMentions: &discord.DefaultAllowedMentions,
	}
}

// Config is the configuration for the webhook client
type Config struct {
	Logger                 log.Logger
	RestClient             rest.Client
	RestClientConfigOpts   []rest.ConfigOpt
	WebhookService         rest.WebhookService
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
	}
}

// WithRestClient sets the rest client for the webhook client
//goland:noinspection GoUnusedExportedFunction
func WithRestClient(restClient rest.Client) ConfigOpt {
	return func(config *Config) {
		config.RestClient = restClient
	}
}

// WithRestClientConfigOpts sets the rest client configuration for the webhook client
//goland:noinspection GoUnusedExportedFunction
func WithRestClientConfigOpts(opts ...rest.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		config.RestClientConfigOpts = append(config.RestClientConfigOpts, opts...)
	}
}

// WithWebhookService sets the webhook service for the webhook client
//goland:noinspection GoUnusedExportedFunction
func WithWebhookService(webhookService rest.WebhookService) ConfigOpt {
	return func(config *Config) {
		config.WebhookService = webhookService
	}
}

// WithDefaultAllowedMentions sets the default allowed mentions for the webhook client
//goland:noinspection GoUnusedExportedFunction
func WithDefaultAllowedMentions(allowedMentions discord.AllowedMentions) ConfigOpt {
	return func(config *Config) {
		config.DefaultAllowedMentions = &allowedMentions
	}
}
