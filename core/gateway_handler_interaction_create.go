package core

import (
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
	return &discord.Interaction{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerInteractionCreate) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	interaction := *v.(*discord.Interaction)

	HandleInteraction(bot, sequenceNumber, nil, interaction)
}

func HandleInteraction(bot *Bot, sequenceNumber int, c chan<- discord.InteractionResponse, interaction discord.Interaction) {
	coreInteraction := bot.EntityBuilder.CreateInteraction(interaction, c, CacheStrategyYes)

	genericEvent := NewGenericEvent(bot, sequenceNumber)

	switch interaction.Type {
	case discord.InteractionTypeApplicationCommandAutoComplete:
		bot.EventManager.Dispatch(&ApplicationCommandAutocompleteEvent{
			GenericEvent: genericEvent,
			ApplicationCommandAutocompleteInteraction: bot.EntityBuilder.CreateApplicationCommandAutocompleteInteraction(bot.EntityBuilder.CreateApplicationCommandOptionsInteraction(bot.EntityBuilder.CreateApplicationCommandInteraction(coreInteraction, CacheStrategyYes))),
		})

	case discord.InteractionTypeCommand:
		applicationCommandInteraction := bot.EntityBuilder.CreateApplicationCommandInteraction(coreInteraction, CacheStrategyYes)

		switch interaction.Data.CommandType {
		case discord.ApplicationCommandTypeSlash:
			bot.EventManager.Dispatch(&SlashCommandEvent{
				GenericEvent:            genericEvent,
				SlashCommandInteraction: bot.EntityBuilder.CreateSlashCommandInteraction(bot.EntityBuilder.CreateApplicationCommandOptionsInteraction(applicationCommandInteraction)),
			})

		case discord.ApplicationCommandTypeUser, discord.ApplicationCommandTypeMessage:
			contextCommandInteraction := bot.EntityBuilder.CreateContextCommandInteraction(applicationCommandInteraction)

			switch interaction.Data.CommandType {
			case discord.ApplicationCommandTypeUser:
				bot.EventManager.Dispatch(&UserCommandEvent{
					GenericEvent:           genericEvent,
					UserCommandInteraction: bot.EntityBuilder.CreateUserCommandInteraction(contextCommandInteraction),
				})

			case discord.ApplicationCommandTypeMessage:
				bot.EventManager.Dispatch(&MessageCommandEvent{
					GenericEvent:              genericEvent,
					MessageCommandInteraction: bot.EntityBuilder.CreateMessageCommandInteraction(contextCommandInteraction),
				})
			}
		}

	case discord.InteractionTypeComponent:
		componentInteraction := bot.EntityBuilder.CreateComponentInteraction(coreInteraction, CacheStrategyYes)

		switch interaction.Data.ComponentType {
		case discord.ComponentTypeButton:
			bot.EventManager.Dispatch(&ButtonClickEvent{
				GenericEvent:      genericEvent,
				ButtonInteraction: bot.EntityBuilder.CreateButtonInteraction(componentInteraction),
			})

		case discord.ComponentTypeSelectMenu:
			bot.EventManager.Dispatch(&SelectMenuSubmitEvent{
				GenericEvent:          genericEvent,
				SelectMenuInteraction: bot.EntityBuilder.CreateSelectMenuInteraction(componentInteraction),
			})
		}
	}
}
