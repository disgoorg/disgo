package gatewayhandlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
)

// gatewayHandlerStageInstanceCreate handles core.GatewayEventMessageCreate
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
func (h *gatewayHandlerStageInstanceCreate) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	stageInstance := *v.(*discord.StageInstance)

	bot.EventManager.Dispatch(&events.StageInstanceCreateEvent{
		GenericStageInstanceEvent: &events.GenericStageInstanceEvent{
			GenericGuildChannelEvent: &events.GenericGuildChannelEvent{
				GenericChannelEvent: &events.GenericChannelEvent{
					GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
					ChannelID:    stageInstance.ChannelID,
					Channel:      bot.Caches.ChannelCache().Get(stageInstance.ChannelID),
				},
				GuildID: stageInstance.GuildID,
			},
			StageInstance: bot.EntityBuilder.CreateStageInstance(stageInstance, core.CacheStrategyYes),
		},
	})
}
