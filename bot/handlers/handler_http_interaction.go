package handlers

import (
	"log/slog"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/httpinteraction"
)

// GetHTTPInteractionHandler returns the default httpserver.Server event handler for processing the raw payload which gets passed into the bot.EventManager
func GetHTTPInteractionHandler() bot.HTTPInteractionEventHandler {
	return &httpInteractionHandler{}
}

type httpInteractionHandler struct{}

func (h *httpInteractionHandler) HandleHTTPInteraction(client *bot.Client, respond httpinteraction.RespondFunc, event httpinteraction.EventInteractionCreate) {
	// we just want to pong all pings automatically
	if event.Type() == discord.InteractionTypePing {
		client.Logger.Debug("received http interaction ping. responding with pong")
		if err := respond(discord.InteractionResponse{
			Type: discord.InteractionResponseTypePong,
		}); err != nil {
			client.Logger.Error("failed to respond to http interaction ping", slog.Any("err", err))
		}
	}
	handleInteraction(client, -1, -1, respond, event.Interaction)
}
