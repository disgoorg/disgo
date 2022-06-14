package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
)

var _ bot.HTTPServerEventHandler = (*httpserverHandlerInteractionCreate)(nil)

type httpserverHandlerInteractionCreate struct{}

func (h *httpserverHandlerInteractionCreate) New() any {
	return &discord.UnmarshalInteraction{}
}

func (h *httpserverHandlerInteractionCreate) HandleHTTPEvent(client bot.Client, respondFunc func(response discord.InteractionResponse) error, v any) {
	interaction := (*v.(*discord.UnmarshalInteraction)).Interaction

	// we just want to pong all pings
	// no need for any event
	if interaction.Type() == discord.InteractionTypePing {
		client.Logger().Debug("received http interaction ping. responding with pong")
		if err := respondFunc(discord.InteractionResponse{
			Type: discord.InteractionResponseTypePong,
		}); err != nil {
			client.Logger().Error("failed to respond to http interaction ping: ", err)
		}
		return
	}
	handleInteraction(client, -1, -1, respondFunc, interaction)
}
