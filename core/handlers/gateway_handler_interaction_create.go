package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerInteractionCreate handles core.InteractionCreateGatewayEvent
type gatewayHandlerInteractionCreate struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *gatewayHandlerInteractionCreate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeInteractionCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerInteractionCreate) New() interface{} {
	return &discord.UnmarshalInteraction{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerInteractionCreate) HandleGatewayEvent(bot core.Bot, sequenceNumber discord.GatewaySequence, v interface{}) {
	HandleInteraction(bot, sequenceNumber, nil, (*v.(*discord.UnmarshalInteraction)).Interaction)
}

func HandleInteraction(bot core.Bot, sequenceNumber discord.GatewaySequence, c chan<- discord.InteractionResponse, interaction discord.Interaction) {

	genericEvent := events.NewGenericEvent(bot, sequenceNumber)

	bot.EventManager().Dispatch(&events.InteractionEvent{
		GenericEvent: genericEvent,
		Interaction:  interaction,
	})

	switch i := interaction.(type) {
	case discord.ApplicationCommandInteraction:
		bot.EventManager().Dispatch(&events.ApplicationCommandInteractionEvent{
			GenericEvent:                  genericEvent,
			ApplicationCommandInteraction: i,
		})

	case discord.ComponentInteraction:
		bot.EventManager().Dispatch(&events.ComponentInteractionEvent{
			GenericEvent:         genericEvent,
			ComponentInteraction: i,
		})

	case discord.AutocompleteInteraction:
		bot.EventManager().Dispatch(&events.AutocompleteInteractionEvent{
			GenericEvent:            genericEvent,
			AutocompleteInteraction: i,
		})

	case discord.ModalSubmitInteraction:
		bot.EventManager().Dispatch(&events.ModalSubmitInteractionEvent{
			GenericEvent:           genericEvent,
			ModalSubmitInteraction: i,
		})

	default:
		bot.Logger().Errorf("unknown interaction with type %d received", interaction.Type())
	}
}
