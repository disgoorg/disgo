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
		channel.Disgo = disgo
		channel.GuildID = guild.ID
		switch channel.Type {
		case api.ChannelTypeText, api.ChannelTypeNews:
			disgo.Cache().CacheTextChannel(&api.TextChannel{
				GuildChannel: *channel,
				MessageChannel: api.MessageChannel{
					Channel: channel.Channel,
				},
			})
		case api.ChannelTypeVoice:
			disgo.Cache().CacheVoiceChannel(&api.VoiceChannel{
				GuildChannel: *channel,
			})
		case api.ChannelTypeCategory:
			disgo.Cache().CacheCategory(&api.Category{
				GuildChannel: *channel,
			})
		case api.ChannelTypeStore:
			disgo.Cache().CacheStoreChannel(&api.StoreChannel{
				GuildChannel: *channel,
			})
		}
	}

	for i := range fullGuild.Roles {
		role := fullGuild.Roles[i]
		role.Disgo = disgo
		role.GuildID = guild.ID
		disgo.Cache().CacheRole(role)
	}

	for i := range fullGuild.Members {
		member := fullGuild.Members[i]
		member.Disgo = disgo
		member.GuildID = guild.ID
		disgo.Cache().CacheMember(member)
	}

	for i := range fullGuild.VoiceStates {
		voiceState := fullGuild.VoiceStates[i]
		voiceState.Disgo = disgo
		disgo.Cache().CacheVoiceState(voiceState)
	}

	/*for i := range fullGuild.Emotes {
		emote := fullGuild.Emotes[i]
		emote.Disgo = disgo
		emote.GuildID = guild.ID
		disgo.Cache().CacheEmote(emote)
	}*/

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
		// guild join
		eventManager.Dispatch(events.GuildJoinEvent{
			GenericGuildEvent: genericGuildEvent,
			Guild:             guild,
		})
	}
}
