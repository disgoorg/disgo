package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// ChannelCreateHandler handles core.GatewayEventChannelCreate
type ChannelCreateHandler struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *ChannelCreateHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeChannelCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *ChannelCreateHandler) New() interface{} {
	return &discord.Channel{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *ChannelCreateHandler) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	channel := *v.(*discord.Channel)

	genericChannelEvent := &GenericChannelEvent{
		GenericEvent: NewGenericEvent(bot, sequenceNumber),
		ChannelID:    channel.ID,
		Channel:      bot.EntityBuilder.CreateChannel(channel, CacheStrategyYes),
	}

	if channel.GuildID != nil {
		bot.EventManager.Dispatch(&GuildChannelCreateEvent{
			GenericGuildChannelEvent: &GenericGuildChannelEvent{
				GuildID:             *channel.GuildID,
				GenericChannelEvent: genericChannelEvent,
			},
		})
	} else {
		bot.EventManager.Dispatch(&DMChannelCreateEvent{
			GenericDMChannelEvent: &GenericDMChannelEvent{
				GenericChannelEvent: genericChannelEvent,
			},
		})
	}
}
