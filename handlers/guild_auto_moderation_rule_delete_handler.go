package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

type gatewayHandlerAutoModerationRuleDelete struct{}

func (h *gatewayHandlerAutoModerationRuleDelete) EventType() gateway.EventType {
	return gateway.EventTypeAutoModerationRuleDelete
}

func (h *gatewayHandlerAutoModerationRuleDelete) New() any {
	return &gateway.EventAutoModerationRuleDelete{}
}

func (h *gatewayHandlerAutoModerationRuleDelete) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	rule := *v.(*gateway.EventAutoModerationRuleDelete)

	client.EventManager().DispatchEvent(&events.AutoModerationRuleDelete{
		GenericAutoModerationRule: &events.GenericAutoModerationRule{
			GenericEvent:       events.NewGenericEvent(client, sequenceNumber, shardID),
			AutoModerationRule: rule.AutoModerationRule,
		},
	})
}
