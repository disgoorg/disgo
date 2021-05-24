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
	return &api.Interaction{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h InteractionCreateHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	interaction, ok := i.(*api.Interaction)
	if !ok {
		return
	}
	handleInteraction(disgo, eventManager, sequenceNumber, interaction, nil)
}

func handleInteraction(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, interaction *api.Interaction, c chan *api.InteractionResponse) {

	interaction = disgo.EntityBuilder().CreateInteraction(interaction, api.CacheStrategyYes)
	genericInteractionEvent := events.GenericInteractionEvent{
		GenericEvent: events.NewEvent(disgo, sequenceNumber),
		Interaction:  interaction,
	}
	eventManager.Dispatch(genericInteractionEvent)

	if interaction.Data != nil {
		options := interaction.Data.Options
		var subCommandName *string
		var subCommandGroupName *string
		if len(options) == 1 {
			option := interaction.Data.Options[0]
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

		eventManager.Dispatch(events.SlashCommandEvent{
			GenericInteractionEvent: genericInteractionEvent,
			ResponseChannel:         c,
			FromWebhook:             c != nil,
			CommandID:               interaction.Data.ID,
			CommandName:             interaction.Data.Name,
			SubCommandName:          subCommandName,
			SubCommandGroupName:     subCommandGroupName,
			Options:                 newOptions,
			Replied:                 false,
		})
	}
}
