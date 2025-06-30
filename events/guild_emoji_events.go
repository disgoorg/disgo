package events

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
)

// EmojisUpdate is dispatched when a guild's emojis are updated.
// This event does not depend on a cache like EmojiCreate, EmojiUpdate or EmojiDelete.
type EmojisUpdate struct {
	*GenericEvent
	gateway.EventGuildEmojisUpdate
}

// GenericEmoji is called upon receiving EmojiCreate , EmojiUpdate or EmojiDelete (requires gateway.IntentGuildExpressions)
type GenericEmoji struct {
	*GenericEvent
	GuildID snowflake.ID
	Emoji   discord.Emoji
}

// EmojiCreate indicates that a new discord.Emoji got created in a discord.Guild (requires gateway.IntentGuildExpressions)
type EmojiCreate struct {
	*GenericEmoji
}

// EmojiUpdate indicates that a discord.Emoji got updated in a discord.Guild (requires gateway.IntentGuildExpressions)
type EmojiUpdate struct {
	*GenericEmoji
	OldEmoji discord.Emoji
}

// EmojiDelete indicates that a discord.Emoji got deleted in a discord.Guild (requires gateway.IntentGuildExpressions)
type EmojiDelete struct {
	*GenericEmoji
}
