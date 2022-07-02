package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

type gatewayHandlerStageInstanceUpdate struct{}

func (h *gatewayHandlerStageInstanceUpdate) EventType() gateway.EventType {
	return gateway.EventTypeStageInstanceUpdate
}

func (h *gatewayHandlerStageInstanceUpdate) New() any {
	return &discord.StageInstance{}
}

func (h *gatewayHandlerStageInstanceUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	stageInstance := *v.(*discord.StageInstance)

	oldStageInstance, _ := client.Caches().StageInstances().Get(stageInstance.GuildID, stageInstance.ID)
	client.Caches().StageInstances().Put(stageInstance.GuildID, stageInstance.ID, stageInstance)

	client.EventManager().DispatchEvent(&events.StageInstanceUpdate{
		GenericStageInstance: &events.GenericStageInstance{
			GenericEvent:    events.NewGenericEvent(client, sequenceNumber, shardID),
			StageInstanceID: stageInstance.ID,
			StageInstance:   stageInstance,
		},
		OldStageInstance: oldStageInstance,
	})
}
