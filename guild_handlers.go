package disgo

import (
	log "github.com/sirupsen/logrus"

	"github.com/DiscoOrg/disgo/models"
)

// GuildCreateEvent payload from GUILD_CREATE gateways event sent by discord
type GuildCreateEvent struct {
	Guild models.Guild
}

type GuildCreateHandler struct {}

func (h GuildCreateHandler) New() interface{} {
	return &GuildCreateEvent{}
}

func (h GuildCreateHandler) Handle(eventManager EventManager, i interface{}) {
	guild, ok := i.(*GuildCreateEvent)
	if !ok {
		return
	}
	log.Infof("GuildCreateEvent: %v", guild)
}

// GuildDeleteEvent payload from GUILD_DELETE gateways event sent by discord
type GuildDeleteEvent struct {
	Guild models.UnavailableGuild
}

type GuildDeleteHandler struct {}

func (h GuildDeleteHandler) New() interface{} {
	return &GuildDeleteEvent{}
}

func (h GuildDeleteHandler) Handle(eventManager EventManager, i interface{}) {
	guild, ok := i.(*GuildDeleteEvent)
	if !ok {
		return
	}
	log.Infof("GuildDeleteEvent: %v", guild)
}

// GuildUpdateEvent payload from GUILD_DELETE gateways event sent by discord
type GuildUpdateEvent struct {
	Guild models.Guild
}

type GuildUpdateHandler struct {}

func (h GuildUpdateHandler) New() interface{} {
	return &GuildUpdateEvent{}
}

func (h GuildUpdateHandler) Handle(eventManager EventManager, i interface{}) {
	guild, ok := i.(*GuildUpdateEvent)
	if !ok {
		return
	}
	log.Infof("GuildUpdateEvent: %v", guild)
}