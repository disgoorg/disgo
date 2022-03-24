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
	return &discord.InviteDeleteGatewayEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerInviteDelete) HandleGatewayEvent(client bot.Client, sequenceNumber discord.GatewaySequence, v any) {
	payload := *v.(*discord.InviteDeleteGatewayEvent)

	client.EventManager().DispatchEvent(&events.InviteDeleteEvent{
		GenericInviteEvent: &events.GenericInviteEvent{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber),
			GuildID:      payload.GuildID,
			ChannelID:    payload.ChannelID,
			Code:         payload.Code,
		},
	})
}
