package events

import (
	"github.com/DisgoOrg/disgo/api"
)

// GenericEmoteEvent is called upon receiving EmoteCreateEvent, EmoteUpdateEvent or EmoteDeleteEvent(requires api.GatewayIntentsGuildEmojis)
type GenericEmoteEvent struct {
	*GenericGuildEvent
	Emote *api.Emoji
}

// EmoteCreateEvent indicates that a new api.Emoji got created in an api.Guild(requires api.GatewayIntentsGuildEmojis)
type EmoteCreateEvent struct {
	*GenericEmoteEvent
}

// EmoteUpdateEvent indicates that an api.Emoji got updated in an api.Guild(requires api.GatewayIntentsGuildEmojis)
type EmoteUpdateEvent struct {
	*GenericEmoteEvent
	OldEmote *api.Emoji
}

// EmoteDeleteEvent indicates that an api.Emoji got deleted in an api.Guild(requires api.GatewayIntentsGuildEmojis)
type EmoteDeleteEvent struct {
	*GenericEmoteEvent
}
