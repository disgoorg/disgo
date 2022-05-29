package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

type gatewayHandlerStageInstanceDelete struct{}

func (h *gatewayHandlerStageInstanceDelete) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeStageInstanceDelete
}

func (h *gatewayHandlerStageInstanceDelete) New() any {
	return &discord.StageInstance{}
}

func (h *gatewayHandlerStageInstanceDelete) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	stageInstance := *v.(*discord.StageInstance)

	client.Caches().StageInstances().Remove(stageInstance.GuildID, stageInstance.ID)

	client.EventManager().DispatchEvent(&events.StageInstanceDelete{
		GenericStageInstance: &events.GenericStageInstance{
			GenericEvent:    events.NewGenericEvent(client, sequenceNumber, shardID),
			StageInstanceID: stageInstance.ID,
			StageInstance:   stageInstance,
		},
	})
}
