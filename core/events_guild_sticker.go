package core

// GenericStickerEvent is called upon receiving StickerCreateEvent, StickerUpdateEvent or StickerDeleteEvent (requires discord.GatewayIntentGuildEmojisAndStickers)
type GenericStickerEvent struct {
	*GenericGuildEvent
	Sticker *Sticker
}

// StickerCreateEvent indicates that a new core.Sticker got created in a core.Guild (requires discord.GatewayIntentGuildEmojisAndStickers)
type StickerCreateEvent struct {
	*GenericStickerEvent
}

// StickerUpdateEvent indicates that a core.Sticker got updated in a core.Guild (requires discord.GatewayIntentGuildEmojisAndStickers)
type StickerUpdateEvent struct {
	*GenericStickerEvent
	OldSticker *Sticker
}

// StickerDeleteEvent indicates that a core.Sticker got deleted in a core.Guild (requires discord.GatewayIntentGuildEmojisAndStickers)
type StickerDeleteEvent struct {
	*GenericStickerEvent
}
