package handlers

import (
	log "github.com/sirupsen/logrus"

	"github.com/DiscoOrg/disgo/api"
	"github.com/DiscoOrg/disgo/api/events"
)

type InteractionCreateHandler struct{}

func (h InteractionCreateHandler) New() interface{} {
	return &api.Interaction{}
}

func (h InteractionCreateHandler) Handle(disgo api.Disgo, eventManager api.EventManager, i interface{}) {
	interaction, ok := i.(*api.Interaction)
	if !ok {
		return
	}
	log.Infof("interaction received: %v", interaction)
	eventManager.Dispatch(events.SlashCommandEvent{
		GenericInteractionEvent: events.GenericInteractionEvent{
			Event:         api.Event{
				Disgo:          disgo,
			},
			Token:         interaction.Token,
			InteractionID: interaction.ID,
			Guild:         disgo.Cache().Guild(interaction.GuildID),
			Member:        nil,
			User:          api.User{},
			Channel:       nil,
		},
		Name:                    interaction.Data.Name,
		CommandID:               interaction.Data.ID,
		Options:                 interaction.Data.Options,
	})
}
