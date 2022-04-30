package discord

import "github.com/disgoorg/snowflake/v2"

type VoiceServerUpdate struct {
	Token    string       `json:"token"`
	GuildID  snowflake.ID `json:"guild_id"`
	Endpoint *string      `json:"endpoint"`
}
