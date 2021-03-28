package handlers

import (
	"github.com/DiscoOrg/disgo/api"
)

// InteractionCreateWebhookHandler handles api.InteractionCreateWebhookEvent
type InteractionCreateWebhookHandler struct{}

// Name returns the raw gateway event name
func (h InteractionCreateWebhookHandler) Name() string {
	return api.InteractionCreateWebhookEvent
}

// New constructs a new payload receiver for the raw gateway event
func (h InteractionCreateWebhookHandler) New() interface{} {
	return &api.Interaction{}
}

// Handle handles the specific raw gateway event
func (h InteractionCreateWebhookHandler) Handle(disgo api.Disgo, eventManager api.EventManager, c chan interface{}, i interface{}) {
	interaction, ok := i.(*api.Interaction)
	if !ok {
		return
	}

	if interaction.Type == api.InteractionTypePing {
		c <- api.InteractionResponse{
			Type: api.InteractionResponseTypePong,
		}
		return
	}
	handleInteractions(disgo, eventManager, c, interaction)
}
