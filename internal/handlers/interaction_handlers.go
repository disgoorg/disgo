package handlers

import (
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
	if interaction.Member != nil {
		disgo.Cache().CacheMember(interaction.Member)
	}
	if interaction.User != nil {
		disgo.Cache().CacheUser(interaction.User)
	}

	genericInteractionEvent := events.GenericInteractionEvent{
		Event: api.Event{
			Disgo: disgo,
		},
		Interaction: *interaction,
	}

	eventManager.Dispatch(genericInteractionEvent)

	if interaction.Data != nil {
		eventManager.Dispatch(events.SlashCommandEvent{
			GenericInteractionEvent: genericInteractionEvent,
			InteractionData: *interaction.Data,
		})
	}

}
