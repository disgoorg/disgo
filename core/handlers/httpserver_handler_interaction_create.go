package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

var _ core.HTTPServerEventHandler = (*httpserverHandlerInteractionCreate)(nil)

// httpserverHandlerInteractionCreate handles core.InteractionCreateWebhookEvent
type httpserverHandlerInteractionCreate struct{}

// New constructs a new payload receiver for the raw gateway event
func (h *httpserverHandlerInteractionCreate) New() interface{} {
	return &discord.UnmarshalInteraction{}
}

// HandleHTTPEvent handles the specific raw gateway event
func (h *httpserverHandlerInteractionCreate) HandleHTTPEvent(bot *core.Bot, c chan<- discord.InteractionResponse, v interface{}) {
	interaction := (*v.(*discord.UnmarshalInteraction)).Interaction

	// we just want to pong all pings
	// no need for any event
	if interaction.InteractionType() == discord.InteractionTypePing {
		bot.Logger.Info("received http interaction ping. responding with pong")
		c <- discord.InteractionResponse{
			Type: discord.InteractionCallbackTypePong,
		}
		return
	}
	HandleInteraction(bot, -1, c, interaction)
}
