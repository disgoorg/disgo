package events

import (
	"github.com/DisgoOrg/disgo/discord"
)

// GenericGuildVoiceEvent is called upon receiving GuildVoiceJoinEvent, GuildVoiceMoveEvent, GuildVoiceLeaveEvent
type GenericGuildVoiceEvent struct {
	*GenericEvent
	VoiceState discord.VoiceState
}

// GuildVoiceStateUpdateEvent indicates that the core.VoiceState of a core.Member has updated(requires core.GatewayIntentsGuildVoiceStates)
type GuildVoiceStateUpdateEvent struct {
	*GenericGuildVoiceEvent
	OldVoiceState discord.VoiceState
}

// GuildVoiceJoinEvent indicates that a core.Member joined a core.Channel(requires core.GatewayIntentsGuildVoiceStates)
type GuildVoiceJoinEvent struct {
	*GenericGuildVoiceEvent
}

// GuildVoiceMoveEvent indicates that a core.Member moved a core.Channel(requires core.GatewayIntentsGuildVoiceStates)
type GuildVoiceMoveEvent struct {
	*GenericGuildVoiceEvent
	OldVoiceState discord.VoiceState
}

// GuildVoiceLeaveEvent indicates that a core.Member left a core.Channel(requires core.GatewayIntentsGuildVoiceStates)
type GuildVoiceLeaveEvent struct {
	*GenericGuildVoiceEvent
	OldVoiceState discord.VoiceState
}

type VoiceServerUpdateEvent struct {
	*GenericEvent
	VoiceServerUpdate discord.VoiceServerUpdate
}
