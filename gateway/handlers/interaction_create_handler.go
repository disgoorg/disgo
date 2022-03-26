package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/rest"
)

// gatewayHandlerInteractionCreate handles core.InteractionCreateGatewayEvent
type gatewayHandlerInteractionCreate struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *gatewayHandlerInteractionCreate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeInteractionCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerInteractionCreate) New() any {
	return &discord.UnmarshalInteraction{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerInteractionCreate) HandleGatewayEvent(client bot.Client, sequenceNumber int, v any) {
	HandleInteraction(client, sequenceNumber, nil, (*v.(*discord.UnmarshalInteraction)).Interaction)
}

func respond(client bot.Client, c chan<- discord.InteractionResponse, interaction discord.BaseInteraction) events.InteractionResponderFunc {
	return func(callbackType discord.InteractionCallbackType, data discord.InteractionCallbackData, opts ...rest.RequestOpt) error {
		response := discord.InteractionResponse{
			Type: callbackType,
			Data: data,
		}
		if c != nil {
			c <- response
			return nil
		}
		return client.Rest().Interactions().CreateInteractionResponse(interaction.ID(), interaction.Token(), response, opts...)
	}
}

func HandleInteraction(client bot.Client, sequenceNumber int, c chan<- discord.InteractionResponse, interaction discord.Interaction) {

	genericEvent := events.NewGenericEvent(client, sequenceNumber)

	client.EventManager().DispatchEvent(&events.InteractionEvent{
		GenericEvent: genericEvent,
		Interaction:  interaction,
		Respond:      respond(client, c, interaction),
	})

	switch i := interaction.(type) {
	case discord.ApplicationCommandInteraction:
		client.EventManager().DispatchEvent(&events.ApplicationCommandInteractionEvent{
			GenericEvent:                  genericEvent,
			ApplicationCommandInteraction: i,
			Respond:                       respond(client, c, interaction),
		})

	case discord.ComponentInteraction:
		client.EventManager().DispatchEvent(&events.ComponentInteractionEvent{
			GenericEvent:         genericEvent,
			ComponentInteraction: i,
			Respond:              respond(client, c, interaction),
		})

	case discord.AutocompleteInteraction:
		client.EventManager().DispatchEvent(&events.AutocompleteInteractionEvent{
			GenericEvent:            genericEvent,
			AutocompleteInteraction: i,
			Respond:                 respond(client, c, interaction),
		})

	case discord.ModalSubmitInteraction:
		client.EventManager().DispatchEvent(&events.ModalSubmitInteractionEvent{
			GenericEvent:           genericEvent,
			ModalSubmitInteraction: i,
			Respond:                respond(client, c, interaction),
		})

	default:
		client.Logger().Errorf("unknown interaction with type %d received", interaction.Type())
	}
}
