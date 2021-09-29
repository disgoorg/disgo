package core

import (
	"github.com/DisgoOrg/disgo/discord"
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
func (h *gatewayHandlerGuildEmojisUpdate) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.GuildEmojisUpdateGatewayEvent)

	if bot.Caches.Config().CacheFlags.Missing(CacheFlagEmojis) {
		return
	}

	var (
		emojiCache    = bot.Caches.EmojiCache().GuildCache(payload.GuildID)
		oldEmojis     map[discord.Snowflake]*Emoji
		newEmojis     map[discord.Snowflake]*Emoji
		updatedEmojis map[discord.Snowflake]*Emoji
	)

	oldEmojis = make(map[discord.Snowflake]*Emoji, len(emojiCache))
	for key, value := range emojiCache {
		va := *value
		oldEmojis[key] = &va
	}

	for _, current := range payload.Emojis {
		emoji, ok := emojiCache[current.ID]
		if ok {
			delete(oldEmojis, current.ID)
			if !cmp.Equal(emoji, current) {
				updatedEmojis[current.ID] = bot.EntityBuilder.CreateEmoji(payload.GuildID, current, CacheStrategyYes)
			}
		} else {
			newEmojis[current.ID] = bot.EntityBuilder.CreateEmoji(payload.GuildID, current, CacheStrategyYes)
		}
	}

	for emojiID := range oldEmojis {
		bot.Caches.EmojiCache().Remove(payload.GuildID, emojiID)
	}

	genericGuildEvent := &GenericGuildEvent{
		GenericEvent: NewGenericEvent(bot, sequenceNumber),
		GuildID:      payload.GuildID,
		Guild:        bot.Caches.GuildCache().Get(payload.GuildID),
	}

	for _, emoji := range newEmojis {
		bot.EventManager.Dispatch(&EmojiCreateEvent{
			GenericEmojiEvent: &GenericEmojiEvent{
				GenericGuildEvent: genericGuildEvent,
				Emoji:             emoji,
			},
		})
	}

	for _, emoji := range updatedEmojis {
		bot.EventManager.Dispatch(&EmojiUpdateEvent{
			GenericEmojiEvent: &GenericEmojiEvent{
				GenericGuildEvent: genericGuildEvent,
				Emoji:             emoji,
			},
		})
	}

	for _, emoji := range oldEmojis {
		bot.EventManager.Dispatch(&EmojiDeleteEvent{
			GenericEmojiEvent: &GenericEmojiEvent{
				GenericGuildEvent: genericGuildEvent,
				Emoji:             emoji,
			},
		})
	}

}
