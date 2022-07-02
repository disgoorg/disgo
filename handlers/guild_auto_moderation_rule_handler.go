package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerAutoModerationRuleCreate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventAutoModerationRuleCreate) {
	client.EventManager().DispatchEvent(&events.AutoModerationRuleCreate{
		GenericAutoModerationRule: &events.GenericAutoModerationRule{
			GenericEvent:       events.NewGenericEvent(client, sequenceNumber, shardID),
			AutoModerationRule: event.AutoModerationRule,
		},
	})
}

func gatewayHandlerAutoModerationRuleUpdate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventAutoModerationRuleUpdate) {
	client.EventManager().DispatchEvent(&events.AutoModerationRuleUpdate{
		GenericAutoModerationRule: &events.GenericAutoModerationRule{
			GenericEvent:       events.NewGenericEvent(client, sequenceNumber, shardID),
			AutoModerationRule: event.AutoModerationRule,
		},
	})
}

func gatewayHandlerAutoModerationRuleDelete(client bot.Client, sequenceNumber int, shardID int, event gateway.EventAutoModerationRuleDelete) {
	client.EventManager().DispatchEvent(&events.AutoModerationRuleDelete{
		GenericAutoModerationRule: &events.GenericAutoModerationRule{
			GenericEvent:       events.NewGenericEvent(client, sequenceNumber, shardID),
			AutoModerationRule: event.AutoModerationRule,
		},
	})
}

func gatewayHandlerAutoModerationActionExecution(client bot.Client, sequenceNumber int, shardID int, event gateway.EventAutoModerationActionExecution) {
	client.EventManager().DispatchEvent(&events.AutoModerationActionExecution{
		GenericEvent:                       events.NewGenericEvent(client, sequenceNumber, shardID),
		EventAutoModerationActionExecution: event,
	})
}
