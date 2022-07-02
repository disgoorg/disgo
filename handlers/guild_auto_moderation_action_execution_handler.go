package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

type gatewayHandlerAutoModerationActionExecution struct{}

func (h *gatewayHandlerAutoModerationActionExecution) EventType() gateway.EventType {
	return gateway.EventTypeAutoModerationActionExecution
}

func (h *gatewayHandlerAutoModerationActionExecution) New() any {
	return &gateway.EventAutoModerationActionExecution{}
}

func (h *gatewayHandlerAutoModerationActionExecution) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	actionExecution := *v.(*gateway.EventAutoModerationActionExecution)

	client.EventManager().DispatchEvent(&events.AutoModerationActionExecution{
		GenericEvent:                       events.NewGenericEvent(client, sequenceNumber, shardID),
		EventAutoModerationActionExecution: actionExecution,
	})
}
