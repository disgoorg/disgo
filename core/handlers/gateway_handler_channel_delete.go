package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	events2 "github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerChannelDelete handles core.GatewayEventChannelDelete
type gatewayHandlerChannelDelete struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *gatewayHandlerChannelDelete) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeChannelDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerChannelDelete) New() interface{} {
	return &discord.UnmarshalChannel{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerChannelDelete) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	payload := v.(*discord.UnmarshalChannel).Channel

	bot.Caches.Channels().Remove(payload.ID())
	channel := bot.EntityBuilder.CreateChannel(payload, core.CacheStrategyNo)

	if ch, ok := channel.(core.GuildChannel); ok {
		bot.EventManager.Dispatch(&events2.GuildChannelDeleteEvent{
			GenericGuildChannelEvent: &events2.GenericGuildChannelEvent{
				GenericEvent: events2.NewGenericEvent(bot, sequenceNumber),
				ChannelID:    channel.ID(),
				Channel:      ch,
				GuildID:      ch.GuildID(),
			},
		})
	} else if ch, ok := channel.(*core.DMChannel); ok {
		bot.EventManager.Dispatch(&events2.DMChannelDeleteEvent{
			GenericDMChannelEvent: &events2.GenericDMChannelEvent{
				GenericEvent: events2.NewGenericEvent(bot, sequenceNumber),
				ChannelID:    channel.ID(),
				Channel:      ch,
			},
		})
	}
}
