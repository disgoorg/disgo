package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// GuildCreateHandler handles api.GuildCreateGatewayEvent
type GuildCreateHandler struct{}

// EventType returns the api.GatewayGatewayEventType
func (h *GuildCreateHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *GuildCreateHandler) New() interface{} {
	return discord.GatewayGuild{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *GuildCreateHandler) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	guild, ok := v.(discord.GatewayGuild)
	if !ok {
		return
	}

	oldCoreGuild := bot.Caches.GuildCache().GetCopy(guild.ID)
	wasUnavailable := true
	if oldCoreGuild != nil {
		wasUnavailable = oldCoreGuild.Unavailable
	}

	genericGuildEvent := &events.GenericGuildEvent{
		GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
		GuildID:      guild.ID,
		Guild:        bot.EntityBuilder.CreateGuild(guild.Guild, core.CacheStrategyYes),
	}

	for _, channel := range guild.Channels {
		channel.GuildID = &guild.ID
		bot.EntityBuilder.CreateChannel(channel, core.CacheStrategyYes)
	}

	for _, role := range guild.Roles {
		role.GuildID = guild.ID
		bot.EntityBuilder.CreateRole(guild.ID, role, core.CacheStrategyYes)
	}

	for _, member := range guild.Members {
		bot.EntityBuilder.CreateMember(guild.ID, member, core.CacheStrategyYes)
	}

	for _, voiceState := range guild.VoiceStates {
		bot.EntityBuilder.CreateVoiceState(guild.ID, voiceState, core.CacheStrategyYes)
	}

	for _, emote := range guild.Emojis {
		bot.EntityBuilder.CreateEmoji(guild.ID, emote, core.CacheStrategyYes)
	}

	for _, stageInstance := range guild.StageInstances {
		bot.EntityBuilder.CreateStageInstance(stageInstance, core.CacheStrategyYes)
	}

	// TODO: presence

	if wasUnavailable {
		bot.EventManager.Dispatch(&events.GuildAvailableEvent{
			GenericGuildEvent: genericGuildEvent,
		})
	} else {
		bot.EventManager.Dispatch(&events.GuildJoinEvent{
			GenericGuildEvent: genericGuildEvent,
		})
	}
}
