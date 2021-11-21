package gatewayhandlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
)

// gatewayHandlerReady handles discord.GatewayEventTypeInvalidSession
type gatewayHandlerInvalidSession struct{}

// EventType returns the gateway.EventType
func (h *gatewayHandlerInvalidSession) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeInvalidSession
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerInvalidSession) New() interface{} {
	return new(bool)
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerInvalidSession) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	resumable := *v.(*bool)

	// TODO: this should be used by the gateway
	bot.EventManager.Dispatch(&events.InvalidSessionEvent{
		GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
		MayResume:    resumable,
	})

}
