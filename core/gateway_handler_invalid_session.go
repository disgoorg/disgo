package core

import (
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
func (h *gatewayHandlerInvalidSession) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	resumable := *v.(*bool)

	bot.EventManager.Dispatch(&InvalidSessionEvent{
		GenericEvent: NewGenericEvent(bot, sequenceNumber),
		MayResume:    resumable,
	})

}
