package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// ChannelUpdateHandler handles api.GatewayEventChannelUpdate
type ChannelUpdateHandler struct{}

// EventType returns the api.GatewayGatewayEventType
func (h *ChannelUpdateHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeChannelUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *ChannelUpdateHandler) New() interface{} {
	return discord.Channel{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *ChannelUpdateHandler) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	channel, ok := v.(discord.Channel)
	if !ok {
		return
	}

	oldCoreChannel := bot.Caches.ChannelCache().GetCopy(channel.ID)

	genericChannelEvent := &events.GenericChannelEvent{
		GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
		ChannelID:    channel.ID,
		Channel:      bot.EntityBuilder.CreateChannel(channel, core.CacheStrategyNo),
	}

	if channel.GuildID != nil {
		bot.EventManager.Dispatch(&events.GuildChannelUpdateEvent{
			GenericGuildChannelEvent: &events.GenericGuildChannelEvent{
				GenericChannelEvent: genericChannelEvent,
				GuildID:             *channel.GuildID,
			},
			OldChannel: oldCoreChannel,
		})
	} else {
		bot.EventManager.Dispatch(&events.DMChannelUpdateEvent{
			GenericDMChannelEvent: &events.GenericDMChannelEvent{
				GenericChannelEvent: genericChannelEvent,
			},
			OldChannel: oldCoreChannel,
		})
	}
}
