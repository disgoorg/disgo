package handlers

import (
	log "github.com/sirupsen/logrus"

	"github.com/DiscoOrg/disgo/api"
	"github.com/DiscoOrg/disgo/api/events"
)

type GuildCreateHandler struct{}

func (h GuildCreateHandler) New() interface{} {
	return &api.Guild{}
}

func (h GuildCreateHandler) Handle(disgo api.Disgo, eventManager api.EventManager, i interface{}) {
	guild, ok := i.(*api.Guild)
	if !ok {
		return
	}
	guild.Disgo = disgo
	wasUnavailable := disgo.Cache().UnavailableGuild(guild.ID) != nil
	disgo.Cache().CacheGuild(guild)
	for i := range guild.Channels {
		channel := guild.Channels[i]
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
			disgo.Cache().CacheCategory(&api.CategoryChannel{
				GuildChannel: *channel,
			})
		case api.ChannelTypeStore:
			disgo.Cache().CacheStoreChannel(&api.StoreChannel{
				GuildChannel: *channel,
			})
		}
	}

	if wasUnavailable {
		disgo.Cache().UncacheUnavailableGuild(guild.ID)
		eventManager.Dispatch(events.GuildAvailableEvent{
			GenericGuildEvent: events.GenericGuildEvent{
				Event: api.Event{
					Disgo: disgo,
				},
				GuildID: guild.ID,
			},
		})
	} else {
		// guild join
		eventManager.Dispatch(events.GuildJoinEvent{
			GenericGuildEvent: events.GenericGuildEvent{
				Event: api.Event{
					Disgo: disgo,
				},
				GuildID: guild.ID,
			},
		})
	}
}

type GuildDeleteHandler struct{}

func (h GuildDeleteHandler) New() interface{} {
	return &api.UnavailableGuild{}
}

func (h GuildDeleteHandler) Handle(disgo api.Disgo, eventManager api.EventManager, i interface{}) {
	unavailableGuild, ok := i.(*api.UnavailableGuild)
	if !ok {
		return
	}
	log.Infof("GuildDeleteEvent: %v", unavailableGuild)
	disgo.Cache().UncacheGuild(unavailableGuild.ID)
}

// GuildUpdateEvent payload from GUILD_DELETE gateways event sent by discord
type GuildUpdateEvent struct {
	Guild api.Guild
}

type GuildUpdateHandler struct{}

func (h GuildUpdateHandler) New() interface{} {
	return &GuildUpdateEvent{}
}

func (h GuildUpdateHandler) Handle(disgo api.Disgo, eventManager api.EventManager, i interface{}) {
	guild, ok := i.(*GuildUpdateEvent)
	if !ok {
		return
	}
	log.Infof("GuildUpdateEvent: %v", guild)
}
