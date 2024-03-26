package events

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

// GenericGuildSoundboardSound is called upon receiving GuildSoundboardSoundCreate and GuildSoundboardSoundUpdate
type GenericGuildSoundboardSound struct {
	*GenericEvent
	discord.SoundboardSound
}

// GuildSoundboardSoundCreate indicates that a discord.SoundboardSound was created in a discord.Guild
type GuildSoundboardSoundCreate struct {
	*GenericGuildSoundboardSound
}

// GuildSoundboardSoundUpdate indicates that a discord.SoundboardSound was updated in a discord.Guild
type GuildSoundboardSoundUpdate struct {
	*GenericGuildSoundboardSound
}

// GuildSoundboardSoundDelete indicates that a discord.SoundboardSound was deleted in a discord.Guild
type GuildSoundboardSoundDelete struct {
	*GenericEvent
	SoundID snowflake.ID
	GuildID snowflake.ID
}

// SoundboardSounds is a response to gateway.OpcodeRequestSoundboardSounds
type SoundboardSounds struct {
	*GenericEvent
	SoundboardSounds []discord.SoundboardSound
	GuildID          snowflake.ID
}
