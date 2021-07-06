package events

import (
	"github.com/DisgoOrg/disgo/api"
)

// GenericVoiceChannelEvent is called upon receiving VoiceChannelCreateEvent, VoiceChannelUpdateEvent or VoiceChannelDeleteEvent
type GenericVoiceChannelEvent struct {
	*GenericGuildChannelEvent
	VoiceChannel *api.VoiceChannel
}

// VoiceChannelCreateEvent indicates that a new api.VoiceChannel got created in a api.Guild
type VoiceChannelCreateEvent struct {
	*GenericVoiceChannelEvent
}

// VoiceChannelUpdateEvent indicates that a api.VoiceChannel got updated in a api.Guild
type VoiceChannelUpdateEvent struct {
	*GenericVoiceChannelEvent
	OldVoiceChannel *api.VoiceChannel
}

// VoiceChannelDeleteEvent indicates that a api.VoiceChannel got deleted in a api.Guild
type VoiceChannelDeleteEvent struct {
	*GenericVoiceChannelEvent
}
