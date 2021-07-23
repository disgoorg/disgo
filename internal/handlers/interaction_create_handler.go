package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

// InteractionCreateHandler handles api.InteractionCreateGatewayEvent
type InteractionCreateHandler struct{}

// Event returns the raw gateway event Event
func (h InteractionCreateHandler) Event() api.GatewayEventType {
	return api.GatewayEventInteractionCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h InteractionCreateHandler) New() interface{} {
	return &api.FullInteraction{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h InteractionCreateHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	fullInteraction, ok := i.(*api.FullInteraction)
	if !ok {
		return
	}
	handleInteraction(disgo, eventManager, sequenceNumber, fullInteraction, nil)
}

func handleInteraction(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, fullInteraction *api.FullInteraction, c chan api.InteractionResponse) {
	interaction := disgo.EntityBuilder().CreateInteraction(fullInteraction, c, api.CacheStrategyYes)

	genericInteractionEvent := &events.GenericInteractionEvent{
		GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
		Interaction:  interaction,
	}

	switch fullInteraction.Type {
	case api.InteractionTypeCommand:
		genericCommandInteraction := disgo.EntityBuilder().CreateGenericCommandInteraction(fullInteraction, interaction, api.CacheStrategyYes)

		genericCommandEvent := &events.GenericCommandEvent{
			GenericInteractionEvent:   genericInteractionEvent,
			GenericCommandInteraction: genericCommandInteraction,
		}

		switch genericCommandInteraction.Data.Type {
		case api.CommandTypeSlashCommand:
			eventManager.Dispatch(&events.SlashCommandEvent{
				GenericCommandEvent: genericCommandEvent,
				CommandInteraction:  disgo.EntityBuilder().CreateSlashCommandInteraction(fullInteraction, genericCommandInteraction),
			})

		case api.CommandTypeUserContext, api.CommandTypeMessageContext:
			genericContextInteraction := disgo.EntityBuilder().CreateGenericContextInteraction(fullInteraction, genericCommandInteraction)

			genericContextEvent := &events.GenericContextEvent{
				GenericCommandEvent:       genericCommandEvent,
				GenericContextInteraction: genericContextInteraction,
			}

			switch genericCommandInteraction.Data.Type {
			case api.CommandTypeUserContext:
				eventManager.Dispatch(&events.UserContextEvent{
					GenericContextEvent:    genericContextEvent,
					UserContextInteraction: disgo.EntityBuilder().CreateUserContextInteraction(fullInteraction, genericContextInteraction),
				})

			case api.CommandTypeMessageContext:
				eventManager.Dispatch(&events.MessageContextEvent{
					GenericContextEvent:       genericContextEvent,
					MessageContextInteraction: disgo.EntityBuilder().CreateMessageContextInteraction(fullInteraction, genericContextInteraction),
				})
			}
		}

	case api.InteractionTypeComponent:
		componentInteraction := disgo.EntityBuilder().CreateComponentInteraction(fullInteraction, interaction, api.CacheStrategyYes)

		genericComponentEvent := &events.GenericComponentEvent{
			GenericInteractionEvent: genericInteractionEvent,
			ComponentInteraction:    componentInteraction,
		}

		switch componentInteraction.Data.ComponentType {
		case api.ComponentTypeButton:
			eventManager.Dispatch(&events.ButtonClickEvent{
				GenericComponentEvent: genericComponentEvent,
				ButtonInteraction:     disgo.EntityBuilder().CreateButtonInteraction(fullInteraction, componentInteraction),
			})

		case api.ComponentTypeSelectMenu:
			eventManager.Dispatch(&events.SelectMenuSubmitEvent{
				GenericComponentEvent: genericComponentEvent,
				SelectMenuInteraction: disgo.EntityBuilder().CreateSelectMenuInteraction(fullInteraction, componentInteraction),
			})
		}
	}
}
