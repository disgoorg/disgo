package discord

import "github.com/disgoorg/snowflake"

type VoiceServerUpdate struct {
	Token    string              `json:"token"`
	GuildID  snowflake.Snowflake `json:"guild_id"`
	Endpoint *string             `json:"endpoint"`
}
