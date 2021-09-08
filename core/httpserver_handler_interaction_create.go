package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

var _ HTTPServerEventHandler = (*InteractionCreateHTTPServerHandler)(nil)

// InteractionCreateHTTPServerHandler handles api.InteractionCreateWebhookEvent
type InteractionCreateHTTPServerHandler struct{}

// New constructs a new payload receiver for the raw gateway event
func (h *InteractionCreateHTTPServerHandler) New() interface{} {
	return discord.Interaction{}
}

// HandleHTTPEvent handles the specific raw gateway event
func (h *InteractionCreateHTTPServerHandler) HandleHTTPEvent(bot *Bot, c chan discord.InteractionResponse, v interface{}) {
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
	HandleInteraction(bot, -1, c, unmarshalInteraction)
}
