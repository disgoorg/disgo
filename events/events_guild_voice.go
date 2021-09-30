package events

import "github.com/DisgoOrg/disgo/core"

// GenericGuildVoiceEvent is called upon receiving GuildVoiceJoinEvent, GuildVoiceUpdateEvent, GuildVoiceLeaveEvent
type GenericGuildVoiceEvent struct {
	*GenericGuildMemberEvent
	VoiceState *core.VoiceState
}

// GuildVoiceJoinEvent indicates that a core.Member joined a core.Channel(requires core.GatewayIntentsGuildVoiceStates)
type GuildVoiceJoinEvent struct {
	*GenericGuildVoiceEvent
}

// GuildVoiceUpdateEvent indicates that a core.Member moved a core.Channel(requires core.GatewayIntentsGuildVoiceStates)
type GuildVoiceUpdateEvent struct {
	*GenericGuildVoiceEvent
	OldVoiceState *core.VoiceState
}

// GuildVoiceLeaveEvent indicates that a core.Member left a core.Channel(requires core.GatewayIntentsGuildVoiceStates)
type GuildVoiceLeaveEvent struct {
	*GenericGuildVoiceEvent
	OldVoiceState *core.VoiceState
}
