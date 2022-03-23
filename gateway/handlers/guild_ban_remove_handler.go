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
	return &discord.GuildBanRemoveGatewayEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildBanRemove) HandleGatewayEvent(client bot.Client, sequenceNumber discord.GatewaySequence, v any) {
	payload := *v.(*discord.GuildBanRemoveGatewayEvent)

	client.EventManager().Dispatch(&events.GuildUnbanEvent{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber),
		GuildID:      payload.GuildID,
		User:         payload.User,
	})
}
