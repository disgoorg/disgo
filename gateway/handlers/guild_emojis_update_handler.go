package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"golang.org/x/exp/slices"
)

// gatewayHandlerGuildEmojisUpdate handles discord.GatewayEventTypeGuildEmojisUpdate
type gatewayHandlerGuildEmojisUpdate struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerGuildEmojisUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildEmojisUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildEmojisUpdate) New() any {
	return &discord.GatewayEventGuildEmojisUpdate{}
}

type updatedEmoji struct {
	old discord.Emoji
	new discord.Emoji
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildEmojisUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*discord.GatewayEventGuildEmojisUpdate)

	client.EventManager().DispatchEvent(&events.EmojisUpdateEvent{
		GenericEvent:                  events.NewGenericEvent(client, sequenceNumber, shardID),
		GatewayEventGuildEmojisUpdate: payload,
	})

	if client.Caches().CacheFlags().Missing(cache.FlagEmojis) {
		return
	}

	createdEmojis := map[snowflake.ID]discord.Emoji{}
	deletedEmojis := client.Caches().Emojis().MapGroupAll(payload.GuildID)
	updatedEmojis := map[snowflake.ID]updatedEmoji{}

	for _, newEmoji := range payload.Emojis {
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
		client.Caches().Emojis().Put(payload.GuildID, emoji.ID, emoji)
		client.EventManager().DispatchEvent(&events.EmojiCreateEvent{
			GenericEmojiEvent: &events.GenericEmojiEvent{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
				GuildID:      payload.GuildID,
				Emoji:        emoji,
			},
		})
	}

	for _, emoji := range updatedEmojis {
		client.Caches().Emojis().Put(payload.GuildID, emoji.new.ID, emoji.new)
		client.EventManager().DispatchEvent(&events.EmojiUpdateEvent{
			GenericEmojiEvent: &events.GenericEmojiEvent{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
				GuildID:      payload.GuildID,
				Emoji:        emoji.new,
			},
			OldEmoji: emoji.old,
		})
	}

	for _, emoji := range deletedEmojis {
		client.Caches().Emojis().Remove(payload.GuildID, emoji.ID)
		client.EventManager().DispatchEvent(&events.EmojiDeleteEvent{
			GenericEmojiEvent: &events.GenericEmojiEvent{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
				GuildID:      payload.GuildID,
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
