package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

type gatewayHandlerAutoModerationActionExecution struct{}

func (h *gatewayHandlerAutoModerationActionExecution) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeAutoModerationActionExecution
}

func (h *gatewayHandlerAutoModerationActionExecution) New() any {
	return &discord.GatewayEventAutoModerationActionExecution{}
}

func (h *gatewayHandlerAutoModerationActionExecution) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	actionExecution := *v.(*discord.GatewayEventAutoModerationActionExecution)

	client.EventManager().DispatchEvent(&events.AutoModerationActionExecution{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		GatewayEventAutoModerationActionExecution: actionExecution,
	})
}
