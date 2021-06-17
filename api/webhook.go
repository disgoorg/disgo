package api

type WebhookType uint8

const (
	Incoming WebhookType = iota
	ChannelFollower
	Application
)

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
