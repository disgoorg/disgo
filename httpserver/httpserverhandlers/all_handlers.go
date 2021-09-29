package httpserverhandlers

import (
	"io"

	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
	"github.com/DisgoOrg/disgo/httpserver"
)

func DefaultHTTPServerEventHandler(bot *core.Bot) httpserver.EventHandlerFunc {
	return func(responseChannel chan<- discord.InteractionResponse, reader io.Reader) {
		reader = events.HandleRawEvent(bot, discord.GatewayEventTypeInteractionCreate, -1, reader)

		bot.EventManager.HandleHTTP(responseChannel, reader)
	}
}

func GetHTTPServerHandler() core.HTTPServerEventHandler {
	return &httpserverHandlerInteractionCreate{}
}
