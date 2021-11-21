package gatewayhandlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
	"github.com/google/go-cmp/cmp"
)

// gatewayHandlerGuildStickersUpdate handles discord.GatewayEventTypeGuildStickersUpdate
type gatewayHandlerGuildStickersUpdate struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *gatewayHandlerGuildStickersUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildStickersUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildStickersUpdate) New() interface{} {
	return &discord.GuildStickersUpdateGatewayEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildStickersUpdate) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.GuildStickersUpdateGatewayEvent)

	if bot.Caches.Config().CacheFlags.Missing(core.CacheFlagStickers) {
		return
	}

	var (
		stickerCache    = bot.Caches.StickerCache().GuildCache(payload.GuildID)
		oldStickers     = map[discord.Snowflake]*core.Sticker{}
		newStickers     = map[discord.Snowflake]*core.Sticker{}
		updatedStickers = map[discord.Snowflake]*core.Sticker{}
	)

	oldStickers = make(map[discord.Snowflake]*core.Sticker, len(stickerCache))
	for key, value := range stickerCache {
		va := *value
		oldStickers[key] = &va
	}

	for _, current := range payload.Stickers {
		sticker, ok := stickerCache[current.ID]
		if ok {
			delete(oldStickers, current.ID)
			if !cmp.Equal(sticker, current) {
				updatedStickers[current.ID] = bot.EntityBuilder.CreateSticker(current, core.CacheStrategyYes)
			}
		} else {
			newStickers[current.ID] = bot.EntityBuilder.CreateSticker(current, core.CacheStrategyYes)
		}
	}

	for stickerID := range oldStickers {
		bot.Caches.StickerCache().Remove(payload.GuildID, stickerID)
	}

	for _, sticker := range newStickers {
		bot.EventManager.Dispatch(&events.StickerCreateEvent{
			GenericStickerEvent: &events.GenericStickerEvent{
				GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
				GuildID:      payload.GuildID,
				Sticker:      sticker,
			},
		})
	}

	for _, sticker := range updatedStickers {
		bot.EventManager.Dispatch(&events.StickerUpdateEvent{
			GenericStickerEvent: &events.GenericStickerEvent{
				GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
				GuildID:      payload.GuildID,
				Sticker:      sticker,
			},
		})
	}

	for _, sticker := range oldStickers {
		bot.EventManager.Dispatch(&events.StickerDeleteEvent{
			GenericStickerEvent: &events.GenericStickerEvent{
				GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
				GuildID:      payload.GuildID,
				Sticker:      sticker,
			},
		})
	}

}
