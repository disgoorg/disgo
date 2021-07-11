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
	return &api.FullInteraction{}
}

// HandleWebhookEvent handles the specific raw gateway event
func (h InteractionCreateWebhookHandler) HandleWebhookEvent(disgo api.Disgo, eventManager api.EventManager, c chan api.InteractionResponse, i interface{}) {
	fullInteraction, ok := i.(*api.FullInteraction)
	if !ok {
		return
	}

	// we just want to pong all pings
	// no need for any event
	if fullInteraction.Type == api.InteractionTypePing {
		disgo.Logger().Debugf("received interaction ping")
		c <- api.InteractionResponse{
			Type: api.InteractionResponseTypePong,
		}
		return
	}
	handleInteraction(disgo, eventManager, -1, fullInteraction, c)
}
