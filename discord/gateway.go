package discord

type Gateway struct {
	URL string `json:"url"`
}

type GatewayBot struct {
	URL               string            `json:"url"`
	Shards            int               `json:"shards"`
	SessionStartLimit SessionStartLimit `json:"session_start_limit"`
}

type SessionStartLimit struct {
	Total          int `json:"total"`
	Remaining      int `json:"remaining"`
	ResetAfter     int `json:"reset_after"`
	MaxConcurrency int `json:"max_concurrency"`
}

type ChannelPinsUpdateGatewayEvent struct {
	GuildID          *Snowflake `json:"guild_id"`
	ChannelID        Snowflake  `json:"channel_id"`
	LastPinTimestamp *Time      `json:"last_pin_timestamp"`
}

type GuildBanAddGatewayEvent struct {
	GuildID Snowflake `json:"guild_id"`
	User    User      `json:"user"`
}

type GuildBanRemoveGatewayEvent struct {
	GuildID Snowflake `json:"guild_id"`
	User    User      `json:"user"`
}

type GuildEmojisUpdateGatewayEvent struct {
	GuildID Snowflake `json:"guild_id"`
	Emojis  []Emoji   `json:"emojis"`
}

type GuildStickersUpdateGatewayEvent struct {
	GuildID  Snowflake `json:"guild_id"`
	Stickers []Sticker `json:"stickers"`
}

type GuildIntegrationsUpdateGatewayEvent struct {
	GuildID Snowflake `json:"guild_id"`
}

type GuildMemberRemoveGatewayEvent struct {
	GuildID Snowflake `json:"guild_id"`
	User    User      `json:"user"`
}

type GuildRoleCreateGatewayEvent struct {
	GuildID Snowflake `json:"guild_id"`
	Role    Role      `json:"role"`
}

type GuildRoleDeleteGatewayEvent struct {
	GuildID Snowflake `json:"guild_id"`
	RoleID  Snowflake `json:"role_id"`
}

type GuildRoleUpdateGatewayEvent struct {
	GuildID Snowflake `json:"guild_id"`
	Role    Role      `json:"role"`
}

type InviteDeleteGatewayEvent struct {
	ChannelID Snowflake  `json:"channel_id"`
	GuildID   *Snowflake `json:"guild_id"`
	Code      string     `json:"code"`
}

type MessageDeleteGatewayEvent struct {
	MessageID Snowflake  `json:"id"`
	GuildID   *Snowflake `json:"guild_id,omitempty"`
	ChannelID Snowflake  `json:"channel_id"`
}

type TypingStartGatewayEvent struct {
	ChannelID Snowflake
	GuildID   *Snowflake
	UserID    Snowflake
	Timestamp Time
	Member    *Member
	User      User
}

type WebhooksUpdateGatewayEvent struct {
	GuildID   Snowflake `json:"guild_id"`
	ChannelID Snowflake `json:"channel_id"`
}

type InvalidSessionGatewayEvent struct {
	bool
}

type IntegrationCreateGatewayEvent struct {
	Integration
	GuildID Snowflake `json:"guild_id"`
}

type IntegrationUpdateGatewayEvent struct {
	Integration
	GuildID Snowflake `json:"guild_id"`
}

type IntegrationDeleteGatewayEvent struct {
	ID            Snowflake `json:"id"`
	GuildID       Snowflake `json:"guild_id"`
	ApplicationID Snowflake `json:"application_id"`
}
