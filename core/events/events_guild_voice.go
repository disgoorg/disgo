package events

import (
	"github.com/DisgoOrg/disgo/discord"
)

// GenericGuildVoiceEvent is called upon receiving GuildVoiceJoinEvent, GuildVoiceMoveEvent, GuildVoiceLeaveEvent
type GenericGuildVoiceEvent struct {
	*GenericEvent
	VoiceState discord.VoiceState
}

// GuildVoiceStateUpdateEvent indicates that the discord.VoiceState of a discord.Member has updated(requires discord.GatewayIntentsGuildVoiceStates)
type GuildVoiceStateUpdateEvent struct {
	*GenericGuildVoiceEvent
	OldVoiceState discord.VoiceState
}

// GuildVoiceJoinEvent indicates that a discord.Member joined a discord.Channel(requires discord.GatewayIntentsGuildVoiceStates)
type GuildVoiceJoinEvent struct {
	*GenericGuildVoiceEvent
}

// GuildVoiceMoveEvent indicates that a discord.Member moved a discord.Channel(requires discord.GatewayIntentsGuildVoiceStates)
type GuildVoiceMoveEvent struct {
	*GenericGuildVoiceEvent
	OldVoiceState discord.VoiceState
}

// GuildVoiceLeaveEvent indicates that a discord.Member left a discord.Channel(requires discord.GatewayIntentsGuildVoiceStates)
type GuildVoiceLeaveEvent struct {
	*GenericGuildVoiceEvent
	OldVoiceState discord.VoiceState
}

type VoiceServerUpdateEvent struct {
	*GenericEvent
	VoiceServerUpdate discord.VoiceServerUpdate
}
