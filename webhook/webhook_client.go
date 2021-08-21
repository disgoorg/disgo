package webhook

import (
	"context"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/log"
)

// Client lets you edit/send WebhookMessage(s) or update/delete the Webhook
type Client interface {
	Logger() log.Logger

	RestClient() rest.Client
	WebhookService() rest.WebhookService
	EntityBuilder() EntityBuilder

	DefaultAllowedMentions() *discord.AllowedMentions
	SetDefaultAllowedMentions(allowedMentions *discord.AllowedMentions)

	GetWebhook(ctx context.Context) (*Webhook, rest.Error)
	UpdateWebhook(ctx context.Context, webhookUpdate discord.WebhookUpdate) (*Webhook, rest.Error)
	DeleteWebhook(ctx context.Context) rest.Error

	CreateMessage(ctx context.Context, messageCreate discord.WebhookMessageCreate) (*Message, rest.Error)
	CreateContent(ctx context.Context, content string) (*Message, rest.Error)
	CreateEmbeds(ctx context.Context, embeds ...discord.Embed) (*Message, rest.Error)

	UpdateMessage(ctx context.Context, messageID discord.Snowflake, messageUpdate discord.WebhookMessageUpdate) (*Message, rest.Error)
	UpdateContent(ctx context.Context, messageID discord.Snowflake, content string) (*Message, rest.Error)
	UpdateEmbeds(ctx context.Context, messageID discord.Snowflake, embeds ...discord.Embed) (*Message, rest.Error)

	DeleteMessage(ctx context.Context, id discord.Snowflake) rest.Error

	ID() discord.Snowflake
	Token() string
	URL() string
}
