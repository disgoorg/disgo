package core

import (
	"github.com/DisgoOrg/disgo/discord"
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
func (h *gatewayHandlerGuildStickersUpdate) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.GuildStickersUpdateGatewayEvent)

	if bot.Caches.Config().CacheFlags.Missing(CacheFlagStickers) {
		return
	}

	var (
		stickerCache    = bot.Caches.StickerCache().GuildCache(payload.GuildID)
		oldStickers     map[discord.Snowflake]*Sticker
		newStickers     map[discord.Snowflake]*Sticker
		updatedStickers map[discord.Snowflake]*Sticker
	)

	oldStickers = make(map[discord.Snowflake]*Sticker, len(stickerCache))
	for key, value := range stickerCache {
		va := *value
		oldStickers[key] = &va
	}

	for _, current := range payload.Stickers {
		sticker, ok := stickerCache[current.ID]
		if ok {
			delete(oldStickers, current.ID)
			if !cmp.Equal(sticker, current) {
				updatedStickers[current.ID] = bot.EntityBuilder.CreateSticker(current, CacheStrategyYes)
			}
		} else {
			newStickers[current.ID] = bot.EntityBuilder.CreateSticker(current, CacheStrategyYes)
		}
	}

	for stickerID := range oldStickers {
		bot.Caches.StickerCache().Remove(payload.GuildID, stickerID)
	}

	genericGuildEvent := &GenericGuildEvent{
		GenericEvent: NewGenericEvent(bot, sequenceNumber),
		GuildID:      payload.GuildID,
		Guild:        bot.Caches.GuildCache().Get(payload.GuildID),
	}

	for _, sticker := range newStickers {
		bot.EventManager.Dispatch(&StickerCreateEvent{
			GenericStickerEvent: &GenericStickerEvent{
				GenericGuildEvent: genericGuildEvent,
				Sticker:           sticker,
			},
		})
	}

	for _, sticker := range updatedStickers {
		bot.EventManager.Dispatch(&StickerUpdateEvent{
			GenericStickerEvent: &GenericStickerEvent{
				GenericGuildEvent: genericGuildEvent,
				Sticker:           sticker,
			},
		})
	}

	for _, sticker := range oldStickers {
		bot.EventManager.Dispatch(&StickerDeleteEvent{
			GenericStickerEvent: &GenericStickerEvent{
				GenericGuildEvent: genericGuildEvent,
				Sticker:           sticker,
			},
		})
	}

}
