package events

import "github.com/DisgoOrg/disgo/api"

type GenericGuildVoiceEvent struct {
	GenericGuildMemberEvent
	Member     *api.Member
	VoiceState *api.VoiceState
}

// GuildMemberVoiceStateUpdateEvent indicates that the api.VoiceState of a api.Member updated
type GuildMemberVoiceStateUpdateEvent struct {
	GenericGuildVoiceEvent
	Left   *api.VoiceChannel
	Joined *api.VoiceChannel
}

type GuildMemberVoiceDeafenEvent struct {
	GenericGuildVoiceEvent
}

func (e GuildMemberVoiceDeafenEvent) Deafened() bool {
	return e.VoiceState.Deafened()
}
