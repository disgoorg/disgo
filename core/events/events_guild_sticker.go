package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/snowflake"
)

// GenericStickerEvent is called upon receiving StickerCreateEvent, StickerUpdateEvent or StickerDeleteEvent (requires discord.GatewayIntentGuildEmojisAndStickers)
type GenericStickerEvent struct {
	*GenericEvent
	GuildID snowflake.Snowflake
	Sticker *core.Sticker
}

// StickerCreateEvent indicates that a new core.Sticker got created in a core.Guild (requires discord.GatewayIntentGuildEmojisAndStickers)
type StickerCreateEvent struct {
	*GenericStickerEvent
}

// StickerUpdateEvent indicates that a core.Sticker got updated in a core.Guild (requires discord.GatewayIntentGuildEmojisAndStickers)
type StickerUpdateEvent struct {
	*GenericStickerEvent
	OldSticker *core.Sticker
}

// StickerDeleteEvent indicates that a core.Sticker got deleted in a core.Guild (requires discord.GatewayIntentGuildEmojisAndStickers)
type StickerDeleteEvent struct {
	*GenericStickerEvent
}
