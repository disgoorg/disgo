package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
	"golang.org/x/exp/slices"
)

type updatedEmoji struct {
	old discord.Emoji
	new discord.Emoji
}

func gatewayHandlerGuildEmojisUpdate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildEmojisUpdate) {
	client.EventManager().DispatchEvent(&events.EmojisUpdate{
		GenericEvent:           events.NewGenericEvent(client, sequenceNumber, shardID),
		EventGuildEmojisUpdate: event,
	})

	if client.Caches().CacheFlags().Missing(cache.FlagEmojis) {
		return
	}

	createdEmojis := map[snowflake.ID]discord.Emoji{}
	deletedEmojis := client.Caches().Emojis().MapGroupAll(event.GuildID)
	updatedEmojis := map[snowflake.ID]updatedEmoji{}

	for _, newEmoji := range event.Emojis {
		oldEmoji, ok := deletedEmojis[newEmoji.ID]
		if ok {
			delete(deletedEmojis, newEmoji.ID)
			if isEmojiUpdated(oldEmoji, newEmoji) {
				updatedEmojis[newEmoji.ID] = updatedEmoji{new: newEmoji, old: oldEmoji}
			}
			continue
		}
		createdEmojis[newEmoji.ID] = newEmoji
	}

	for _, emoji := range createdEmojis {
		client.Caches().Emojis().Put(event.GuildID, emoji.ID, emoji)
		client.EventManager().DispatchEvent(&events.EmojiCreate{
			GenericEmoji: &events.GenericEmoji{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
				GuildID:      event.GuildID,
				Emoji:        emoji,
			},
		})
	}

	for _, emoji := range updatedEmojis {
		client.Caches().Emojis().Put(event.GuildID, emoji.new.ID, emoji.new)
		client.EventManager().DispatchEvent(&events.EmojiUpdate{
			GenericEmoji: &events.GenericEmoji{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
				GuildID:      event.GuildID,
				Emoji:        emoji.new,
			},
			OldEmoji: emoji.old,
		})
	}

	for _, emoji := range deletedEmojis {
		client.Caches().Emojis().Remove(event.GuildID, emoji.ID)
		client.EventManager().DispatchEvent(&events.EmojiDelete{
			GenericEmoji: &events.GenericEmoji{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
				GuildID:      event.GuildID,
				Emoji:        emoji,
			},
		})
	}

}

func isEmojiUpdated(old discord.Emoji, new discord.Emoji) bool {
	if old.Name != new.Name {
		return true
	}
	if !slices.Equal(old.Roles, new.Roles) {
		return true
	}
	return false
}
