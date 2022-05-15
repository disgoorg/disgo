package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// gatewayHandlerChannelCreate handles discord.GatewayEventTypeApplicationCommandPermissionsUpdate
type gatewayHandlerApplicationCommandPermissionsUpdate struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *gatewayHandlerApplicationCommandPermissionsUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeApplicationCommandPermissionsUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerApplicationCommandPermissionsUpdate) New() any {
	return &discord.ApplicationCommandPermissions{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerApplicationCommandPermissionsUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber int, v any) {
	permissions := *v.(*discord.ApplicationCommandPermissions)

	client.EventManager().DispatchEvent(&events.GuildApplicationCommandPermissionsUpdate{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber),
		Permissions:  permissions,
	})
}
