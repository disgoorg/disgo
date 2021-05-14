package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

// GuildCreateHandler handles api.GuildCreateGatewayEvent
type GuildCreateHandler struct{}

// Event returns the raw gateway event Event
func (h GuildCreateHandler) Event() api.GatewayEventType {
	return api.GatewayEventGuildCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h GuildCreateHandler) New() interface{} {
	return &api.FullGuild{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h GuildCreateHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	fullGuild, ok := i.(*api.FullGuild)
	if !ok {
		return
	}

	oldGuild := disgo.Cache().Guild(fullGuild.ID)
	wasUnavailable := true
	if oldGuild != nil {
		oldGuild = &*oldGuild
		wasUnavailable = oldGuild.Unavailable
	}
	guild := disgo.EntityBuilder().CreateGuild(fullGuild.Guild, api.CacheStrategyYes)

	for _, channel := range fullGuild.Channels {
		channel.GuildID_ = &guild.ID
		switch channel.Type() {
		case api.ChannelTypeText, api.ChannelTypeNews:
			disgo.EntityBuilder().CreateTextChannel(channel, api.CacheStrategyYes)
		case api.ChannelTypeVoice:
			disgo.EntityBuilder().CreateVoiceChannel(channel, api.CacheStrategyYes)
		case api.ChannelTypeCategory:
			disgo.EntityBuilder().CreateCategory(channel, api.CacheStrategyYes)
		case api.ChannelTypeStore:
			disgo.EntityBuilder().CreateStoreChannel(channel, api.CacheStrategyYes)
		}
	}

	for _, thread := range fullGuild.Threads {
		disgo.EntityBuilder().CreateThread(thread, api.CacheStrategyYes)
	}

	for _, role := range fullGuild.Roles {
		disgo.EntityBuilder().CreateRole(guild.ID, role, api.CacheStrategyYes)
	}

	for _, member := range fullGuild.Members {
		disgo.EntityBuilder().CreateMember(guild.ID, member, api.CacheStrategyYes)
	}

	for _, voiceState := range fullGuild.VoiceStates {
		disgo.EntityBuilder().CreateVoiceState(guild.ID, voiceState, api.CacheStrategyYes)
	}

	for _, emote := range fullGuild.Emotes {
		disgo.EntityBuilder().CreateEmote(guild.ID, emote, api.CacheStrategyYes)
	}

	// TODO: presence
	/*for i := range fullGuild.Presences {
		presence := fullGuild.Presences[i]
		presence.Disgo = disgo
		disgo.Cache().CachePresence(presence)
	}*/

	genericGuildEvent := events.GenericGuildEvent{
		GenericEvent: events.NewEvent(disgo, sequenceNumber),
		Guild:        guild,
	}
	eventManager.Dispatch(genericGuildEvent)

	if wasUnavailable {
		eventManager.Dispatch(events.GuildAvailableEvent{
			GenericGuildEvent: genericGuildEvent,
		})
	} else {
		eventManager.Dispatch(events.GuildJoinEvent{
			GenericGuildEvent: genericGuildEvent,
		})
	}
}
