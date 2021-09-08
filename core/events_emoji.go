package core

// GenericEmojiEvent is called upon receiving EmojiCreateEvent, EmojiUpdateEvent or EmojiDeleteEvent(requires api.GatewayIntentsGuildEmojis)
type GenericEmojiEvent struct {
	*GenericGuildEvent
	Emoji *Emoji
}

// EmojiCreateEvent indicates that a new api.Emoji got created in an api.Guild(requires api.GatewayIntentsGuildEmojis)
type EmojiCreateEvent struct {
	*GenericEmojiEvent
}

// EmojiUpdateEvent indicates that an api.Emoji got updated in an api.Guild(requires api.GatewayIntentsGuildEmojis)
type EmojiUpdateEvent struct {
	*GenericEmojiEvent
	OldEmoji *Emoji
}

// EmojiDeleteEvent indicates that an api.Emoji got deleted in an api.Guild(requires api.GatewayIntentsGuildEmojis)
type EmojiDeleteEvent struct {
	*GenericEmojiEvent
}
