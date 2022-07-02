package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

type gatewayHandlerVoiceServerUpdate struct{}

func (h *gatewayHandlerVoiceServerUpdate) EventType() gateway.EventType {
	return gateway.EventTypeVoiceServerUpdate
}

func (h *gatewayHandlerVoiceServerUpdate) New() any {
	return &discord.VoiceServerUpdate{}
}

func (h *gatewayHandlerVoiceServerUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*discord.VoiceServerUpdate)

	client.EventManager().DispatchEvent(&events.VoiceServerUpdate{
		GenericEvent:      events.NewGenericEvent(client, sequenceNumber, shardID),
		VoiceServerUpdate: payload,
	})
}
