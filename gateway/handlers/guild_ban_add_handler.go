package handlers

import (
	"github.com/DisgoOrg/disgo/bot"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
)

// gatewayHandlerGuildBanAdd handles discord.GatewayEventTypeGuildBanAdd
type gatewayHandlerGuildBanAdd struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerGuildBanAdd) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildBanAdd
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildBanAdd) New() any {
	return &discord.GuildBanAddGatewayEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildBanAdd) HandleGatewayEvent(client bot.Client, sequenceNumber discord.GatewaySequence, v any) {
	payload := *v.(*discord.GuildBanAddGatewayEvent)

	client.EventManager().Dispatch(&events.GuildBanEvent{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber),
		GuildID:      payload.GuildID,
		User:         payload.User,
	})
}
