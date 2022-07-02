package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerVoiceServerUpdate struct {}

func (h *gatewayHandlerVoiceServerUpdate) EventType() gateway.EventType {
	return gateway.EventTypeVoiceServerUpdate
}

func (h *gatewayHandlerVoiceServerUpdate) New() any {
	return &gateway.EventVoiceServerUpdate{}
}

func (h *gatewayHandlerVoiceServerUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*gateway.EventVoiceServerUpdate)

	client.EventManager().DispatchEvent(&events.VoiceServerUpdate{
		GenericEvent:           events.NewGenericEvent(client, sequenceNumber, shardID),
		EventVoiceServerUpdate: payload,
	})
}
