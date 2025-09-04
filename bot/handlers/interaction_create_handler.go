package handlers

import (
	"fmt"
	"log/slog"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/httpinteraction"
	"github.com/disgoorg/disgo/rest"
)

func gatewayHandlerInteractionCreate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventInteractionCreate) {
	handleInteraction(client, sequenceNumber, shardID, nil, event.Interaction)
}

func respondFunc(client *bot.Client, respond httpinteraction.RespondFunc, interaction discord.Interaction) events.InteractionResponderFunc {
	return func(responseType discord.InteractionResponseType, data discord.InteractionResponseData, opts ...rest.RequestOpt) error {
		response := discord.InteractionResponse{
			Type: responseType,
			Data: data,
		}
		if respond != nil {
			return respond(response)
		}
		return client.Rest.CreateInteractionResponse(interaction.ID(), interaction.Token(), response, opts...)
	}
}

func handleInteraction(client *bot.Client, sequenceNumber int, shardID int, respond httpinteraction.RespondFunc, interaction discord.Interaction) {
	client.EventManager.DispatchEvent(&events.InteractionCreate{
		Event:        events.NewEvent(client),
		GatewayEvent: events.NewGatewayEvent(sequenceNumber, shardID),
		Interaction:  interaction,
		Respond:      respondFunc(client, respond, interaction),
	})

	switch i := interaction.(type) {
	case discord.ApplicationCommandInteraction:
		client.EventManager.DispatchEvent(&events.ApplicationCommandInteractionCreate{
			Event:                         events.NewEvent(client),
			GatewayEvent:                  events.NewGatewayEvent(sequenceNumber, shardID),
			ApplicationCommandInteraction: i,
			Respond:                       respondFunc(client, respond, interaction),
		})

	case discord.ComponentInteraction:
		client.EventManager.DispatchEvent(&events.ComponentInteractionCreate{
			Event:                events.NewEvent(client),
			GatewayEvent:         events.NewGatewayEvent(sequenceNumber, shardID),
			ComponentInteraction: i,
			Respond:              respondFunc(client, respond, interaction),
		})

	case discord.AutocompleteInteraction:
		client.EventManager.DispatchEvent(&events.AutocompleteInteractionCreate{
			Event:                   events.NewEvent(client),
			GatewayEvent:            events.NewGatewayEvent(sequenceNumber, shardID),
			AutocompleteInteraction: i,
			Respond:                 respondFunc(client, respond, interaction),
		})

	case discord.ModalSubmitInteraction:
		client.EventManager.DispatchEvent(&events.ModalSubmitInteractionCreate{
			Event:                  events.NewEvent(client),
			GatewayEvent:           events.NewGatewayEvent(sequenceNumber, shardID),
			ModalSubmitInteraction: i,
			Respond:                respondFunc(client, respond, interaction),
		})

	default:
		client.Logger.Error("unknown interaction", slog.String("type", fmt.Sprintf("%T", interaction)))
	}
}
