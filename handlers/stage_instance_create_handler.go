package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// gatewayHandlerStageInstanceCreate handles discord.GatewayEventTypeStageInstanceCreate
type gatewayHandlerStageInstanceCreate struct{}

// EventType returns the discord.GatewayEventTypeStageInstanceCreate
func (h *gatewayHandlerStageInstanceCreate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeStageInstanceCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerStageInstanceCreate) New() any {
	return &discord.StageInstance{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerStageInstanceCreate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	stageInstance := *v.(*discord.StageInstance)

	client.Caches().StageInstances().Put(stageInstance.GuildID, stageInstance.ID, stageInstance)

	client.EventManager().DispatchEvent(&events.StageInstanceCreate{
		GenericStageInstance: &events.GenericStageInstance{
			GenericEvent:    events.NewGenericEvent(client, sequenceNumber, shardID),
			StageInstanceID: stageInstance.ID,
			StageInstance:   stageInstance,
		},
	})
}
