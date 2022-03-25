package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
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

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildEmojisUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber int, v any) {
	/*payload := *v.(*discord.GatewayEventGuildEmojisUpdate)

	if client.Caches().Config().CacheFlags.Missing(cache.CacheFlagEmojis) {
		return
	}

	var (
		emojiCache    = client.Caches().Emojis().GuildCache(payload.GuildID)
		oldEmojis     = map[snowflake.Snowflake]*core.Emoji{}
		newEmojis     = map[snowflake.Snowflake]*core.Emoji{}
		updatedEmojis = map[snowflake.Snowflake]*core.Emoji{}
	)

	oldEmojis = make(map[snowflake.Snowflake]*core.Emoji, len(emojiCache))
	for key, value := range emojiCache {
		va := *value
		oldEmojis[key] = &va
	}

	for _, current := range payload.Emojis {
		emoji, ok := emojiCache[current.ID]
		if ok {
			delete(oldEmojis, current.ID)
			if !cmp.Equal(emoji, current) {
				updatedEmojis[current.ID] = bot.EntityBuilder.CreateEmoji(payload.GuildID, current, core.CacheStrategyYes)
			}
		} else {
			newEmojis[current.ID] = bot.EntityBuilder.CreateEmoji(payload.GuildID, current, core.CacheStrategyYes)
		}
	}

	for emojiID := range oldEmojis {
		client.Caches().Emojis().Remove(payload.GuildID, emojiID)
	}

	for _, emoji := range newEmojis {
		client.EventManager().DispatchEvent(&events.EmojiCreateEvent{
			GenericEmojiEvent: &events.GenericEmojiEvent{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber),
				GuildID:      payload.GuildID,
				Emoji:        emoji,
			},
		})
	}

	for _, emoji := range updatedEmojis {
		client.EventManager().DispatchEvent(&events.EmojiUpdateEvent{
			GenericEmojiEvent: &events.GenericEmojiEvent{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber),
				GuildID:      payload.GuildID,
				Emoji:        emoji,
			},
		})
	}

	for _, emoji := range oldEmojis {
		client.EventManager().DispatchEvent(&events.EmojiDeleteEvent{
			GenericEmojiEvent: &events.GenericEmojiEvent{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber),
				GuildID:      payload.GuildID,
				Emoji:        emoji,
			},
		})
	}*/
}
