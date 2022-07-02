package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

type gatewayHandlerAutoModerationRuleCreate struct{}

func (h *gatewayHandlerAutoModerationRuleCreate) EventType() gateway.EventType {
	return gateway.EventTypeAutoModerationRuleCreate
}

func (h *gatewayHandlerAutoModerationRuleCreate) New() any {
	return &gateway.EventAutoModerationRuleCreate{}
}

func (h *gatewayHandlerAutoModerationRuleCreate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	rule := *v.(*gateway.EventAutoModerationRuleCreate)

	client.EventManager().DispatchEvent(&events.AutoModerationRuleCreate{
		GenericAutoModerationRule: &events.GenericAutoModerationRule{
			GenericEvent:       events.NewGenericEvent(client, sequenceNumber, shardID),
			AutoModerationRule: rule.AutoModerationRule,
		},
	})
}
