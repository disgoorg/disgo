package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/snowflake"
)

// GenericEmojiEvent is called upon receiving EmojiCreateEvent, EmojiUpdateEvent or EmojiDeleteEvent (requires discord.GatewayIntentGuildEmojisAndStickers)
type GenericEmojiEvent struct {
	*GenericEvent
	GuildID snowflake.Snowflake
	Emoji   *core.Emoji
}

// EmojiCreateEvent indicates that a new core.Emoji got created in a core.Guild (requires discord.GatewayIntentGuildEmojisAndStickers)
type EmojiCreateEvent struct {
	*GenericEmojiEvent
}

// EmojiUpdateEvent indicates that a core.Emoji got updated in a core.Guild (requires discord.GatewayIntentGuildEmojisAndStickers)
type EmojiUpdateEvent struct {
	*GenericEmojiEvent
	OldEmoji *core.Emoji
}

// EmojiDeleteEvent indicates that a core.Emoji got deleted in a core.Guild (requires discord.GatewayIntentGuildEmojisAndStickers)
type EmojiDeleteEvent struct {
	*GenericEmojiEvent
}