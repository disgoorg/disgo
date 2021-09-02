package webhook

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/log"
)

// DefaultAllowedMentions gives you the default AllowedMentions for a Message
var DefaultAllowedMentions = discord.AllowedMentions{
	Parse:       []discord.AllowedMentionType{discord.AllowedMentionTypeUsers, discord.AllowedMentionTypeRoles, discord.AllowedMentionTypeEveryone},
	Roles:       []discord.Snowflake{},
	Users:       []discord.Snowflake{},
	RepliedUser: true,
}

// Client lets you edit/send WebhookMessage(s) or update/delete the Webhook
type Client interface {
	Logger() log.Logger

	RestClient() rest.Client
	WebhookService() rest.WebhookService
	EntityBuilder() EntityBuilder

	DefaultAllowedMentions() *discord.AllowedMentions

	GetWebhook(opts ...rest.RequestOpt) (*Webhook, rest.Error)
	UpdateWebhook(webhookUpdate discord.WebhookUpdate, opts ...rest.RequestOpt) (*Webhook, rest.Error)
	DeleteWebhook(opts ...rest.RequestOpt) rest.Error

	CreateMessage(messageCreate discord.WebhookMessageCreate, opts ...rest.RequestOpt) (*Message, rest.Error)
	CreateContent(content string, opts ...rest.RequestOpt) (*Message, rest.Error)
	CreateEmbeds(embeds []discord.Embed, opts ...rest.RequestOpt) (*Message, rest.Error)

	UpdateMessage(messageID discord.Snowflake, messageUpdate discord.WebhookMessageUpdate, opts ...rest.RequestOpt) (*Message, rest.Error)
	UpdateContent(messageID discord.Snowflake, content string, opts ...rest.RequestOpt) (*Message, rest.Error)
	UpdateEmbeds(messageID discord.Snowflake, embeds []discord.Embed, opts ...rest.RequestOpt) (*Message, rest.Error)

	DeleteMessage(messageID discord.Snowflake, opts ...rest.RequestOpt) rest.Error

	ID() discord.Snowflake
	Token() string
	URL() string
}
