package core

import "github.com/DisgoOrg/disgo/discord"

// StageInstanceCreateHandler handles core.GatewayEventMessageCreate
type StageInstanceCreateHandler struct{}

// EventType returns the discord.GatewayEventTypeStageInstanceCreate
func (h *StageInstanceCreateHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeStageInstanceCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *StageInstanceCreateHandler) New() interface{} {
	return &discord.StageInstance{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *StageInstanceCreateHandler) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
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

// StageInstanceUpdateHandler handles core.GatewayEventMessageCreate
type StageInstanceUpdateHandler struct{}

// EventType returns the discord.GatewayEventTypeStageInstanceUpdate
func (h *StageInstanceUpdateHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeStageInstanceUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *StageInstanceUpdateHandler) New() interface{} {
	return &discord.StageInstance{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *StageInstanceUpdateHandler) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	stageInstance := *v.(*discord.StageInstance)

	oldStageInstance := bot.Caches.StageInstanceCache().GetCopy(stageInstance.ID)

	bot.EventManager.Dispatch(&StageInstanceUpdateEvent{
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
		OldStageInstance: oldStageInstance,
	})
}

// StageInstanceDeleteHandler handles discord.GatewayEventTypeStageInstanceDelete
type StageInstanceDeleteHandler struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *StageInstanceDeleteHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeStageInstanceDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h *StageInstanceDeleteHandler) New() interface{} {
	return &discord.StageInstance{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *StageInstanceDeleteHandler) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
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
