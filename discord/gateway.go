package discord

import (
	"time"

	"github.com/DisgoOrg/disgo/json"
)

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

type GuildMembersChunkGatewayEvent struct {
	GuildID    Snowflake   `json:"guild_id"`
	Members    []Member    `json:"members"`
	ChunkIndex int         `json:"chunk_index"`
	ChunkCount int         `json:"chunk_count"`
	NotFound   []Snowflake `json:"not_found"`
	Presences  []Presence  `json:"presences"`
	Nonce      string      `json:"nonce"`
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
	ID        Snowflake  `json:"id"`
	ChannelID Snowflake  `json:"channel_id"`
	GuildID   *Snowflake `json:"guild_id,omitempty"`
}

type MessageDeleteBulkGatewayEvent struct {
	IDs       []Snowflake `json:"id"`
	ChannelID Snowflake   `json:"channel_id"`
	GuildID   *Snowflake  `json:"guild_id,omitempty"`
}

type TypingStartGatewayEvent struct {
	ChannelID Snowflake
	GuildID   *Snowflake
	UserID    Snowflake
	Timestamp time.Time
	Member    *Member
	User      User
}

func (e *TypingStartGatewayEvent) UnmarshalJSON(data []byte) error {
	type typingStartGatewayEvent TypingStartGatewayEvent
    var v struct {
        Timestamp int64 `json:"timestamp"`
        typingStartGatewayEvent
    }
    if err := json.Unmarshal(data, &v); err != nil {
        return err
    }
    e.Timestamp = time.Unix(v.Timestamp, 0)
    return nil
}

type WebhooksUpdateGatewayEvent struct {
	GuildID   Snowflake `json:"guild_id"`
	ChannelID Snowflake `json:"channel_id"`
}

type IntegrationCreateGatewayEvent struct {
	Integration
	GuildID Snowflake `json:"guild_id"`
}

func (e *IntegrationCreateGatewayEvent) UnmarshalJSON(data []byte) error {
	type integrationCreateGatewayEvent IntegrationCreateGatewayEvent
	var v struct {
		UnmarshalIntegration
		integrationCreateGatewayEvent
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*e = IntegrationCreateGatewayEvent(v.integrationCreateGatewayEvent)

	e.Integration = v.UnmarshalIntegration.Integration
	return nil
}

type IntegrationUpdateGatewayEvent struct {
	Integration
	GuildID Snowflake `json:"guild_id"`
}

func (e *IntegrationUpdateGatewayEvent) UnmarshalJSON(data []byte) error {
	type integrationUpdateGatewayEvent IntegrationUpdateGatewayEvent
	var v struct {
		UnmarshalIntegration
		integrationUpdateGatewayEvent
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*e = IntegrationUpdateGatewayEvent(v.integrationUpdateGatewayEvent)

	e.Integration = v.UnmarshalIntegration.Integration
	return nil
}

type IntegrationDeleteGatewayEvent struct {
	ID            Snowflake  `json:"id"`
	GuildID       Snowflake  `json:"guild_id"`
	ApplicationID *Snowflake `json:"application_id"`
}
