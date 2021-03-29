package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

// GuildCreateHandler handles api.GuildCreateGatewayEvent
type GuildCreateHandler struct{}

// Name returns the raw gateway event name
func (h GuildCreateHandler) Name() string {
	return api.GuildCreateGatewayEvent
}

// New constructs a new payload receiver for the raw gateway event
func (h GuildCreateHandler) New() interface{} {
	return &api.FullGuild{}
}

// Handle handles the specific raw gateway event
func (h GuildCreateHandler) Handle(disgo api.Disgo, eventManager api.EventManager, i interface{}) {
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

	for _, channel := range fullGuild.Channels {
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

	for _, role := range fullGuild.Roles {
		role.Disgo = disgo
		role.GuildID = guild.ID
		disgo.Cache().CacheRole(role)
	}

	for _, voiceState := range fullGuild.VoiceStates {
		voiceState.Disgo = disgo
		voiceState.GuildID = guild.ID
		disgo.Cache().CacheVoiceState(voiceState)
	}

	genericGuildEvent := events.GenericGuildEvent{
		Event: api.Event{
			Disgo: disgo,
		},
		GuildID: guild.ID,
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
