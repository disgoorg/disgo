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
	ID            Snowflake             `json:"id"`
	Type          WebhookType           `json:"type"`
	GuildID       *Snowflake            `json:"guild_id"`
	ChannelID     *Snowflake            `json:"channel_id"`
	User          *User                 `json:"user"`
	Name          *string               `json:"name"`
	Avatar        *string               `json:"avatar"`
	Token         string                `json:"token,omitempty"`
	ApplicationID *Snowflake            `json:"application_id"`
	SourceGuild   *WebhookSourceGuild   `json:"source_guild"`
	SourceChannel *WebhookSourceChannel `json:"source_channel"`
	URL           string                `json:"url,omitempty"`
}

type WebhookSourceGuild struct {
	ID   Snowflake `json:"id"`
	Name string    `json:"name"`
	Icon *string   `json:"icon"`
}

type WebhookSourceChannel struct {
	ID   Snowflake `json:"id"`
	Name string    `json:"name"`
}

// WebhookCreate is used to create a Webhook
type WebhookCreate struct {
	Name   string `json:"name"`
	Avatar Icon   `json:"avatar,omitempty"`
}

// WebhookUpdate is used to update a Webhook
type WebhookUpdate struct {
	Name   *string       `json:"name,omitempty"`
	Avatar *OptionalIcon `json:"avatar,omitempty"`
}
