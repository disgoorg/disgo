package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

type guildEmojisUpdatePayload struct {
	GuildID discord.Snowflake `json:"guild_id"`
	Emojis  []discord.Emoji   `json:"emojis"`
}

// GuildEmojisUpdateHandler handles discord.GatewayEventTypeGuildEmojisUpdate
type GuildEmojisUpdateHandler struct{}

// EventType returns the api.GatewayGatewayEventType
func (h *GuildEmojisUpdateHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildEmojisUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *GuildEmojisUpdateHandler) New() interface{} {
	return guildEmojisUpdatePayload{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *GuildEmojisUpdateHandler) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	payload, ok := v.(guildEmojisUpdatePayload)
	if !ok {
		return
	}

	if bot.Caches.Config().CacheFlags.Missing(core.CacheFlagEmojis) {
		return
	}

	var (
		emojiCache    = bot.Caches.EmojiCache().GuildCache(payload.GuildID)
		oldEmojis     map[discord.Snowflake]*core.Emoji
		newEmojis     map[discord.Snowflake]*core.Emoji
		updatedEmojis map[discord.Snowflake]*core.Emoji
	)

	oldEmojis = make(map[discord.Snowflake]*core.Emoji, len(emojiCache))
	for key, value := range emojiCache {
		oldEmojis[key] = &*value
	}

	for _, current := range payload.Emojis {
		emoji, ok := emojiCache[current.ID]
		if ok {
			delete(oldEmojis, current.ID)
			if isEmojiUpdated(emoji, current) {
				updatedEmojis[current.ID] = bot.EntityBuilder.CreateEmoji(payload.GuildID, current, core.CacheStrategyYes)
			}
		} else {
			newEmojis[current.ID] = bot.EntityBuilder.CreateEmoji(payload.GuildID, current, core.CacheStrategyYes)
		}
	}

	for emojiID := range oldEmojis {
		bot.Caches.EmojiCache().Remove(payload.GuildID, emojiID)
	}

	genericGuildEvent := &events.GenericGuildEvent{
		GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
		GuildID:      payload.GuildID,
		Guild:        bot.Caches.GuildCache().Get(payload.GuildID),
	}

	for _, emoji := range newEmojis {
		bot.EventManager.Dispatch(&events.EmojiCreateEvent{
			GenericEmojiEvent: &events.GenericEmojiEvent{
				GenericGuildEvent: genericGuildEvent,
				Emoji:             emoji,
			},
		})
	}

	for _, emoji := range updatedEmojis {
		bot.EventManager.Dispatch(&events.EmojiUpdateEvent{
			GenericEmojiEvent: &events.GenericEmojiEvent{
				GenericGuildEvent: genericGuildEvent,
				Emoji:             emoji,
			},
		})
	}

	for _, emoji := range oldEmojis {
		bot.EventManager.Dispatch(&events.EmojiDeleteEvent{
			GenericEmojiEvent: &events.GenericEmojiEvent{
				GenericGuildEvent: genericGuildEvent,
				Emoji:             emoji,
			},
		})
	}

}

func isEmojiUpdated(oldEmoji *core.Emoji, newEmoji discord.Emoji) bool {
	// TODO: actual check here
	return false
}
