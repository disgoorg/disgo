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
	responseState := events.NewInteractionResponseState()
	wrappedResponder := events.WrapInteractionResponder(respond(client, respondFunc, interaction), responseState)

	interactionEvent := &events.InteractionCreate{
		GenericEvent: genericEvent,
		Interaction:  interaction,
		Respond:      wrappedResponder,
	}
	events.SetInteractionCreateResponseState(interactionEvent, responseState)
	client.EventManager.DispatchEvent(interactionEvent)

	switch i := interaction.(type) {
	case discord.ApplicationCommandInteraction:
		appEvent := &events.ApplicationCommandInteractionCreate{
			GenericEvent:                  genericEvent,
			ApplicationCommandInteraction: i,
			Respond:                       wrappedResponder,
		}
		events.SetApplicationCommandInteractionCreateResponseState(appEvent, responseState)
		client.EventManager.DispatchEvent(appEvent)

	case discord.ComponentInteraction:
		componentEvent := &events.ComponentInteractionCreate{
			GenericEvent:         genericEvent,
			ComponentInteraction: i,
			Respond:              wrappedResponder,
		}
		events.SetComponentInteractionCreateResponseState(componentEvent, responseState)
		client.EventManager.DispatchEvent(componentEvent)

	case discord.AutocompleteInteraction:
		autocompleteEvent := &events.AutocompleteInteractionCreate{
			GenericEvent:            genericEvent,
			AutocompleteInteraction: i,
			Respond:                 wrappedResponder,
		}
		events.SetAutocompleteInteractionCreateResponseState(autocompleteEvent, responseState)
		client.EventManager.DispatchEvent(autocompleteEvent)

	case discord.ModalSubmitInteraction:
		modalEvent := &events.ModalSubmitInteractionCreate{
			GenericEvent:           genericEvent,
			ModalSubmitInteraction: i,
			Respond:                wrappedResponder,
		}
		events.SetModalSubmitInteractionCreateResponseState(modalEvent, responseState)
		client.EventManager.DispatchEvent(modalEvent)

	default:
		client.Logger.Error("unknown interaction", slog.String("type", fmt.Sprintf("%T", interaction)))
	}
}
