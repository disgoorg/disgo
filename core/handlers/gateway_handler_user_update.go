package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerVoiceServerUpdate handles discord.GatewayEventTypeUserUpdate
type gatewayHandlerUserUpdate struct{}

// EventType returns the discord.GatewayEventTypeUserUpdate
func (h *gatewayHandlerUserUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeUserUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerUserUpdate) New() any {
	return &discord.OAuth2User{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerUserUpdate) HandleGatewayEvent(bot core.Bot, sequenceNumber discord.GatewaySequence, v any) {
	user := *v.(*discord.OAuth2User)

	var oldUser discord.OAuth2User
	if bot.SelfUser() != nil {
		oldUser = *bot.SelfUser()
	}
	bot.SetSelfUser(user)

	bot.EventManager().Dispatch(&events.SelfUpdateEvent{
		GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
		SelfUser:     user,
		OldSelfUser:  oldUser,
	})

}
