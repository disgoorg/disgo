package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// gatewayHandlerStageInstanceUpdate handles discord.GatewayEventTypeStageInstanceUpdate
type gatewayHandlerStageInstanceUpdate struct{}

// EventType returns the discord.GatewayEventTypeStageInstanceUpdate
func (h *gatewayHandlerStageInstanceUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeStageInstanceUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerStageInstanceUpdate) New() any {
	return &discord.StageInstance{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerStageInstanceUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber int, v any) {
	stageInstance := *v.(*discord.StageInstance)

	oldStageInstance, _ := client.Caches().StageInstances().Get(stageInstance.GuildID, stageInstance.ID)
	client.Caches().StageInstances().Put(stageInstance.GuildID, stageInstance.ID, stageInstance)

	client.EventManager().DispatchEvent(&events.StageInstanceUpdateEvent{
		GenericStageInstanceEvent: &events.GenericStageInstanceEvent{
			GenericEvent:    events.NewGenericEvent(client, sequenceNumber),
			StageInstanceID: stageInstance.ID,
			StageInstance:   stageInstance,
		},
		OldStageInstance: oldStageInstance,
	})
}
