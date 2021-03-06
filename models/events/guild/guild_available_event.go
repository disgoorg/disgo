package guild

import (
	"github.com/DiscoOrg/disgo/models"
	"github.com/DiscoOrg/disgo/models/events"
)

type GuildAvailableEvent struct {
	events.Event
	guild models.Guild
}

func (g GuildAvailableEvent) Guild() models.Guild {
	return g.guild
}