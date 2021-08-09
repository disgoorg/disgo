package events

import (
	
)

// GenericEmoteEvent is called upon receiving EmoteCreateEvent, EmoteUpdateEvent or EmoteDeleteEvent(requires api.GatewayIntentsGuildEmojis)
type GenericEmoteEvent struct {
	*GenericGuildEvent
	Emote *core.Emoji
}

// EmoteCreateEvent indicates that a new api.Emoji got created in an api.Guild(requires api.GatewayIntentsGuildEmojis)
type EmoteCreateEvent struct {
	*GenericEmoteEvent
}

// EmoteUpdateEvent indicates that an api.Emoji got updated in an api.Guild(requires api.GatewayIntentsGuildEmojis)
type EmoteUpdateEvent struct {
	*GenericEmoteEvent
	OldEmote *core.Emoji
}

// EmoteDeleteEvent indicates that an api.Emoji got deleted in an api.Guild(requires api.GatewayIntentsGuildEmojis)
type EmoteDeleteEvent struct {
	*GenericEmoteEvent
}
