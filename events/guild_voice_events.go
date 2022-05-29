package events

import (
	"github.com/disgoorg/disgo/discord"
)

// GenericGuildVoiceState is called upon receiving GuildVoiceJoin , GuildVoiceMove , GuildVoiceLeave
type GenericGuildVoiceState struct {
	*GenericEvent
	VoiceState discord.VoiceState
	Member     discord.Member
}

// GuildVoiceStateUpdate indicates that the discord.VoiceState of a discord.Member has updated(requires discord.GatewayIntentsGuildVoiceStates)
type GuildVoiceStateUpdate struct {
	*GenericGuildVoiceState
	OldVoiceState discord.VoiceState
}

// GuildVoiceJoin indicates that a discord.Member joined a discord.Channel(requires discord.GatewayIntentsGuildVoiceStates)
type GuildVoiceJoin struct {
	*GenericGuildVoiceState
}

// GuildVoiceMove indicates that a discord.Member moved a discord.Channel(requires discord.GatewayIntentsGuildVoiceStates)
type GuildVoiceMove struct {
	*GenericGuildVoiceState
	OldVoiceState discord.VoiceState
}

// GuildVoiceLeave indicates that a discord.Member left a discord.Channel(requires discord.GatewayIntentsGuildVoiceStates)
type GuildVoiceLeave struct {
	*GenericGuildVoiceState
	OldVoiceState discord.VoiceState
}

// VoiceServerUpdate indicates that a voice server the bot is connected to has been changed
type VoiceServerUpdate struct {
	*GenericEvent
	discord.VoiceServerUpdate
}
