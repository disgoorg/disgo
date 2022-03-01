package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/snowflake"
	"github.com/google/go-cmp/cmp"
)

// gatewayHandlerGuildEmojisUpdate handles discord.GatewayEventTypeGuildEmojisUpdate
type gatewayHandlerGuildEmojisUpdate struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerGuildEmojisUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildEmojisUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildEmojisUpdate) New() interface{} {
	return &discord.GuildEmojisUpdateGatewayEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildEmojisUpdate) HandleGatewayEvent(bot *core.Bot, sequenceNumber discord.GatewaySequence, shardID int, v interface{}) {
	payload := *v.(*discord.GuildEmojisUpdateGatewayEvent)

	if bot.Caches.Config().CacheFlags.Missing(core.CacheFlagEmojis) {
		return
	}

	var (
		emojiCache    = bot.Caches.Emojis().GuildCache(payload.GuildID)
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
		bot.Caches.Emojis().Remove(payload.GuildID, emojiID)
	}

	for _, emoji := range newEmojis {
		bot.EventManager.Dispatch(&events.EmojiCreateEvent{
			GenericEmojiEvent: &events.GenericEmojiEvent{
				GenericEvent: events.NewGenericEvent(bot, sequenceNumber, shardID),
				GuildID:      payload.GuildID,
				Emoji:        emoji,
			},
		})
	}

	for _, emoji := range updatedEmojis {
		bot.EventManager.Dispatch(&events.EmojiUpdateEvent{
			GenericEmojiEvent: &events.GenericEmojiEvent{
				GenericEvent: events.NewGenericEvent(bot, sequenceNumber, shardID),
				GuildID:      payload.GuildID,
				Emoji:        emoji,
			},
		})
	}

	for _, emoji := range oldEmojis {
		bot.EventManager.Dispatch(&events.EmojiDeleteEvent{
			GenericEmojiEvent: &events.GenericEmojiEvent{
				GenericEvent: events.NewGenericEvent(bot, sequenceNumber, shardID),
				GuildID:      payload.GuildID,
				Emoji:        emoji,
			},
		})
	}

}
