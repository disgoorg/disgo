package webhook

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/log"
)

// Client lets you edit/send WebhookMessage(s) or update/delete the Webhook
type Client interface {
	Logger() log.Logger

	HTTPClient() rest.HTTPClient
	WebhookService() rest.WebhookService
	EntityBuilder() EntityBuilder

	DefaultAllowedMentions() *discord.AllowedMentions
	SetDefaultAllowedMentions(allowedMentions *discord.AllowedMentions)

	GetWebhook() (*Webhook, rest.Error)
	UpdateWebhook(webhookUpdate discord.WebhookUpdate) (*Webhook, rest.Error)
	DeleteWebhook() rest.Error

	CreateMessage(messageCreate discord.MessageCreate) (*Message, rest.Error)
	CreateContent(content string) (*Message, rest.Error)
	CreateEmbeds(embeds ...discord.Embed) (*Message, rest.Error)

	UpdateMessage(messageID discord.Snowflake, messageUpdate discord.MessageUpdate) (*Message, rest.Error)
	UpdateContent(messageID discord.Snowflake, content string) (*Message, rest.Error)
	UpdateEmbeds(messageID discord.Snowflake, embeds ...discord.Embed) (*Message, rest.Error)

	DeleteMessage(id discord.Snowflake) rest.Error

	ID() discord.Snowflake
	Token() string
	URL() string
}
