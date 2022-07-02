package events

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
)

// GenericGuildVoiceState is called upon receiving GuildVoiceJoin , GuildVoiceMove , GuildVoiceLeave
type GenericGuildVoiceState struct {
	*GenericEvent
	VoiceState discord.VoiceState
	Member     discord.Member
}

// GuildVoiceStateUpdate indicates that the discord.VoiceState of a discord.Member has updated(requires gateway.IntentsGuildVoiceStates)
type GuildVoiceStateUpdate struct {
	*GenericGuildVoiceState
	OldVoiceState discord.VoiceState
}

// GuildVoiceJoin indicates that a discord.Member joined a discord.Channel(requires gateway.IntentsGuildVoiceStates)
type GuildVoiceJoin struct {
	*GenericGuildVoiceState
}

// GuildVoiceMove indicates that a discord.Member moved a discord.Channel(requires gateway.IntentsGuildVoiceStates)
type GuildVoiceMove struct {
	*GenericGuildVoiceState
	OldVoiceState discord.VoiceState
}

// GuildVoiceLeave indicates that a discord.Member left a discord.Channel(requires gateway.IntentsGuildVoiceStates)
type GuildVoiceLeave struct {
	*GenericGuildVoiceState
	OldVoiceState discord.VoiceState
}

// VoiceServerUpdate indicates that a voice server the bot is connected to has been changed
type VoiceServerUpdate struct {
	*GenericEvent
	gateway.EventVoiceServerUpdate
}
