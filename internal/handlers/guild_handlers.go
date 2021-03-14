package handlers

import (
	log "github.com/sirupsen/logrus"

	"github.com/DiscoOrg/disgo/api"
)


type GuildCreateHandler struct {}

func (h GuildCreateHandler) New() interface{} {
	return &api.Guild{}
}

func (h GuildCreateHandler) Handle(disgo api.Disgo, eventManager api.EventManager, i interface{}) {
	guild, ok := i.(*api.Guild)
	if !ok {
		return
	}
	log.Infof("GuildCreateEvent: %v", guild)
	guild.Disgo = disgo
	disgo.Cache().CacheGuild(guild)
}

type GuildDeleteHandler struct {}

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

type GuildUpdateHandler struct {}

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