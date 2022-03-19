package webhook

import (
	"context"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/snowflake"
)

type Client interface {
	// ID returns the configured Webhook id
	ID() snowflake.Snowflake
	// Token returns the configured Webhook token
	Token() string
	// URL returns the full Webhook URL
	URL() string
	// Config returns the configured Config
	Config() Config
	// Close closes all connections the Webhook Client has open
	Close(ctx context.Context)
	// WebhookService returns the underlying rest.WebhookService
	WebhookService() rest.WebhookService

	// GetWebhook fetches the current Webhook from discord
	GetWebhook(opts ...rest.RequestOpt) (*discord.IncomingWebhook, error)
	// UpdateWebhook updates the current Webhook
	UpdateWebhook(webhookUpdate discord.WebhookUpdateWithToken, opts ...rest.RequestOpt) (*discord.IncomingWebhook, error)
	// DeleteWebhook deletes the current Webhook
	DeleteWebhook(opts ...rest.RequestOpt) error

	// CreateMessage creates a new Message from the discord.WebhookMessageCreate
	CreateMessage(messageCreate discord.WebhookMessageCreate, opts ...rest.RequestOpt) (*discord.Message, error)
	// CreateMessageInThread creates a new Message from the discord.WebhookMessageCreate in the provided thread
	CreateMessageInThread(messageCreate discord.WebhookMessageCreate, threadID snowflake.Snowflake, opts ...rest.RequestOpt) (*discord.Message, error)
	// CreateContent creates a new Message from the provided content
	CreateContent(content string, opts ...rest.RequestOpt) (*discord.Message, error)
	// CreateEmbeds creates a new Message from the provided discord.Embed(s)
	CreateEmbeds(embeds []discord.Embed, opts ...rest.RequestOpt) (*discord.Message, error)

	// UpdateMessage updates an already sent Webhook Message with the discord.WebhookMessageUpdate
	UpdateMessage(messageID snowflake.Snowflake, messageUpdate discord.WebhookMessageUpdate, opts ...rest.RequestOpt) (*discord.Message, error)
	// UpdateMessageInThread updates an already sent Webhook Message with the discord.WebhookMessageUpdate in the provided thread
	UpdateMessageInThread(messageID snowflake.Snowflake, messageUpdate discord.WebhookMessageUpdate, threadID snowflake.Snowflake, opts ...rest.RequestOpt) (*discord.Message, error)
	// UpdateContent updates an already sent Webhook Message with the content
	UpdateContent(messageID snowflake.Snowflake, content string, opts ...rest.RequestOpt) (*discord.Message, error)
	// UpdateEmbeds updates an already sent Webhook Message with the discord.Embed(s)
	UpdateEmbeds(messageID snowflake.Snowflake, embeds []discord.Embed, opts ...rest.RequestOpt) (*discord.Message, error)

	// DeleteMessage deletes an already sent Webhook Message
	DeleteMessage(messageID snowflake.Snowflake, opts ...rest.RequestOpt) error
	// DeleteMessageInThread deletes an already sent Webhook Message in the provided thread
	DeleteMessageInThread(messageID snowflake.Snowflake, threadID snowflake.Snowflake, opts ...rest.RequestOpt) error
}
