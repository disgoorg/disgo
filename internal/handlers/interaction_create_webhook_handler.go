package handlers

import (
	"github.com/DisgoOrg/disgo/api"
)

// InteractionCreateWebhookHandler handles api.InteractionCreateWebhookEvent
type InteractionCreateWebhookHandler struct{}

// Event returns the raw gateway event Event
func (h InteractionCreateWebhookHandler) Event() api.GatewayEventType {
	return api.WebhookEventInteractionCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h InteractionCreateWebhookHandler) New() interface{} {
	return &api.Interaction{}
}

// Handle handles the specific raw gateway event
func (h InteractionCreateWebhookHandler) HandleWebhookEvent(disgo api.Disgo, eventManager api.EventManager, c chan interface{}, i interface{}) {
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
	handleInteraction(disgo, eventManager, -1, interaction, c)
}
