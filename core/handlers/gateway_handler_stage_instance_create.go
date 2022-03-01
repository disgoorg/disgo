package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerStageInstanceCreate handles discord.GatewayEventTypeStageInstanceCreate
type gatewayHandlerStageInstanceCreate struct{}

// EventType returns the discord.GatewayEventTypeStageInstanceCreate
func (h *gatewayHandlerStageInstanceCreate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeStageInstanceCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerStageInstanceCreate) New() interface{} {
	return &discord.StageInstance{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerStageInstanceCreate) HandleGatewayEvent(bot *core.Bot, sequenceNumber discord.GatewaySequence, v interface{}) {
	payload := *v.(*discord.StageInstance)

	bot.EventManager.Dispatch(&events.StageInstanceCreateEvent{
		GenericStageInstanceEvent: &events.GenericStageInstanceEvent{
			GenericEvent:    events.NewGenericEvent(bot, sequenceNumber),
			StageInstanceID: payload.ID,
			StageInstance:   bot.EntityBuilder.CreateStageInstance(payload, core.CacheStrategyYes),
		},
	})
}
