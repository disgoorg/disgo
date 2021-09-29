package gatewayhandlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
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
func (h *gatewayHandlerChannelDelete) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	channel := *v.(*discord.Channel)

	genericChannelEvent := &events.GenericChannelEvent{
		GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
		ChannelID:    channel.ID,
		Channel:      bot.EntityBuilder.CreateChannel(channel, core.CacheStrategyNo),
	}

	if channel.GuildID != nil {
		bot.EventManager.Dispatch(&events.GuildChannelDeleteEvent{
			GenericGuildChannelEvent: &events.GenericGuildChannelEvent{
				GenericChannelEvent: genericChannelEvent,
				GuildID:             *channel.GuildID,
			},
		})
	} else {
		bot.EventManager.Dispatch(&events.DMChannelDeleteEvent{
			GenericDMChannelEvent: &events.GenericDMChannelEvent{
				GenericChannelEvent: genericChannelEvent,
			},
		})
	}
}
