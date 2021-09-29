package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerChannelDelete handles discord.GatewayEventTypeChannelDelete
type gatewayHandlerChannelDelete struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerChannelDelete) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeChannelDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerChannelDelete) New() interface{} {
	return &discord.Channel{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerChannelDelete) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	channel := *v.(*discord.Channel)

	genericChannelEvent := &GenericChannelEvent{
		GenericEvent: NewGenericEvent(bot, sequenceNumber),
		ChannelID:    channel.ID,
		Channel:      bot.EntityBuilder.CreateChannel(channel, CacheStrategyNo),
	}

	if channel.GuildID != nil {
		bot.EventManager.Dispatch(&GuildChannelDeleteEvent{
			GenericGuildChannelEvent: &GenericGuildChannelEvent{
				GenericChannelEvent: genericChannelEvent,
				GuildID:             *channel.GuildID,
			},
		})
	} else {
		bot.EventManager.Dispatch(&DMChannelDeleteEvent{
			GenericDMChannelEvent: &GenericDMChannelEvent{
				GenericChannelEvent: genericChannelEvent,
			},
		})
	}
}
