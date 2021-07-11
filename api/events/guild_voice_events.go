package events

import (
	"github.com/DisgoOrg/disgo/api"
)

// GenericGuildVoiceEvent is called upon receiving GuildVoiceJoinEvent, GuildVoiceUpdateEvent, GuildVoiceLeaveEvent
type GenericGuildVoiceEvent struct {
	*GenericGuildMemberEvent
	VoiceState *api.VoiceState
}

// GuildVoiceJoinEvent indicates that a api.Member joined a api.VoiceChannel(requires api.GatewayIntentsGuildVoiceStates)
type GuildVoiceJoinEvent struct {
	*GenericGuildVoiceEvent
}

// GuildVoiceUpdateEvent indicates that a api.Member moved a api.VoiceChannel(requires api.GatewayIntentsGuildVoiceStates)
type GuildVoiceUpdateEvent struct {
	*GenericGuildVoiceEvent
	OldVoiceState *api.VoiceState
}

// GuildVoiceLeaveEvent indicates that a api.Member left a api.VoiceChannel(requires api.GatewayIntentsGuildVoiceStates)
type GuildVoiceLeaveEvent struct {
	*GenericGuildVoiceEvent
}
