package handlers

import (
	log "github.com/sirupsen/logrus"

	"github.com/DiscoOrg/disgo/api"
)

// GuildCreatePayload payload from GUILD_CREATE gateways event sent by discord
type InteractionCreateEvent struct {
	api.Interaction
}

type InteractionCreateHandler struct{}

func (h InteractionCreateHandler) New() interface{} {
	return &InteractionCreateEvent{}
}

func (h InteractionCreateHandler) Handle(disgo api.Disgo, eventManager api.EventManager, i interface{}) {
	interaction, ok := i.(*InteractionCreateEvent)
	if !ok {
		return
	}
	log.Infof("interaction received: %v", interaction)
}
