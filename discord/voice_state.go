package discord

// VoiceState from Discord
type VoiceState struct {
	GuildID                 Snowflake  `json:"guild_id,omitempty"`
	ChannelID               *Snowflake `json:"channel_id"`
	UserID                  Snowflake  `json:"user_id"`
	Member                  *Member    `json:"member,omitempty"`
	SessionID               string     `json:"session_id"`
	GuildDeaf               bool       `json:"deaf"`
	GuildMute               bool       `json:"mute"`
	SelfDeaf                bool       `json:"self_deaf"`
	SelfMute                bool       `json:"self_mute"`
	SelfStream              bool       `json:"self_stream"`
	SelfVideo               bool       `json:"self_video"`
	Suppress                bool       `json:"suppress"`
	RequestToSpeakTimestamp Time       `json:"request_to_speak_timestamp"`
}

type UserVoiceStateUpdate struct {
	ChannelID               Snowflake     `json:"channel_id"`
	Suppress                *OptionalBool `json:"suppress,omitempty"`
	RequestToSpeakTimestamp *OptionalTime `json:"request_to_speak_timestamp,omitempty"`
}
