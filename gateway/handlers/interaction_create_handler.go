package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// InteractionCreateHandler handles api.InteractionCreateGatewayEvent
type InteractionCreateHandler struct{}

// EventType returns the api.GatewayGatewayEventType
func (h *InteractionCreateHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeInteractionCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *InteractionCreateHandler) New() interface{} {
	return discord.Interaction{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *InteractionCreateHandler) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	interaction, ok := v.(discord.Interaction)
	if !ok {
		return
	}
	HandleInteraction(bot, sequenceNumber, nil, interaction)
}

func HandleInteraction(bot *core.Bot, sequenceNumber int, c chan discord.InteractionResponse, interaction discord.Interaction) {
	coreInteraction := bot.EntityBuilder.CreateInteraction(interaction, c, core.CacheStrategyYes)

	genericEvent := events.NewGenericEvent(bot, sequenceNumber)

	switch interaction.Type {
	case discord.InteractionTypeCommand:
		applicationCommandInteraction := bot.EntityBuilder.CreateApplicationCommandInteraction(coreInteraction, core.CacheStrategyYes)

		switch interaction.Data.CommandType {
		case discord.ApplicationCommandTypeSlash:
			bot.EventManager.Dispatch(&events.SlashCommandEvent{
				GenericEvent:            genericEvent,
				SlashCommandInteraction: bot.EntityBuilder.CreateSlashCommandInteraction(applicationCommandInteraction),
			})

		case discord.ApplicationCommandTypeUser, discord.ApplicationCommandTypeMessage:
			contextCommandInteraction := bot.EntityBuilder.CreateContextCommandInteraction(applicationCommandInteraction)

			switch interaction.Data.CommandType {
			case discord.ApplicationCommandTypeUser:
				bot.EventManager.Dispatch(&events.UserCommandEvent{
					GenericEvent:           genericEvent,
					UserCommandInteraction: bot.EntityBuilder.CreateUserCommandInteraction(contextCommandInteraction),
				})

			case discord.ApplicationCommandTypeMessage:
				bot.EventManager.Dispatch(&events.MessageCommandEvent{
					GenericEvent:              genericEvent,
					MessageCommandInteraction: bot.EntityBuilder.CreateMessageCommandInteraction(contextCommandInteraction),
				})
			}
		}

	case discord.InteractionTypeComponent:
		componentInteraction := bot.EntityBuilder.CreateComponentInteraction(coreInteraction, core.CacheStrategyYes)

		switch interaction.Data.ComponentType {
		case discord.ComponentTypeButton:
			bot.EventManager.Dispatch(&events.ButtonClickEvent{
				GenericEvent:      genericEvent,
				ButtonInteraction: bot.EntityBuilder.CreateButtonInteraction(componentInteraction),
			})

		case discord.ComponentTypeSelectMenu:
			bot.EventManager.Dispatch(&events.SelectMenuSubmitEvent{
				GenericEvent:          genericEvent,
				SelectMenuInteraction: bot.EntityBuilder.CreateSelectMenuInteraction(componentInteraction),
			})
		}
	}
}
