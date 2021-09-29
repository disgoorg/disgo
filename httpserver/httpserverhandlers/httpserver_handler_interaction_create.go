package httpserverhandlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway/gatewayhandlers"
)

var _ core.HTTPServerEventHandler = (*httpserverHandlerInteractionCreate)(nil)

type httpserverHandlerInteractionCreate struct{}

func (h *httpserverHandlerInteractionCreate) New() interface{} {
	return &discord.Interaction{}
}

func (h *httpserverHandlerInteractionCreate) HandleHTTPEvent(bot *core.Bot, c chan<- discord.InteractionResponse, v interface{}) {
	unmarshalInteraction := *v.(*discord.Interaction)

	// we just want to pong all pings
	// no need for any event
	if unmarshalInteraction.Type == discord.InteractionTypePing {
		bot.Logger.Info("received http interaction ping. responding with pong")
		c <- discord.InteractionResponse{
			Type: discord.InteractionCallbackTypePong,
		}
		return
	}
	gatewayhandlers.HandleInteraction(bot, -1, c, unmarshalInteraction)
}
