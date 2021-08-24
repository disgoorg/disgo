package discord

// WebhookType (https: //discord.com/developers/docs/resources/webhook#webhook-object-webhook-types)
type WebhookType int

// All WebhookType(s)
//goland:noinspection GoUnusedConst
const (
	WebhookTypeIncoming WebhookType = iota + 1
	WebhookTypeChannelFollower
	WebhookTypeApplication
)

// Webhook (https://discord.com/developers/docs/resources/webhook) is a way to post messages to Discord using the Discord API which do not require bot authentication or use.
type Webhook struct {
	ID            Snowflake       `json:"id"`
	Type          WebhookType     `json:"type"`
	Username      Snowflake       `json:"username"`
	GuildID       *Snowflake      `json:"guild_id"`
	ChannelID     Snowflake       `json:"channel_id"`
	User          User            `json:"user"`
	Name          string          `json:"name"`
	Avatar        *string         `json:"avatar"`
	Token         *string         `json:"token"`
	ApplicationID *Snowflake      `json:"application_id"`
	SourceGuild   *PartialGuild   `json:"source_guild"`
	SourceChannel *PartialChannel `json:"source_channel"`
	URL           *string         `json:"url"`
}

// WebhookCreate is used to create a Webhook
type WebhookCreate struct {
	Name   string `json:"name"`
	Avatar *Icon  `json:"avatar,omitempty"`
}

// WebhookUpdate is used to update a Webhook
type WebhookUpdate struct {
	Name   *string `json:"name,omitempty"`
	Avatar Icon    `json:"avatar,omitempty"`
}
