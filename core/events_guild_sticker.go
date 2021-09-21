package core

// GenericStickerEvent is called upon receiving StickerCreateEvent, StickerUpdateEvent or StickerDeleteEvent(requires core.GatewayIntentsGuildStickers)
type GenericStickerEvent struct {
	*GenericGuildEvent
	Sticker *Sticker
}

// StickerCreateEvent indicates that a new core.Sticker got created in a core.Guild(requires core.GatewayIntentsGuildStickers)
type StickerCreateEvent struct {
	*GenericStickerEvent
}

// StickerUpdateEvent indicates that a core.Sticker got updated in a core.Guild(requires core.GatewayIntentsGuildStickers)
type StickerUpdateEvent struct {
	*GenericStickerEvent
	OldSticker *Sticker
}

// StickerDeleteEvent indicates that a core.Sticker got deleted in a core.Guild(requires core.GatewayIntentsGuildStickers)
type StickerDeleteEvent struct {
	*GenericStickerEvent
}
