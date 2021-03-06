package guild

import (
	"github.com/DiscoOrg/disgo/models"
	"github.com/DiscoOrg/disgo/models/events"
)

type GuildJoinEvent struct {
	events.Event
	guild models.Guild
}

func (g GuildJoinEvent) Guild() models.Guild {
	return g.guild
}