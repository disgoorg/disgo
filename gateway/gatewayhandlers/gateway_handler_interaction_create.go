package gatewayhandlers

import (
	"github.com/DisgoOrg/disgo/core"
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

	bot.EventManager.Dispatch(&events.InteractionCreateEvent{
		GenericEvent: genericEvent,
		Interaction:  coreInteraction,
	})

	switch i := coreInteraction.(type) {
	case core.ApplicationCommandInteraction:
		bot.EventManager.Dispatch(&events.ApplicationCommandInteractionCreateEvent{
			GenericEvent:                  genericEvent,
			ApplicationCommandInteraction: i,
		})

		switch ii := i.(type) {
		case *core.SlashCommandInteraction:
			bot.EventManager.Dispatch(&events.SlashCommandEvent{
				GenericEvent:            genericEvent,
				SlashCommandInteraction: ii,
			})

		case *core.UserCommandInteraction:
			bot.EventManager.Dispatch(&events.UserCommandEvent{
				GenericEvent:           genericEvent,
				UserCommandInteraction: ii,
			})

		case *core.MessageCommandInteraction:
			bot.EventManager.Dispatch(&events.MessageCommandEvent{
				GenericEvent:              genericEvent,
				MessageCommandInteraction: ii,
			})

		default:
			bot.Logger.Errorf("unknown application command interaction with type %d received", ii.ApplicationCommandType())
		}

	case core.ComponentInteraction:
		bot.EventManager.Dispatch(&events.ComponentInteractionCreateEvent{
			GenericEvent:         genericEvent,
			ComponentInteraction: i,
		})

		switch ii := i.(type) {
		case *core.ButtonInteraction:
			bot.EventManager.Dispatch(&events.ButtonClickEvent{
				GenericEvent:      genericEvent,
				ButtonInteraction: ii,
			})

		case *core.SelectMenuInteraction:
			bot.EventManager.Dispatch(&events.SelectMenuSubmitEvent{
				GenericEvent:          genericEvent,
				SelectMenuInteraction: ii,
			})

		default:
			bot.Logger.Errorf("unknown component interaction with type %d received", ii.ComponentType())
		}

	case *core.AutocompleteInteraction:
		bot.EventManager.Dispatch(&events.AutocompleteEvent{
			GenericEvent:            genericEvent,
			AutocompleteInteraction: i,
		})

	default:
		bot.Logger.Errorf("unknown interaction with type %d received", interaction.InteractionType())
	}
}
