package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake"
)

// gatewayHandlerGuildStickersUpdate handles discord.GatewayEventTypeGuildStickersUpdate
type gatewayHandlerGuildStickersUpdate struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerGuildStickersUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildStickersUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildStickersUpdate) New() any {
	return &discord.GatewayEventGuildStickersUpdate{}
}

type updatedSticker struct {
	old discord.Sticker
	new discord.Sticker
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildStickersUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber int, v any) {
	payload := *v.(*discord.GatewayEventGuildStickersUpdate)

	client.EventManager().DispatchEvent(&events.StickersUpdateEvent{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber),
		GuildID:      payload.GuildID,
		Stickers:     payload.Stickers,
	})

	if client.Caches().CacheFlags().Missing(cache.FlagStickers) {
		return
	}

	createdStickers := map[snowflake.Snowflake]discord.Sticker{}
	deletedStickers := client.Caches().Stickers().MapGroupAll(payload.GuildID)
	updatedStickers := map[snowflake.Snowflake]updatedSticker{}

	for _, newSticker := range payload.Stickers {
		oldSticker, ok := deletedStickers[newSticker.ID]
		if ok {
			delete(deletedStickers, newSticker.ID)
			if isStickerUpdated(oldSticker, newSticker) {
				updatedStickers[newSticker.ID] = updatedSticker{new: newSticker, old: oldSticker}
			}
			continue
		}
		createdStickers[newSticker.ID] = newSticker
	}

	for _, emoji := range createdStickers {
		client.EventManager().DispatchEvent(&events.StickerCreateEvent{
			GenericStickerEvent: &events.GenericStickerEvent{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber),
				GuildID:      payload.GuildID,
				Sticker:      emoji,
			},
		})
	}

	for _, emoji := range updatedStickers {
		client.EventManager().DispatchEvent(&events.StickerUpdateEvent{
			GenericStickerEvent: &events.GenericStickerEvent{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber),
				GuildID:      payload.GuildID,
				Sticker:      emoji.new,
			},
			OldSticker: emoji.old,
		})
	}

	for _, emoji := range deletedStickers {
		client.EventManager().DispatchEvent(&events.StickerDeleteEvent{
			GenericStickerEvent: &events.GenericStickerEvent{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber),
				GuildID:      payload.GuildID,
				Sticker:      emoji,
			},
		})
	}
}

func isStickerUpdated(old discord.Sticker, new discord.Sticker) bool {
	if old.Name != new.Name {
		return true
	}
	if old.Description != new.Description {
		return true
	}
	if old.Tags != new.Tags {
		return true
	}
	return false
}
