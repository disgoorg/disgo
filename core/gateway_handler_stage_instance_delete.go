package core

import "github.com/DisgoOrg/disgo/discord"

// gatewayHandlerStageInstanceDelete handles discord.GatewayEventTypeStageInstanceDelete
type gatewayHandlerStageInstanceDelete struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *gatewayHandlerStageInstanceDelete) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeStageInstanceDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerStageInstanceDelete) New() interface{} {
	return &discord.StageInstance{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerStageInstanceDelete) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	stageInstance := *v.(*discord.StageInstance)

	bot.Caches.StageInstanceCache().Remove(stageInstance.ID)

	if channel := bot.Caches.ChannelCache().Get(stageInstance.ChannelID); channel != nil {
		channel.StageInstanceID = nil
	}

	bot.EventManager.Dispatch(&StageInstanceDeleteEvent{
		GenericStageInstanceEvent: &GenericStageInstanceEvent{
			GenericGuildChannelEvent: &GenericGuildChannelEvent{
				GenericChannelEvent: &GenericChannelEvent{
					GenericEvent: NewGenericEvent(bot, sequenceNumber),
					ChannelID:    stageInstance.ChannelID,
					Channel:      bot.Caches.ChannelCache().Get(stageInstance.ChannelID),
				},
				GuildID: stageInstance.GuildID,
			},
			StageInstance: bot.EntityBuilder.CreateStageInstance(stageInstance, CacheStrategyNo),
		},
	})
}
