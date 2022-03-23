package events

import (
	"github.com/disgoorg/disgo/discord"
)

// GenericGuildVoiceStateEvent is called upon receiving GuildVoiceJoinEvent, GuildVoiceMoveEvent, GuildVoiceLeaveEvent
type GenericGuildVoiceStateEvent struct {
	*GenericEvent
	VoiceState discord.VoiceState
}

// GuildVoiceStateUpdateEvent indicates that the discord.VoiceState of a discord.Member has updated(requires discord.GatewayIntentsGuildVoiceStates)
type GuildVoiceStateUpdateEvent struct {
	*GenericGuildVoiceStateEvent
	OldVoiceState discord.VoiceState
}

// GuildVoiceJoinEvent indicates that a discord.Member joined a discord.Channel(requires discord.GatewayIntentsGuildVoiceStates)
type GuildVoiceJoinEvent struct {
	*GenericGuildVoiceStateEvent
}

// GuildVoiceMoveEvent indicates that a discord.Member moved a discord.Channel(requires discord.GatewayIntentsGuildVoiceStates)
type GuildVoiceMoveEvent struct {
	*GenericGuildVoiceStateEvent
	OldVoiceState discord.VoiceState
}

// GuildVoiceLeaveEvent indicates that a discord.Member left a discord.Channel(requires discord.GatewayIntentsGuildVoiceStates)
type GuildVoiceLeaveEvent struct {
	*GenericGuildVoiceStateEvent
	OldVoiceState discord.VoiceState
}

type VoiceServerUpdateEvent struct {
	*GenericEvent
	VoiceServerUpdate discord.VoiceServerUpdate
}
