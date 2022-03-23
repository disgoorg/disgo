package discord

import (
	"github.com/DisgoOrg/snowflake"
	"github.com/disgoorg/disgo/json"
)

// VoiceState from Discord
type VoiceState struct {
	GuildID                 snowflake.Snowflake  `json:"guild_id,omitempty"`
	ChannelID               *snowflake.Snowflake `json:"channel_id"`
	UserID                  snowflake.Snowflake  `json:"user_id"`
	Member                  *Member              `json:"member,omitempty"`
	SessionID               string               `json:"session_id"`
	GuildDeaf               bool                 `json:"deaf"`
	GuildMute               bool                 `json:"mute"`
	SelfDeaf                bool                 `json:"self_deaf"`
	SelfMute                bool                 `json:"self_mute"`
	SelfStream              bool                 `json:"self_stream"`
	SelfVideo               bool                 `json:"self_video"`
	Suppress                bool                 `json:"suppress"`
	RequestToSpeakTimestamp *Time                `json:"request_to_speak_timestamp"`
}

type UserVoiceStateUpdate struct {
	ChannelID               snowflake.Snowflake  `json:"channel_id"`
	Suppress                *bool                `json:"suppress,omitempty"`
	RequestToSpeakTimestamp *json.Nullable[Time] `json:"request_to_speak_timestamp,omitempty"`
}
