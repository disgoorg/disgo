package guild

import (
	"github.com/DiscoOrg/disgo/models"
	"github.com/DiscoOrg/disgo/models/events"
)

type GuildUnavailableEvent struct {
	events.Event
	guild models.Guild
}

func (g GuildUnavailableEvent) Guild() models.Guild {
	return g.guild
}