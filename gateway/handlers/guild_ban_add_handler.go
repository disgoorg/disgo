package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// gatewayHandlerGuildBanAdd handles discord.GatewayEventTypeGuildBanAdd
type gatewayHandlerGuildBanAdd struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerGuildBanAdd) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildBanAdd
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildBanAdd) New() any {
	return &discord.GatewayEventGuildBanAdd{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildBanAdd) HandleGatewayEvent(client bot.Client, sequenceNumber int, v any) {
	payload := *v.(*discord.GatewayEventGuildBanAdd)

	client.EventManager().DispatchEvent(&events.GuildBanEvent{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber),
		GuildID:      payload.GuildID,
		User:         payload.User,
	})
}
