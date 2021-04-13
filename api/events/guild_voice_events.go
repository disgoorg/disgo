package events

import (
	"github.com/DisgoOrg/disgo/api"
)

type GenericGuildVoiceEvent struct {
	GenericGuildMemberEvent
	VoiceState *api.VoiceState
}

type GuildVoiceJoinEvent struct {
	GenericGuildVoiceEvent
	GenericVoiceChannelEvent
}

type GuildVoiceUpdateEvent struct {
	GenericGuildVoiceEvent
	OldVoiceState *api.VoiceState
}

type GuildVoiceLeaveEvent struct {
	GenericGuildVoiceEvent
	GenericVoiceChannelEvent
}
