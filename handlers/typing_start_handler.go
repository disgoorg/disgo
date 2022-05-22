package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// gatewayHandlerTypingStart handles discord.GatewayEventTypeInviteDelete
type gatewayHandlerTypingStart struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *gatewayHandlerTypingStart) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeTypingStart
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerTypingStart) New() any {
	return &discord.GatewayEventTypingStart{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerTypingStart) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*discord.GatewayEventTypingStart)

	client.EventManager().DispatchEvent(&events.UserTypingStart{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		ChannelID:    payload.ChannelID,
		GuildID:      payload.GuildID,
		UserID:       payload.UserID,
		Timestamp:    payload.Timestamp,
	})

	if payload.GuildID == nil {
		client.EventManager().DispatchEvent(&events.DMUserTypingStart{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			ChannelID:    payload.ChannelID,
			UserID:       payload.UserID,
			Timestamp:    payload.Timestamp,
		})
	} else {
		client.Caches().Members().Put(*payload.GuildID, payload.UserID, *payload.Member)
		client.EventManager().DispatchEvent(&events.GuildMemberTypingStart{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			ChannelID:    payload.ChannelID,
			UserID:       payload.UserID,
			GuildID:      *payload.GuildID,
			Timestamp:    payload.Timestamp,
			Member:       *payload.Member,
		})
	}
}
