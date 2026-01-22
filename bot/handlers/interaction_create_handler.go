package handlers

import (
	"fmt"
	"log/slog"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/httpserver"
	"github.com/disgoorg/disgo/rest"
)

func gatewayHandlerInteractionCreate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventInteractionCreate) {
	handleInteraction(client, sequenceNumber, shardID, nil, event.Interaction)
}

func respond(client *bot.Client, respondFunc httpserver.RespondFunc, interaction discord.Interaction) events.InteractionResponderFunc {
	return func(responseType discord.InteractionResponseType, data discord.InteractionResponseData, opts ...rest.RequestOpt) error {
		response := discord.InteractionResponse{
			Type: responseType,
			Data: data,
		}
		if respondFunc != nil {
			return respondFunc(response)
		}
		return client.Rest.CreateInteractionResponse(interaction.ID(), interaction.Token(), response, opts...)
	}
}

func handleInteraction(client *bot.Client, sequenceNumber int, shardID int, respondFunc httpserver.RespondFunc, interaction discord.Interaction) {
	genericEvent := events.NewGenericEvent(client, sequenceNumber, shardID)
	responseState := &events.InteractionResponseState{}
	wrappedResponder := events.WrapInteractionResponder(respond(client, respondFunc, interaction), responseState)

	client.EventManager.DispatchEvent(&events.InteractionCreate{
		GenericEvent:  genericEvent,
		Interaction:   interaction,
		ResponseState: responseState,
		Respond:       wrappedResponder,
	})

	switch i := interaction.(type) {
	case discord.ApplicationCommandInteraction:
		client.EventManager.DispatchEvent(&events.ApplicationCommandInteractionCreate{
			GenericEvent:                  genericEvent,
			ApplicationCommandInteraction: i,
			ResponseState:                 responseState,
			Respond:                       wrappedResponder,
		})

	case discord.ComponentInteraction:
		client.EventManager.DispatchEvent(&events.ComponentInteractionCreate{
			GenericEvent:         genericEvent,
			ComponentInteraction: i,
			ResponseState:        responseState,
			Respond:              wrappedResponder,
		})

	case discord.AutocompleteInteraction:
		client.EventManager.DispatchEvent(&events.AutocompleteInteractionCreate{
			GenericEvent:            genericEvent,
			AutocompleteInteraction: i,
			ResponseState:           responseState,
			Respond:                 wrappedResponder,
		})

	case discord.ModalSubmitInteraction:
		client.EventManager.DispatchEvent(&events.ModalSubmitInteractionCreate{
			GenericEvent:           genericEvent,
			ModalSubmitInteraction: i,
			ResponseState:          responseState,
			Respond:                wrappedResponder,
		})

	default:
		client.Logger.Error("unknown interaction", slog.String("type", fmt.Sprintf("%T", interaction)))
	}
}
