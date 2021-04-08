package events

import (
	"github.com/DisgoOrg/disgo/api"
)

type GenericGuildVoiceEvent struct {
	GenericGuildMemberEvent
	Member *api.Member
}

type GuildVoiceUpdateEvent struct {
	GenericGuildVoiceEvent
	NewVoiceState *api.VoiceState
	OldVoiceState *api.VoiceState
}

type GuildVoiceJoinEvent struct {
	GenericGuildVoiceEvent
	GenericVoiceChannelEvent
}

type GuildVoiceLeaveEvent struct {
	GenericGuildVoiceEvent
	GenericVoiceChannelEvent
}
