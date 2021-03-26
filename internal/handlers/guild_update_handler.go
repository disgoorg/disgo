package handlers

import (
	"github.com/DiscoOrg/disgo/api"
	"github.com/DiscoOrg/disgo/api/events"
)

type GuildUpdateHandler struct{}

func (h GuildUpdateHandler) New() interface{} {
	return &api.Guild{}
}

func (h GuildUpdateHandler) Handle(disgo api.Disgo, eventManager api.EventManager, i interface{}) {
	guild, ok := i.(*api.Guild)
	if !ok {
		return
	}

	oldGuild := *disgo.Cache().Guild(guild.ID)
	disgo.Cache().CacheGuild(guild)

	genericGuildEvent := events.GenericGuildEvent{
		Event:   api.Event{
			Disgo: disgo,
		},
		GuildID: guild.ID,
	}

	eventManager.Dispatch(genericGuildEvent)
	eventManager.Dispatch(events.GuildUpdateEvent{
		GenericGuildEvent: genericGuildEvent,
		Guild:             guild,
		OldGuild:          &oldGuild,
	})

}
