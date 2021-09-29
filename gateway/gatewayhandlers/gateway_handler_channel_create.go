package gatewayhandlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
)

type gatewayHandlerChannelCreate struct{}

func (h *gatewayHandlerChannelCreate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeChannelCreate
}

func (h *gatewayHandlerChannelCreate) New() interface{} {
	return &discord.Channel{}
}

func (h *gatewayHandlerChannelCreate) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	channel := *v.(*discord.Channel)

	genericChannelEvent := &events.GenericChannelEvent{
		GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
		ChannelID:    channel.ID,
		Channel:      bot.EntityBuilder.CreateChannel(channel, core.CacheStrategyYes),
	}

	if channel.GuildID != nil {
		bot.EventManager.Dispatch(&events.GuildChannelCreateEvent{
			GenericGuildChannelEvent: &events.GenericGuildChannelEvent{
				GuildID:             *channel.GuildID,
				GenericChannelEvent: genericChannelEvent,
			},
		})
	} else {
		bot.EventManager.Dispatch(&events.DMChannelCreateEvent{
			GenericDMChannelEvent: &events.GenericDMChannelEvent{
				GenericChannelEvent: genericChannelEvent,
			},
		})
	}
}
