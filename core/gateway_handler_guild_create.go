package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// GuildCreateHandler handles core.GuildCreateGatewayEvent
type GuildCreateHandler struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *GuildCreateHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *GuildCreateHandler) New() interface{} {
	return &discord.GatewayGuild{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *GuildCreateHandler) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	guild := *v.(*discord.GatewayGuild)

	oldCoreGuild := bot.Caches.GuildCache().GetCopy(guild.ID)
	wasUnavailable := true
	if oldCoreGuild != nil {
		wasUnavailable = oldCoreGuild.Unavailable
	}

	genericGuildEvent := &GenericGuildEvent{
		GenericEvent: NewGenericEvent(bot, sequenceNumber),
		GuildID:      guild.ID,
		Guild:        bot.EntityBuilder.CreateGuild(guild.Guild, CacheStrategyYes),
	}

	for _, channel := range guild.Channels {
		channel.GuildID = &guild.ID
		bot.EntityBuilder.CreateChannel(channel, CacheStrategyYes)
	}

	for _, role := range guild.Roles {
		role.GuildID = guild.ID
		bot.EntityBuilder.CreateRole(guild.ID, role, CacheStrategyYes)
	}

	for _, member := range guild.Members {
		bot.EntityBuilder.CreateMember(guild.ID, member, CacheStrategyYes)
	}

	for _, voiceState := range guild.VoiceStates {
		bot.EntityBuilder.CreateVoiceState(guild.ID, voiceState, CacheStrategyYes)
	}

	for _, emote := range guild.Emojis {
		bot.EntityBuilder.CreateEmoji(guild.ID, emote, CacheStrategyYes)
	}

	for _, stageInstance := range guild.StageInstances {
		bot.EntityBuilder.CreateStageInstance(stageInstance, CacheStrategyYes)
	}

	// TODO: presence

	if wasUnavailable {
		bot.EventManager.Dispatch(&GuildAvailableEvent{
			GenericGuildEvent: genericGuildEvent,
		})
	} else {
		bot.EventManager.Dispatch(&GuildJoinEvent{
			GenericGuildEvent: genericGuildEvent,
		})
	}
}
