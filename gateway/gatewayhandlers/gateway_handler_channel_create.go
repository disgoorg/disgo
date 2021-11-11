package gatewayhandlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
)

// gatewayHandlerChannelCreate handles core.GatewayEventChannelCreate
type gatewayHandlerChannelCreate struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *gatewayHandlerChannelCreate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeChannelCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerChannelCreate) New() interface{} {
	return &discord.UnmarshalChannel{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerChannelCreate) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	channel := v.(*discord.UnmarshalChannel).Channel

	genericChannelEvent := &events.GenericChannelEvent{
		GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
		ChannelID:    channel.ID(),
		Channel:      bot.EntityBuilder.CreateChannel(channel, core.CacheStrategyYes),
	}

	if ch, ok := channel.(discord.GuildChannel); ok {
		bot.EventManager.Dispatch(&events.GuildChannelCreateEvent{
			GenericGuildChannelEvent: &events.GenericGuildChannelEvent{
				GuildID:             ch.GuildID(),
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
