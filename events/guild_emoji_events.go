package events

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

type EmojisUpdate struct {
	*GenericEvent
	discord.GatewayEventGuildEmojisUpdate
}

// GenericEmoji is called upon receiving EmojiCreate , EmojiUpdate or EmojiDelete (requires discord.GatewayIntentGuildEmojisAndStickers)
type GenericEmoji struct {
	*GenericEvent
	GuildID snowflake.ID
	Emoji   discord.Emoji
}

// EmojiCreate indicates that a new discord.Emoji got created in a discord.Guild (requires discord.GatewayIntentGuildEmojisAndStickers)
type EmojiCreate struct {
	*GenericEmoji
}

// EmojiUpdate indicates that a discord.Emoji got updated in a discord.Guild (requires discord.GatewayIntentGuildEmojisAndStickers)
type EmojiUpdate struct {
	*GenericEmoji
	OldEmoji discord.Emoji
}

// EmojiDelete indicates that a discord.Emoji got deleted in a discord.Guild (requires discord.GatewayIntentGuildEmojisAndStickers)
type EmojiDelete struct {
	*GenericEmoji
}
