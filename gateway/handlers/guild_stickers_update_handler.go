package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
)

// gatewayHandlerGuildStickersUpdate handles discord.GatewayEventTypeGuildStickersUpdate
type gatewayHandlerGuildStickersUpdate struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerGuildStickersUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildStickersUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildStickersUpdate) New() any {
	return &discord.GuildStickersUpdateGatewayEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildStickersUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber discord.GatewaySequence, v any) {
	/*payload := *v.(*discord.GuildStickersUpdateGatewayEvent)

	if client.Caches().Config().CacheFlags.Missing(cache.CacheFlagStickers) {
		return
	}

	var (
		stickerCache    = client.Caches().Stickers().GuildCache(payload.GuildID)
		oldStickers     = map[snowflake.Snowflake]*core.Sticker{}
		newStickers     = map[snowflake.Snowflake]*core.Sticker{}
		updatedStickers = map[snowflake.Snowflake]*core.Sticker{}
	)

	oldStickers = make(map[snowflake.Snowflake]*core.Sticker, len(stickerCache))
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
		client.Caches().Stickers().Remove(payload.GuildID, stickerID)
	}

	for _, sticker := range newStickers {
		client.EventManager().Dispatch(&events.StickerCreateEvent{
			GenericStickerEvent: &events.GenericStickerEvent{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber),
				GuildID:      payload.GuildID,
				Sticker:      sticker,
			},
		})
	}

	for _, sticker := range updatedStickers {
		client.EventManager().Dispatch(&events.StickerUpdateEvent{
			GenericStickerEvent: &events.GenericStickerEvent{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber),
				GuildID:      payload.GuildID,
				Sticker:      sticker,
			},
		})
	}

	for _, sticker := range oldStickers {
		client.EventManager().Dispatch(&events.StickerDeleteEvent{
			GenericStickerEvent: &events.GenericStickerEvent{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber),
				GuildID:      payload.GuildID,
				Sticker:      sticker,
			},
		})
	}*/
}
