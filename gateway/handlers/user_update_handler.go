package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
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
func (h *gatewayHandlerUserUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber int, v any) {
	user := *v.(*discord.OAuth2User)

	var oldUser discord.OAuth2User
	if client.SelfUser() != nil {
		oldUser = *client.SelfUser()
	}
	client.SetSelfUser(user)

	client.EventManager().DispatchEvent(&events.SelfUpdateEvent{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber),
		SelfUser:     user,
		OldSelfUser:  oldUser,
	})

}
