package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

type gatewayHandlerAutoModerationRuleUpdate struct{}

func (h *gatewayHandlerAutoModerationRuleUpdate) EventType() gateway.EventType {
	return gateway.EventTypeAutoModerationRuleUpdate
}

func (h *gatewayHandlerAutoModerationRuleUpdate) New() any {
	return &gateway.EventAutoModerationRuleUpdate{}
}

func (h *gatewayHandlerAutoModerationRuleUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	rule := *v.(*gateway.EventAutoModerationRuleUpdate)

	client.EventManager().DispatchEvent(&events.AutoModerationRuleUpdate{
		GenericAutoModerationRule: &events.GenericAutoModerationRule{
			GenericEvent:       events.NewGenericEvent(client, sequenceNumber, shardID),
			AutoModerationRule: rule.AutoModerationRule,
		},
	})
}
