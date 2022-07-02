package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

type gatewayHandlerAutoModerationRuleCreate struct{}

func (h *gatewayHandlerAutoModerationRuleCreate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeAutoModerationRuleCreate
}

func (h *gatewayHandlerAutoModerationRuleCreate) New() any {
	return &discord.GatewayEventAutoModerationRuleCreate{}
}

func (h *gatewayHandlerAutoModerationRuleCreate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	rule := *v.(*discord.GatewayEventAutoModerationRuleCreate)

	client.EventManager().DispatchEvent(&events.AutoModerationRuleCreate{
		GenericAutoModerationRule: &events.GenericAutoModerationRule{
			GenericEvent:       events.NewGenericEvent(client, sequenceNumber, shardID),
			AutoModerationRule: rule.AutoModerationRule,
		},
	})
}
