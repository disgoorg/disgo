package webhook

import (
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
)

func defaultConfig() config {
	return config{
		Logger:                 slog.Default(),
		DefaultAllowedMentions: &discord.DefaultAllowedMentions,
	}
}

type config struct {
	Logger                 *slog.Logger
	RestClient             rest.Client
	RestClientConfigOpts   []rest.ConfigOpt
	Webhooks               rest.Webhooks
	DefaultAllowedMentions *discord.AllowedMentions
}

// ConfigOpt is used to provide optional parameters to the webhook client
type ConfigOpt func(config *config)

func (c *config) apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
	c.Logger = c.Logger.With(slog.String("name", "webhook"))
	if c.RestClient == nil {
		c.RestClient = rest.NewClient("", append([]rest.ConfigOpt{rest.WithLogger(c.Logger)}, c.RestClientConfigOpts...)...)
	}
	if c.Webhooks == nil {
		c.Webhooks = rest.NewWebhooks(c.RestClient)
	}
}

// WithLogger sets the logger for the webhook client
func WithLogger(logger *slog.Logger) ConfigOpt {
	return func(config *config) {
		config.Logger = logger
	}
}

// WithRestClient sets the rest client for the webhook client
func WithRestClient(restClient rest.Client) ConfigOpt {
	return func(config *config) {
		config.RestClient = restClient
	}
}

// WithRestClientConfigOpts sets the rest client configuration for the webhook client
func WithRestClientConfigOpts(opts ...rest.ConfigOpt) ConfigOpt {
	return func(config *config) {
		config.RestClientConfigOpts = append(config.RestClientConfigOpts, opts...)
	}
}

// WithWebhooks sets the webhook service for the webhook client
func WithWebhooks(webhooks rest.Webhooks) ConfigOpt {
	return func(config *config) {
		config.Webhooks = webhooks
	}
}

// WithDefaultAllowedMentions sets the default allowed mentions for the webhook client
func WithDefaultAllowedMentions(allowedMentions discord.AllowedMentions) ConfigOpt {
	return func(config *config) {
		config.DefaultAllowedMentions = &allowedMentions
	}
}
