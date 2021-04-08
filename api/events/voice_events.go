package events

import (
	"github.com/DisgoOrg/disgo/api"
)

type GenericVoiceChannelEvent struct {
	GenericChannelEvent
}

func (e GenericVoiceChannelEvent) VoiceChannel() *api.VoiceChannel {
	return e.Disgo().Cache().VoiceChannel(e.ChannelID)
}

type VoiceChannelCreateEvent struct {
	GenericVoiceChannelEvent
	VoiceChannel *api.VoiceChannel
}

type VoiceChannelUpdateEvent struct {
	GenericVoiceChannelEvent
	NewVoiceChannel *api.VoiceChannel
	OldVoiceChannel *api.VoiceChannel
}

type VoiceChannelDeleteEvent struct {
	GenericVoiceChannelEvent
	VoiceChannel *api.VoiceChannel
}
