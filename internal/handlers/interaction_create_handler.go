package handlers

import (
	"encoding/json"

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
	interaction, ok := i.(*api.FullInteraction)
	if !ok {
		return
	}
	handleInteraction(disgo, eventManager, sequenceNumber, interaction, nil)
}

func handleInteraction(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, fullInteraction *api.FullInteraction, c chan *api.InteractionResponse) {

	interaction := disgo.EntityBuilder().CreateInteraction(fullInteraction.Interaction, api.CacheStrategyYes)
	genericInteractionEvent := events.GenericInteractionEvent{
		GenericEvent:    events.NewEvent(disgo, sequenceNumber),
		Interaction:     interaction,
		ResponseChannel: c,
		FromWebhook:     c != nil,
		Replied:         false,
	}
	eventManager.Dispatch(genericInteractionEvent)

	switch interaction.Type {
	case api.InteractionTypeApplicationCommand:
		var data *api.SlashCommandInteractionData
		if err := json.Unmarshal(fullInteraction.Data, &data); err != nil {
			disgo.Logger().Errorf("failed to unmarshal SlashCommandInteractionData: %s", err)
			return
		}
		options := data.Options
		var subCommandName *string
		var subCommandGroupName *string
		if len(options) == 1 {
			option := data.Options[0]
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
				Resolved: data.Resolved,
				Name:     optionData.Name,
				Type:     optionData.Type,
				Value:    optionData.Value,
			})
		}

		eventManager.Dispatch(events.SlashCommandEvent{
			GenericInteractionEvent: genericInteractionEvent,
			CommandID:               data.ID,
			CommandName:             data.Name,
			SubCommandName:          subCommandName,
			SubCommandGroupName:     subCommandGroupName,
			Options:                 newOptions,
		})
	case api.InteractionTypeComponent:
		var data *api.ComponentInteractionData
		if err := json.Unmarshal(fullInteraction.Data, &data); err != nil {
			disgo.Logger().Errorf("failed to unmarshal ComponentInteractionData: %s", err)
			return
		}
	}
}
