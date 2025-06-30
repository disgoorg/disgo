package events

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
)

// GuildVoiceChannelEffectSend indicates that a discord.Member sent an effect in a discord.GuildVoiceChannel (requires gateway.IntentGuildVoiceStates)
type GuildVoiceChannelEffectSend struct {
	*GenericEvent
	gateway.EventVoiceChannelEffectSend
}

// GenericGuildVoiceState is called upon receiving GuildVoiceJoin, GuildVoiceMove and GuildVoiceLeave
type GenericGuildVoiceState struct {
	*GenericEvent
	VoiceState discord.VoiceState
	Member     discord.Member
}

// GuildVoiceStateUpdate indicates that the discord.VoiceState of a discord.Member has updated (requires gateway.IntentGuildVoiceStates)
type GuildVoiceStateUpdate struct {
	*GenericGuildVoiceState
	OldVoiceState discord.VoiceState
}

// GuildVoiceJoin indicates that a discord.Member joined a discord.GuildVoiceChannel (requires gateway.IntentGuildVoiceStates)
type GuildVoiceJoin struct {
	*GenericGuildVoiceState
}

// GuildVoiceMove indicates that a discord.Member was moved to a different discord.GuildVoiceChannel (requires gateway.IntentGuildVoiceStates)
type GuildVoiceMove struct {
	*GenericGuildVoiceState
	OldVoiceState discord.VoiceState
}

// GuildVoiceLeave indicates that a discord.Member left a discord.GuildVoiceChannel (requires gateway.IntentGuildVoiceStates)
type GuildVoiceLeave struct {
	*GenericGuildVoiceState
	OldVoiceState discord.VoiceState
}

// VoiceServerUpdate indicates that a voice server the bot is connected to has been changed
type VoiceServerUpdate struct {
	*GenericEvent
	gateway.EventVoiceServerUpdate
}
