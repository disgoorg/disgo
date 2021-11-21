package gatewayhandlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
)

// gatewayHandlerStageInstanceUpdate handles core.GatewayEventMessageCreate
type gatewayHandlerStageInstanceUpdate struct{}

// EventType returns the discord.GatewayEventTypeStageInstanceUpdate
func (h *gatewayHandlerStageInstanceUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeStageInstanceUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerStageInstanceUpdate) New() interface{} {
	return &discord.StageInstance{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerStageInstanceUpdate) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	stageInstance := *v.(*discord.StageInstance)

	oldStageInstance := bot.Caches.StageInstances().GetCopy(stageInstance.ID)

	bot.EventManager.Dispatch(&events.StageInstanceUpdateEvent{
		GenericStageInstanceEvent: &events.GenericStageInstanceEvent{
			GenericEvent:    events.NewGenericEvent(bot, sequenceNumber),
			StageInstanceID: stageInstance.ID,
			StageInstance:   bot.EntityBuilder.CreateStageInstance(stageInstance, core.CacheStrategyYes),
		},
		OldStageInstance: oldStageInstance,
	})
}
