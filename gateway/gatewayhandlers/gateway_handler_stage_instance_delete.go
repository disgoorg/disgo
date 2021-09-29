package gatewayhandlers

import (
	"github.com/DisgoOrg/disgo/core"
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
func (h *gatewayHandlerStageInstanceDelete) New() interface{} {
	return &discord.StageInstance{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerStageInstanceDelete) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	stageInstance := *v.(*discord.StageInstance)

	bot.Caches.StageInstanceCache().Remove(stageInstance.ID)

	if channel := bot.Caches.ChannelCache().Get(stageInstance.ChannelID); channel != nil {
		channel.StageInstanceID = nil
	}

	bot.EventManager.Dispatch(&events.StageInstanceDeleteEvent{
		GenericStageInstanceEvent: &events.GenericStageInstanceEvent{
			GenericGuildChannelEvent: &events.GenericGuildChannelEvent{
				GenericChannelEvent: &events.GenericChannelEvent{
					GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
					ChannelID:    stageInstance.ChannelID,
					Channel:      bot.Caches.ChannelCache().Get(stageInstance.ChannelID),
				},
				GuildID: stageInstance.GuildID,
			},
			StageInstance: bot.EntityBuilder.CreateStageInstance(stageInstance, core.CacheStrategyNo),
		},
	})
}
