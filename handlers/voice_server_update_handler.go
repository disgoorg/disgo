package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

type gatewayHandlerVoiceServerUpdate struct{}

func (h *gatewayHandlerVoiceServerUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeVoiceServerUpdate
}

func (h *gatewayHandlerVoiceServerUpdate) New() any {
	return &discord.VoiceServerUpdate{}
}

func (h *gatewayHandlerVoiceServerUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*discord.VoiceServerUpdate)

	if client.VoiceManager() != nil {
		client.VoiceManager().HandleVoiceServerUpdate(payload)
	}

	client.EventManager().DispatchEvent(&events.VoiceServerUpdate{
		GenericEvent:      events.NewGenericEvent(client, sequenceNumber, shardID),
		VoiceServerUpdate: payload,
	})
}
