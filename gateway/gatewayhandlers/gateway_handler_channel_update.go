package gatewayhandlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
)

// gatewayHandlerChannelUpdate handles core.GatewayEventChannelUpdate
type gatewayHandlerChannelUpdate struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *gatewayHandlerChannelUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeChannelUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerChannelUpdate) New() interface{} {
	return &discord.UnmarshalChannel{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerChannelUpdate) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	channel := v.(*discord.UnmarshalChannel).Channel

	oldCoreChannel := bot.Caches.ChannelCache().GetCopy(channel.ID())

	genericChannelEvent := &events.GenericChannelEvent{
		GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
		ChannelID:    channel.ID(),
		Channel:      bot.EntityBuilder.CreateChannel(channel, core.CacheStrategyNo),
	}

	if ch, ok := channel.(discord.GuildChannel); ok {
		bot.EventManager.Dispatch(&events.GuildChannelUpdateEvent{
			GenericGuildChannelEvent: &events.GenericGuildChannelEvent{
				GenericChannelEvent: genericChannelEvent,
				GuildID:             ch.GuildID(),
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
