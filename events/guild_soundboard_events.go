package events

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

// GenericGuildSoundboardSound is called upon receiving GuildSoundboardSoundCreate and GuildSoundboardSoundUpdate (requires gateway.IntentGuildExpressions)
type GenericGuildSoundboardSound struct {
	*GenericEvent
	discord.SoundboardSound
}

// GuildSoundboardSoundCreate indicates that a discord.SoundboardSound was created in a discord.Guild (requires gateway.IntentGuildExpressions)
type GuildSoundboardSoundCreate struct {
	*GenericGuildSoundboardSound
}

// GuildSoundboardSoundUpdate indicates that a discord.SoundboardSound was updated in a discord.Guild (requires gateway.IntentGuildExpressions)
type GuildSoundboardSoundUpdate struct {
	*GenericGuildSoundboardSound
	OldGuildSoundboardSound discord.SoundboardSound
}

// GuildSoundboardSoundDelete indicates that a discord.SoundboardSound was deleted in a discord.Guild (requires gateway.IntentGuildExpressions)
type GuildSoundboardSoundDelete struct {
	*GenericEvent
	SoundID snowflake.ID
	GuildID snowflake.ID
}

// GuildSoundboardSoundsUpdate indicates when multiple discord.Guild soundboard sounds were updated (requires gateway.IntentGuildExpressions)
type GuildSoundboardSoundsUpdate struct {
	*GenericEvent
	SoundboardSounds []discord.SoundboardSound
}

// SoundboardSounds is a response to gateway.OpcodeRequestSoundboardSounds
type SoundboardSounds struct {
	*GenericEvent
	SoundboardSounds []discord.SoundboardSound
	GuildID          snowflake.ID
}
