package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

type gatewayHandlerAutoModerationRuleUpdate struct{}

func (h *gatewayHandlerAutoModerationRuleUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeAutoModerationRuleUpdate
}

func (h *gatewayHandlerAutoModerationRuleUpdate) New() any {
	return &discord.GatewayEventAutoModerationRuleUpdate{}
}

func (h *gatewayHandlerAutoModerationRuleUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	rule := *v.(*discord.GatewayEventAutoModerationRuleUpdate)

	client.EventManager().DispatchEvent(&events.AutoModerationRuleUpdate{
		GenericAutoModerationRule: &events.GenericAutoModerationRule{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			AutoModerationRule:         rule.AutoModerationRule,
		},
	})
}
