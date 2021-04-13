package events

import (
	"github.com/DisgoOrg/disgo/api"
)

type GenericEmoteEvent struct {
	GenericGuildEvent
	Emote *api.Emote
}

type EmoteCreateEvent struct {
	GenericEmoteEvent
}

type EmoteUpdateEvent struct {
	GenericEmoteEvent
	OldEmote *api.Emote
}

type EmoteDeleteEvent struct {
	GenericEmoteEvent
}
