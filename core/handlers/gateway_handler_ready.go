package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerReady handles discord.GatewayEventTypeReady
type gatewayHandlerReady struct{}

// EventType returns the gateway.EventType
func (h *gatewayHandlerReady) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeReady
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerReady) New() interface{} {
	return &discord.GatewayEventReady{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerReady) HandleGatewayEvent(bot core.Bot, sequenceNumber discord.GatewaySequence, v interface{}) {
	readyEvent := *v.(*discord.GatewayEventReady)

	var shardID int
	if readyEvent.Shard != nil {
		shardID = readyEvent.Shard[0]
	}

	bot.HandleReadyEvent(readyEvent)

	for _, guild := range readyEvent.Guilds {
		bot.Caches().Guilds().SetUnready(shardID, guild.ID)
	}

	bot.EventManager().Dispatch(&events.ReadyEvent{
		GenericEvent:      events.NewGenericEvent(bot, sequenceNumber),
		GatewayEventReady: readyEvent,
	})

}
