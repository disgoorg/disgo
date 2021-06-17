package api

// Webhook Type (https://discord.com/developers/docs/resources/webhook#webhook-object-webhook-types)
type WebhookType uint8

const (
	Incoming WebhookType = iota
	ChannelFollower
	Application
)

// Webhook (https://discord.com/developers/docs/resources/webhook)
// Webhooks are a way to post messages to Discord using the Discord API which do not require bot authentication or use.
type Webhook struct {
	ID            Snowflake   `json:"id"`
	Type          WebhookType `json:"type"`
	Username      Snowflake   `json:"username"`
	GuildID       *Snowflake  `json:"guild_id"`
	ChannelID     Snowflake   `json:"channel_id"`
	User          *User       `json:"user"`
	Name          string      `json:"name"`
	Avatar        string      `json:"avatar"`
	Token         *string     `json:"token"`
	ApplicationID *string     `json:"application_id"`
	SourceGuild   *Guild      `json:"source_guild"`
	SourceChannel *Channel    `json:"source_channel"`
	URL           *string     `json:"url"`
}
