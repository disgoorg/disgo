package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

type gatewayHandlerAutoModerationRuleDelete struct{}

func (h *gatewayHandlerAutoModerationRuleDelete) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeAutoModerationRuleDelete
}

func (h *gatewayHandlerAutoModerationRuleDelete) New() any {
	return &discord.GatewayEventAutoModerationRuleDelete{}
}

func (h *gatewayHandlerAutoModerationRuleDelete) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	rule := *v.(*discord.GatewayEventAutoModerationRuleDelete)

	client.EventManager().DispatchEvent(&events.AutoModerationRuleDelete{
		GenericAutoModerationRule: &events.GenericAutoModerationRule{
			GenericEvent:       events.NewGenericEvent(client, sequenceNumber, shardID),
			AutoModerationRule: rule.AutoModerationRule,
		},
	})
}
