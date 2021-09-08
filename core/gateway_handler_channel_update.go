package core

import (
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
	return &discord.Channel{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *ChannelUpdateHandler) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	channel := *v.(*discord.Channel)

	oldCoreChannel := bot.Caches.ChannelCache().GetCopy(channel.ID)

	genericChannelEvent := &GenericChannelEvent{
		GenericEvent: NewGenericEvent(bot, sequenceNumber),
		ChannelID:    channel.ID,
		Channel:      bot.EntityBuilder.CreateChannel(channel, CacheStrategyNo),
	}

	if channel.GuildID != nil {
		bot.EventManager.Dispatch(&GuildChannelUpdateEvent{
			GenericGuildChannelEvent: &GenericGuildChannelEvent{
				GenericChannelEvent: genericChannelEvent,
				GuildID:             *channel.GuildID,
			},
			OldChannel: oldCoreChannel,
		})
	} else {
		bot.EventManager.Dispatch(&DMChannelUpdateEvent{
			GenericDMChannelEvent: &GenericDMChannelEvent{
				GenericChannelEvent: genericChannelEvent,
			},
			OldChannel: oldCoreChannel,
		})
	}
}
