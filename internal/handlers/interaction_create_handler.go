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

func handleInteraction(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, fullInteraction *api.FullInteraction, c chan *api.InteractionResponse) {
	genericInteractionEvent := events.GenericInteractionEvent{
		GenericEvent: events.NewEvent(disgo, sequenceNumber),
	}

	switch fullInteraction.Type {
	case api.InteractionTypeCommand:
		interaction := disgo.EntityBuilder().CreateCommandInteraction(fullInteraction, c, api.CacheStrategyYes)

		genericInteractionEvent.Interaction = interaction.Interaction
		eventManager.Dispatch(genericInteractionEvent)

		options := interaction.Data.Options
		var subCommandName *string
		var subCommandGroupName *string
		if len(options) == 1 {
			option := options[0]
			if option.Type == api.CommandOptionTypeSubCommandGroup {
				subCommandGroupName = &option.Name
				options = option.Options
				option = option.Options[0]
			}
			if option.Type == api.CommandOptionTypeSubCommand {
				subCommandName = &option.Name
				options = option.Options
			}
		}

		var newOptions []*api.Option
		for _, optionData := range options {
			newOptions = append(newOptions, &api.Option{
				Resolved: interaction.Data.Resolved,
				Name:     optionData.Name,
				Type:     optionData.Type,
				Value:    optionData.Value,
			})
		}

		eventManager.Dispatch(events.CommandEvent{
			GenericInteractionEvent: genericInteractionEvent,
			CommandInteraction:      interaction,
			CommandID:               interaction.Data.ID,
			CommandName:             interaction.Data.Name,
			SubCommandName:          subCommandName,
			SubCommandGroupName:     subCommandGroupName,
			Options:                 newOptions,
		})

	case api.InteractionTypeComponent:
		componentInteraction := disgo.EntityBuilder().CreateComponentInteraction(fullInteraction, c, api.CacheStrategyYes)

		genericInteractionEvent.Interaction = componentInteraction.Interaction
		eventManager.Dispatch(genericInteractionEvent)

		genericComponentEvent := events.GenericComponentEvent{
			GenericInteractionEvent: genericInteractionEvent,
			ComponentInteraction:    componentInteraction,
		}
		eventManager.Dispatch(genericComponentEvent)

		switch componentInteraction.Data.ComponentType {
		case api.ComponentTypeButton:
			eventManager.Dispatch(events.ButtonClickEvent{
				GenericComponentEvent: genericComponentEvent,
				ButtonInteraction:     disgo.EntityBuilder().CreateButtonInteraction(fullInteraction, componentInteraction),
			})

		case api.ComponentTypeDropdown:
			eventManager.Dispatch(events.DropdownSubmitEvent{
				GenericComponentEvent: genericComponentEvent,
				DropdownInteraction:   disgo.EntityBuilder().CreateDropdownInteraction(fullInteraction, componentInteraction),
			})
		}

	}
}
