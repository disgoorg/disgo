package handlers

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

type updatedSticker struct {
	old discord.Sticker
	new discord.Sticker
}

func gatewayHandlerGuildStickersUpdate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildStickersUpdate) {
	client.EventManager().DispatchEvent(&events.StickersUpdate{
		GenericEvent:             events.NewGenericEvent(client, sequenceNumber, shardID),
		EventGuildStickersUpdate: event,
	})

	if client.Caches().CacheFlags().Missing(cache.FlagStickers) {
		return
	}

	createdStickers := map[snowflake.ID]discord.Sticker{}
	deletedStickers := map[snowflake.ID]discord.Sticker{}
	updatedStickers := map[snowflake.ID]updatedSticker{}

	client.Caches().StickersForEach(event.GuildID, func(sticker discord.Sticker) {
		deletedStickers[sticker.ID] = sticker
	})

	for _, newSticker := range event.Stickers {
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
		client.EventManager().DispatchEvent(&events.StickerCreate{
			GenericSticker: &events.GenericSticker{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
				GuildID:      event.GuildID,
				Sticker:      emoji,
			},
		})
	}

	for _, emoji := range updatedStickers {
		client.EventManager().DispatchEvent(&events.StickerUpdate{
			GenericSticker: &events.GenericSticker{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
				GuildID:      event.GuildID,
				Sticker:      emoji.new,
			},
			OldSticker: emoji.old,
		})
	}

	for _, emoji := range deletedStickers {
		client.EventManager().DispatchEvent(&events.StickerDelete{
			GenericSticker: &events.GenericSticker{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
				GuildID:      event.GuildID,
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
