package events

import (
	"github.com/DisgoOrg/disgo/api"
)

// GenericEmoteEvent is called upon receiving EmoteCreateEvent, EmoteUpdateEvent or EmoteDeleteEvent(requires api.IntentsGuildEmojis)
type GenericEmoteEvent struct {
	GenericGuildEvent
	Emote *api.Emote
}

// EmoteCreateEvent indicates that a new api.Emote got created in a api.Guild(requires api.IntentsGuildEmojis)
type EmoteCreateEvent struct {
	GenericEmoteEvent
}

// EmoteUpdateEvent indicates that a api.Emote got updated in a api.Guild(requires api.IntentsGuildEmojis)
type EmoteUpdateEvent struct {
	GenericEmoteEvent
	OldEmote *api.Emote
}

// EmoteDeleteEvent indicates that a api.Emote got deleted in a api.Guild(requires api.IntentsGuildEmojis)
type EmoteDeleteEvent struct {
	GenericEmoteEvent
}
