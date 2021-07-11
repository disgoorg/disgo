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
		commandInteraction := disgo.EntityBuilder().CreateCommandInteraction(fullInteraction, interaction, api.CacheStrategyYes)

		options := commandInteraction.Data.RawOptions
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
				Resolved: commandInteraction.Data.Resolved,
				Name:     optionData.Name,
				Type:     optionData.Type,
				Value:    optionData.Value,
			})
		}

		commandInteraction.Data.Options = newOptions
		commandInteraction.Data.SubCommandName = subCommandName
		commandInteraction.Data.SubCommandGroupName = subCommandGroupName

		eventManager.Dispatch(&events.CommandEvent{
			GenericInteractionEvent: genericInteractionEvent,
			CommandInteraction:      commandInteraction,
		})

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
