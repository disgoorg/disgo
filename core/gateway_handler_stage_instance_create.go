package core

import "github.com/DisgoOrg/disgo/discord"

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
func (h *gatewayHandlerStageInstanceCreate) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	stageInstance := *v.(*discord.StageInstance)

	bot.EventManager.Dispatch(&StageInstanceCreateEvent{
		GenericStageInstanceEvent: &GenericStageInstanceEvent{
			GenericGuildChannelEvent: &GenericGuildChannelEvent{
				GenericChannelEvent: &GenericChannelEvent{
					GenericEvent: NewGenericEvent(bot, sequenceNumber),
					ChannelID:    stageInstance.ChannelID,
					Channel:      bot.Caches.ChannelCache().Get(stageInstance.ChannelID),
				},
				GuildID: stageInstance.GuildID,
			},
			StageInstance: bot.EntityBuilder.CreateStageInstance(stageInstance, CacheStrategyYes),
		},
	})
}
