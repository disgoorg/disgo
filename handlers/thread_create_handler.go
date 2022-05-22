package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// gatewayHandlerChannelCreate handles discord.GatewayEventTypeThreadCreate
type gatewayHandlerThreadCreate struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerThreadCreate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeThreadCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerThreadCreate) New() any {
	return &discord.GatewayEventThreadCreate{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerThreadCreate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*discord.GatewayEventThreadCreate)

	client.Caches().Channels().Put(payload.ID(), payload.GuildThread)
	client.Caches().ThreadMembers().Put(payload.ID(), payload.ThreadMember.UserID, payload.ThreadMember)

	client.EventManager().DispatchEvent(&events.ThreadCreate{
		GenericThread: &events.GenericThread{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			ThreadID:     payload.ID(),
			GuildID:      payload.GuildID(),
			Thread:       payload.GuildThread,
		},
		ThreadMember: payload.ThreadMember,
	})
}
