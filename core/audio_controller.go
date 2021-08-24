package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// context support?

// AudioController lets you Connect / Disconnect from a VoiceChannel
type AudioController interface {
	Disgo() Disgo
	Connect(guildID discord.Snowflake, channelID discord.Snowflake) error
	Disconnect(guildID discord.Snowflake) error
}
