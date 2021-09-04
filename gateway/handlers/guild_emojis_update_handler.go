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
func (h *GuildEmojisUpdateHandler) HandleGatewayEvent(disgo core.Disgo, eventManager core.EventManager, sequenceNumber int, v interface{}) {
	payload, ok := v.(guildEmojisUpdatePayload)
	if !ok {
		return
	}

	if disgo.Caches().Config().CacheFlags.Missing(core.CacheFlagEmotes) {
		return
	}

	var (
		emojiCache    = disgo.Caches().EmojiCache().EmojiCache(payload.GuildID)
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
				updatedEmojis[current.ID] = disgo.EntityBuilder().CreateEmoji(payload.GuildID, current, core.CacheStrategyYes)
			}
		} else {
			newEmojis[current.ID] = disgo.EntityBuilder().CreateEmoji(payload.GuildID, current, core.CacheStrategyYes)
		}
	}

	for emojiID := range oldEmojis {
		disgo.Caches().EmojiCache().Uncache(payload.GuildID, emojiID)
	}

	genericGuildEvent := &events.GenericGuildEvent{
		GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
		GuildID:      payload.GuildID,
		Guild:        disgo.Caches().GuildCache().Get(payload.GuildID),
	}

	for _, emoji := range newEmojis {
		eventManager.Dispatch(&events.EmojiCreateEvent{
			GenericEmojiEvent: &events.GenericEmojiEvent{
				GenericGuildEvent: genericGuildEvent,
				Emoji:             emoji,
			},
		})
	}

	for _, emoji := range updatedEmojis {
		eventManager.Dispatch(&events.EmojiUpdateEvent{
			GenericEmojiEvent: &events.GenericEmojiEvent{
				GenericGuildEvent: genericGuildEvent,
				Emoji:             emoji,
			},
		})
	}

	for _, emoji := range oldEmojis {
		eventManager.Dispatch(&events.EmojiDeleteEvent{
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
