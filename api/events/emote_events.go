package events

import (
	"github.com/DisgoOrg/disgo/api"
)

type GenericEmoteEvent struct {
	GenericGuildEvent
	EmoteID api.Snowflake
}

func (e GenericEmoteEvent) Emote() *api.Emote {
	return e.Disgo().Cache().Emote(e.EmoteID)
}

type EmoteCreateEvent struct {
	GenericEmoteEvent
	Emote *api.Emote
}

type EmoteUpdateEvent struct {
	GenericEmoteEvent
	NewEmote *api.Emote
	OldEmote *api.Emote
}

type EmoteDeleteEvent struct {
	GenericEmoteEvent
	Emote *api.Emote
}
