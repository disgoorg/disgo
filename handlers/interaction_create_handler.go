package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/rest"
)

type gatewayHandlerInteractionCreate struct{}

func (h *gatewayHandlerInteractionCreate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeInteractionCreate
}

func (h *gatewayHandlerInteractionCreate) New() any {
	return &discord.UnmarshalInteraction{}
}

func (h *gatewayHandlerInteractionCreate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	handleInteraction(client, sequenceNumber, shardID, nil, (*v.(*discord.UnmarshalInteraction)).Interaction)
}

func respond(client bot.Client, respondFunc func(response discord.InteractionResponse) error, interaction discord.BaseInteraction) events.InteractionResponderFunc {
	return func(responseType discord.InteractionResponseType, data discord.InteractionResponseData, opts ...rest.RequestOpt) error {
		response := discord.InteractionResponse{
			Type: responseType,
			Data: data,
		}
		if respondFunc != nil {
			return respondFunc(response)
		}
		return client.Rest().CreateInteractionResponse(interaction.ID(), interaction.Token(), response, opts...)
	}
}

func handleInteraction(client bot.Client, sequenceNumber int, shardID int, respondFunc func(response discord.InteractionResponse) error, interaction discord.Interaction) {

	genericEvent := events.NewGenericEvent(client, sequenceNumber, shardID)

	client.EventManager().DispatchEvent(&events.InteractionCreate{
		GenericEvent: genericEvent,
		Interaction:  interaction,
		Respond:      respond(client, respondFunc, interaction),
	})

	switch i := interaction.(type) {
	case discord.ApplicationCommandInteraction:
		client.EventManager().DispatchEvent(&events.ApplicationCommandInteractionCreate{
			GenericEvent:                  genericEvent,
			ApplicationCommandInteraction: i,
			Respond:                       respond(client, respondFunc, interaction),
		})

	case discord.ComponentInteraction:
		client.EventManager().DispatchEvent(&events.ComponentInteractionCreate{
			GenericEvent:         genericEvent,
			ComponentInteraction: i,
			Respond:              respond(client, respondFunc, interaction),
		})

	case discord.AutocompleteInteraction:
		client.EventManager().DispatchEvent(&events.AutocompleteInteractionCreate{
			GenericEvent:            genericEvent,
			AutocompleteInteraction: i,
			Respond:                 respond(client, respondFunc, interaction),
		})

	case discord.ModalSubmitInteraction:
		client.EventManager().DispatchEvent(&events.ModalSubmitInteractionCreate{
			GenericEvent:           genericEvent,
			ModalSubmitInteraction: i,
			Respond:                respond(client, respondFunc, interaction),
		})

	default:
		client.Logger().Errorf("unknown interaction with type %d received", interaction.Type())
	}
}
