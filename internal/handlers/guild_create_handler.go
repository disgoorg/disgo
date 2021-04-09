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

	guild := fullGuild.Guild

	guild.Disgo = disgo
	oldGuild := disgo.Cache().Guild(guild.ID)
	var wasUnavailable bool
	if oldGuild == nil {
		wasUnavailable = true
	} else {
		wasUnavailable = oldGuild.Unavailable
	}

	disgo.Cache().CacheGuild(guild)
	for i := range fullGuild.Channels {
		channel := fullGuild.Channels[i]
		channel.GuildID = &guild.ID
		switch channel.Type {
		case api.ChannelTypeText, api.ChannelTypeNews:
			disgo.EntityBuilder().CreateTextChannel(channel, true)
		case api.ChannelTypeVoice:
			disgo.EntityBuilder().CreateVoiceChannel(channel, true)
		case api.ChannelTypeCategory:
			disgo.EntityBuilder().CreateCategory(channel, true)
		case api.ChannelTypeStore:
			disgo.EntityBuilder().CreateStoreChannel(channel, true)
		}
	}

	for i := range fullGuild.Roles {
		disgo.EntityBuilder().CreateRole(guild.ID, fullGuild.Roles[i], true)
	}

	for i := range fullGuild.Members {
		disgo.EntityBuilder().CreateMember(guild.ID, fullGuild.Members[i], true)
	}

	for i := range fullGuild.VoiceStates {
		disgo.EntityBuilder().CreateVoiceState(fullGuild.VoiceStates[i], true)
	}

	for i := range fullGuild.Emotes {
		disgo.EntityBuilder().CreateEmote(guild.ID, fullGuild.Emotes[i], true)
	}

	// TODO: presence
	/*for i := range fullGuild.Presences {
		presence := fullGuild.Presences[i]
		presence.Disgo = disgo
		disgo.Cache().CachePresence(presence)
	}*/

	genericGuildEvent := events.GenericGuildEvent{
		GenericEvent: events.NewEvent(disgo, sequenceNumber),
		GuildID:      guild.ID,
	}

	eventManager.Dispatch(genericGuildEvent)

	if wasUnavailable {
		eventManager.Dispatch(events.GuildAvailableEvent{
			GenericGuildEvent: genericGuildEvent,
			Guild:             guild,
		})
	} else {
		eventManager.Dispatch(events.GuildJoinEvent{
			GenericGuildEvent: genericGuildEvent,
			Guild:             guild,
		})
	}
}
