package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway/handlers"
	"github.com/DisgoOrg/disgo/httpserver"
)

// InteractionCreateWebhookHandler handles api.InteractionCreateWebhookEvent
type InteractionCreateWebhookHandler struct{}

// EventType returns the api.GatewayEventType
func (h *InteractionCreateWebhookHandler) EventType() httpserver.EventType {
	return httpserver.EventTypeInteractionCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *InteractionCreateWebhookHandler) New() interface{} {
	return discord.UnmarshalInteraction{}
}

// HandleHTTPEvent handles the specific raw gateway event
func (h *InteractionCreateWebhookHandler) HandleHTTPEvent(disgo core.Disgo, eventManager core.EventManager, c chan discord.InteractionResponse, i interface{}) {
	unmarshalInteraction, ok := i.(discord.UnmarshalInteraction)
	if !ok {
		return
	}

	// we just want to pong all pings
	// no need for any event
	if unmarshalInteraction.Type == discord.InteractionTypePing {
		disgo.Logger().Debugf("received interaction ping")
		c <- discord.InteractionResponse{
			Type: discord.InteractionResponseTypePong,
		}
		return
	}
	handlers.HandleInteraction(disgo, eventManager, -1, c, unmarshalInteraction)
}
