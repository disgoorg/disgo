package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	events2 "github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
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
	bot.EventManager.Dispatch(&events2.InvalidSessionEvent{
		GenericEvent: events2.NewGenericEvent(bot, sequenceNumber),
		MayResume:    resumable,
	})

}
