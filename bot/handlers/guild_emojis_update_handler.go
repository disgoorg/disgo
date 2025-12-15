package handlers

import (
	"log/slog"
	"slices"

	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

type updatedEmoji struct {
	old discord.Emoji
	new discord.Emoji
}

func gatewayHandlerGuildEmojisUpdate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildEmojisUpdate) {
	client.EventManager.DispatchEvent(&events.EmojisUpdate{
		GenericEvent:           events.NewGenericEvent(client, sequenceNumber, shardID),
		EventGuildEmojisUpdate: event,
	})

	if client.Caches.CacheFlags().Missing(cache.FlagEmojis) {
		return
	}

	createdEmojis := map[snowflake.ID]discord.Emoji{}
	deletedEmojis := map[snowflake.ID]discord.Emoji{}
	updatedEmojis := map[snowflake.ID]updatedEmoji{}

	emojis, err := client.Caches.Emojis(event.GuildID)
	if err == nil {
		for emoji := range emojis {
			deletedEmojis[emoji.ID] = emoji
		}
	}

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
		if err := client.Caches.AddEmoji(emoji); err != nil {
			client.Logger.Error("failed to add emoji to cache", slog.Any("err", err), slog.String("emoji_id", emoji.ID.String()))
		}
		client.EventManager.DispatchEvent(&events.EmojiCreate{
			GenericEmoji: &events.GenericEmoji{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
				GuildID:      event.GuildID,
				Emoji:        emoji,
			},
		})
	}

	for _, emoji := range updatedEmojis {
		if err := client.Caches.AddEmoji(emoji.new); err != nil {
			client.Logger.Error("failed to add emoji to cache", slog.Any("err", err), slog.String("emoji_id", emoji.new.ID.String()))
		}
		client.EventManager.DispatchEvent(&events.EmojiUpdate{
			GenericEmoji: &events.GenericEmoji{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
				GuildID:      event.GuildID,
				Emoji:        emoji.new,
			},
			OldEmoji: emoji.old,
		})
	}

	for _, emoji := range deletedEmojis {
		if _, err := client.Caches.RemoveEmoji(event.GuildID, emoji.ID); err != nil {
			client.Logger.Error("failed to remove emoji from cache", slog.Any("err", err), slog.String("emoji_id", emoji.ID.String()))
		}
		client.EventManager.DispatchEvent(&events.EmojiDelete{
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
