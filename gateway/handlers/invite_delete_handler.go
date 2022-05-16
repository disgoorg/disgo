package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// gatewayHandlerInviteDelete handles discord.GatewayEventTypeInviteDelete
type gatewayHandlerInviteDelete struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerInviteDelete) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeInviteDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerInviteDelete) New() any {
	return &discord.GatewayEventInviteDelete{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerInviteDelete) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*discord.GatewayEventInviteDelete)

	client.EventManager().DispatchEvent(&events.InviteDeleteEvent{
		GenericInviteEvent: &events.GenericInviteEvent{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildID:      payload.GuildID,
			ChannelID:    payload.ChannelID,
			Code:         payload.Code,
		},
	})
}
