package events

import (
	"github.com/DisgoOrg/disgo/api"
)

type GenericVoiceChannelEvent struct {
	GenericChannelEvent
	VoiceChannel *api.VoiceChannel
}

type VoiceChannelCreateEvent struct {
	GenericVoiceChannelEvent
}

type VoiceChannelUpdateEvent struct {
	GenericVoiceChannelEvent
	OldVoiceChannel *api.VoiceChannel
}

type VoiceChannelDeleteEvent struct {
	GenericVoiceChannelEvent
}
