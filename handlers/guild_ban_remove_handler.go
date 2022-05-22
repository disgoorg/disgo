package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// gatewayHandlerGuildBanRemove handles discord.GatewayEventTypeGuildBanRemove
type gatewayHandlerGuildBanRemove struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerGuildBanRemove) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildBanRemove
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildBanRemove) New() any {
	return &discord.GatewayEventGuildBanRemove{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildBanRemove) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*discord.GatewayEventGuildBanRemove)

	client.EventManager().DispatchEvent(&events.GuildUnban{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		GuildID:      payload.GuildID,
		User:         payload.User,
	})
}
