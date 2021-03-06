package guild

import (
	"github.com/DiscoOrg/disgo/models"
	"github.com/DiscoOrg/disgo/models/events"
)

type GuildLeaveEvent struct {
	events.Event
	guild models.Guild
}

func (g GuildLeaveEvent) Guild() models.Guild {
	return g.guild
}