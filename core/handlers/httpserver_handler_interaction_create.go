package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

var _ core.HTTPServerEventHandler = (*httpserverHandlerInteractionCreate)(nil)

type httpserverHandlerInteractionCreate struct{}

func (h *httpserverHandlerInteractionCreate) New() interface{} {
	return &discord.UnmarshalInteraction{}
}

func (h *httpserverHandlerInteractionCreate) HandleHTTPEvent(bot core.Bot, c chan<- discord.InteractionResponse, v interface{}) {
	interaction := (*v.(*discord.UnmarshalInteraction)).Interaction

	// we just want to pong all pings
	// no need for any event
	if interaction.Type() == discord.InteractionTypePing {
		bot.Logger().Info("received http interaction ping. responding with pong")
		c <- discord.InteractionResponse{
			Type: discord.InteractionCallbackTypePong,
		}
		return
	}
	HandleInteraction(bot, -1, c, interaction)
}
