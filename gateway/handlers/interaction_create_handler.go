package handlers

import (
	"github.com/DisgoOrg/disgo/bot"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
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
func (h *gatewayHandlerInteractionCreate) HandleGatewayEvent(client bot.Client, sequenceNumber discord.GatewaySequence, v any) {
	HandleInteraction(client, sequenceNumber, nil, (*v.(*discord.UnmarshalInteraction)).Interaction)
}

func respond(client bot.Client, c chan<- discord.InteractionResponse, interaction discord.BaseInteraction, response discord.InteractionResponse) error {
	if c != nil {
		c <- response
		return nil
	}
	return client.Rest().InteractionService().CreateInteractionResponse(interaction.ID(), interaction.Token(), response)
}

func HandleInteraction(client bot.Client, sequenceNumber discord.GatewaySequence, c chan<- discord.InteractionResponse, interaction discord.Interaction) {

	genericEvent := events.NewGenericEvent(client, sequenceNumber)

	client.EventManager().Dispatch(&events.InteractionEvent{
		GenericEvent: genericEvent,
		Interaction:  interaction,
		Respond: func(callbackType discord.InteractionCallbackType, data discord.InteractionCallbackData) error {
			return respond(client, c, interaction, discord.InteractionResponse{
				Type: callbackType,
				Data: data,
			})
		}})

	switch i := interaction.(type) {
	case discord.ApplicationCommandInteraction:
		client.EventManager().Dispatch(&events.ApplicationCommandInteractionEvent{
			GenericEvent:                  genericEvent,
			ApplicationCommandInteraction: i,
			Respond: func(callbackType discord.InteractionCallbackType, data discord.CommandInteractionCallbackData) error {
				return respond(client, c, interaction, discord.InteractionResponse{
					Type: callbackType,
					Data: data,
				})
			},
		})

	case discord.ComponentInteraction:
		client.EventManager().Dispatch(&events.ComponentInteractionEvent{
			GenericEvent:         genericEvent,
			ComponentInteraction: i,
			Respond: func(callbackType discord.InteractionCallbackType, data discord.ComponentInteractionCallbackData) error {
				return respond(client, c, interaction, discord.InteractionResponse{
					Type: callbackType,
					Data: data,
				})
			},
		})

	case discord.AutocompleteInteraction:
		client.EventManager().Dispatch(&events.AutocompleteInteractionEvent{
			GenericEvent:            genericEvent,
			AutocompleteInteraction: i,
			Respond: func(data discord.AutocompleteResult) error {
				return respond(client, c, interaction, discord.InteractionResponse{
					Type: discord.InteractionCallbackTypeAutocompleteResult,
					Data: data,
				})
			},
		})

	case discord.ModalSubmitInteraction:
		client.EventManager().Dispatch(&events.ModalSubmitInteractionEvent{
			GenericEvent:           genericEvent,
			ModalSubmitInteraction: i,
			Respond: func(callbackType discord.InteractionCallbackType, data discord.ModalInteractionCallbackData) error {
				return respond(client, c, interaction, discord.InteractionResponse{
					Type: callbackType,
					Data: data,
				})
			},
		})

	default:
		client.Logger().Errorf("unknown interaction with type %d received", interaction.Type())
	}
}
