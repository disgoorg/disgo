package events

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

type StickersUpdateEvent struct {
	*GenericEvent
	GuildID  snowflake.ID
	Stickers []discord.Sticker
}

// GenericStickerEvent is called upon receiving StickerCreateEvent, StickerUpdateEvent or StickerDeleteEvent (requires discord.GatewayIntentGuildEmojisAndStickers)
type GenericStickerEvent struct {
	*GenericEvent
	GuildID snowflake.ID
	Sticker discord.Sticker
}

// StickerCreateEvent indicates that a new discord.Sticker got created in a discord.Guild (requires discord.GatewayIntentGuildEmojisAndStickers)
type StickerCreateEvent struct {
	*GenericStickerEvent
}

// StickerUpdateEvent indicates that a discord.Sticker got updated in a discord.Guild (requires discord.GatewayIntentGuildEmojisAndStickers)
type StickerUpdateEvent struct {
	*GenericStickerEvent
	OldSticker discord.Sticker
}

// StickerDeleteEvent indicates that a discord.Sticker got deleted in a discord.Guild (requires discord.GatewayIntentGuildEmojisAndStickers)
type StickerDeleteEvent struct {
	*GenericStickerEvent
}
