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
func (h *gatewayHandlerInteractionCreate) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	HandleInteraction(bot, sequenceNumber, nil, (*v.(*discord.UnmarshalInteraction)).Interaction)
}

func HandleInteraction(bot *core.Bot, sequenceNumber int, c chan<- discord.InteractionResponse, interaction discord.Interaction) {
	coreInteraction := bot.EntityBuilder.CreateInteraction(interaction, c, core.CacheStrategyYes)

	genericEvent := events.NewGenericEvent(bot, sequenceNumber)

	bot.EventManager.Dispatch(&events.InteractionEvent{
		GenericEvent: genericEvent,
		Interaction:  coreInteraction,
	})

	switch i := coreInteraction.(type) {
	case *core.ApplicationCommandInteraction:
		bot.EventManager.Dispatch(&events.ApplicationCommandInteractionEvent{
			GenericEvent:                  genericEvent,
			ApplicationCommandInteraction: i,
		})

	case *core.ComponentInteraction:
		bot.EventManager.Dispatch(&events.ComponentInteractionEvent{
			GenericEvent:         genericEvent,
			ComponentInteraction: i,
		})

	case *core.AutocompleteInteraction:
		bot.EventManager.Dispatch(&events.AutocompleteInteractionEvent{
			GenericEvent:            genericEvent,
			AutocompleteInteraction: i,
		})

	case *core.ModalSubmitInteraction:
		bot.EventManager.Dispatch(&events.ModalSubmitInteractionEvent{
			GenericEvent:           genericEvent,
			ModalSubmitInteraction: i,
		})

	default:
		bot.Logger.Errorf("unknown interaction with type %d received", interaction.Type())
	}
}
