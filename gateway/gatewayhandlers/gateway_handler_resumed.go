package gatewayhandlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
)

// gatewayHandlerReady handles discord.GatewayEventTypeResumed
type gatewayHandlerResumed struct{}

// EventType returns the gateway.EventType
func (h *gatewayHandlerResumed) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeResumed
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerResumed) New() interface{} {
	return nil
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerResumed) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, _ interface{}) {
	bot.EventManager.Dispatch(&events.ResumedEvent{
		GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
	})

}
