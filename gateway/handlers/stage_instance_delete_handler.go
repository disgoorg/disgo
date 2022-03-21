package handlers

import (
	"github.com/DisgoOrg/disgo/bot"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
)

// gatewayHandlerStageInstanceDelete handles discord.GatewayEventTypeStageInstanceDelete
type gatewayHandlerStageInstanceDelete struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerStageInstanceDelete) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeStageInstanceDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerStageInstanceDelete) New() any {
	return &discord.StageInstance{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerStageInstanceDelete) HandleGatewayEvent(client bot.Client, sequenceNumber discord.GatewaySequence, v any) {
	stageInstance := *v.(*discord.StageInstance)

	client.Caches().StageInstances().Remove(stageInstance.GuildID, stageInstance.ID)

	client.EventManager().Dispatch(&events.StageInstanceDeleteEvent{
		GenericStageInstanceEvent: &events.GenericStageInstanceEvent{
			GenericEvent:    events.NewGenericEvent(client, sequenceNumber),
			StageInstanceID: stageInstance.ID,
			StageInstance:   stageInstance,
		},
	})
}
