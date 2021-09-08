package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway/handlers"
)

func init() {
	core.HTTPServerEventHandler = &InteractionCreateWebhookHandler{}
}

var _ core.HTTPEventHandler = (*InteractionCreateWebhookHandler)(nil)

// InteractionCreateWebhookHandler handles api.InteractionCreateWebhookEvent
type InteractionCreateWebhookHandler struct{}

// New constructs a new payload receiver for the raw gateway event
func (h *InteractionCreateWebhookHandler) New() interface{} {
	return discord.Interaction{}
}

// HandleHTTPEvent handles the specific raw gateway event
func (h *InteractionCreateWebhookHandler) HandleHTTPEvent(bot *core.Bot, c chan discord.InteractionResponse, v interface{}) {
	unmarshalInteraction, ok := v.(discord.Interaction)
	if !ok {
		return
	}

	// we just want to pong all pings
	// no need for any event
	if unmarshalInteraction.Type == discord.InteractionTypePing {
		bot.Logger.Debugf("received interaction ping")
		c <- discord.InteractionResponse{
			Type: discord.InteractionResponseTypePong,
		}
		return
	}
	handlers.HandleInteraction(bot, -1, c, unmarshalInteraction)
}
