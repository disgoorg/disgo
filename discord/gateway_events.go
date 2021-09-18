package discord

import (
	"time"

	"github.com/DisgoOrg/disgo/json"
)

// GatewayPayload raw GatewayEvent type
type GatewayPayload struct {
	Op Op               `json:"op"`
	S  int              `json:"s,omitempty"`
	T  GatewayEventType `json:"t,omitempty"`
	D  json.RawMessage  `json:"d,omitempty"`
}

// GatewayEventReady is the event sent by discord when you successfully Identify
type GatewayEventReady struct {
	Version     int                `json:"v"`
	SelfUser    OAuth2User         `json:"user"`
	Guilds      []Guild            `json:"guilds"`
	SessionID   string             `json:"session_id"`
	Shard       [2]int             `json:"shard,omitempty"`
	Application PartialApplication `json:"application"`
}

// GatewayEventHello is sent when we connect to the gateway
type GatewayEventHello struct {
	HeartbeatInterval time.Duration `json:"heartbeat_interval"`
}

type GatewayEventMessageReactionAdd struct {
	UserID    Snowflake     `json:"user_id"`
	ChannelID Snowflake     `json:"channel_id"`
	MessageID Snowflake     `json:"message_id"`
	GuildID   *Snowflake    `json:"guild_id"`
	Member    *Member       `json:"member"`
	Emoji     ReactionEmoji `json:"emoji"`
}

type GatewayEventMessageReactionRemove struct {
	UserID    Snowflake     `json:"user_id"`
	ChannelID Snowflake     `json:"channel_id"`
	MessageID Snowflake     `json:"message_id"`
	GuildID   *Snowflake    `json:"guild_id"`
	Emoji     ReactionEmoji `json:"emoji"`
}
