package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

var _ HTTPServerEventHandler = (*httpserverHandlerInteractionCreate)(nil)

// httpserverHandlerInteractionCreate handles core.InteractionCreateWebhookEvent
type httpserverHandlerInteractionCreate struct{}

// New constructs a new payload receiver for the raw gateway event
func (h *httpserverHandlerInteractionCreate) New() interface{} {
	return &discord.Interaction{}
}

// HandleHTTPEvent handles the specific raw gateway event
func (h *httpserverHandlerInteractionCreate) HandleHTTPEvent(bot *Bot, c chan<- discord.InteractionResponse, v interface{}) {
	unmarshalInteraction := *v.(*discord.Interaction)

	// we just want to pong all pings
	// no need for any event
	if unmarshalInteraction.Type == discord.InteractionTypePing {
		bot.Logger.Infof("received interaction ping. responding pong")
		c <- discord.InteractionResponse{
			Type: discord.InteractionResponseTypePong,
		}
		return
	}
	HandleInteraction(bot, -1, c, unmarshalInteraction)
}
