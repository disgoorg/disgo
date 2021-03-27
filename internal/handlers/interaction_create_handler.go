package handlers

import (
	"github.com/DiscoOrg/disgo/api"
	"github.com/DiscoOrg/disgo/api/events"
)

type InteractionCreateHandler struct{}

// Name returns the raw gateway event name
func (h InteractionCreateHandler) Name() string {
	return api.InteractionCreateGatewayEvent
}

// New constructs a new payload receiver for the raw gateway event
func (h InteractionCreateHandler) New() interface{} {
	return &api.Interaction{}
}

// Handle handles the specific raw gateway event
func (h InteractionCreateHandler) Handle(disgo api.Disgo, eventManager api.EventManager, i interface{}) {
	interaction, ok := i.(*api.Interaction)
	if !ok {
		return
	}
	if interaction.Member != nil {
		disgo.Cache().CacheMember(interaction.Member)
	}
	if interaction.User != nil {
		disgo.Cache().CacheUser(interaction.User)
	}

	if interaction.Data != nil && interaction.Data.Resolved != nil {
		resolved := interaction.Data.Resolved
		if resolved.Users != nil {
			for _, user := range resolved.Users {
				disgo.Cache().CacheUser(user)
			}
		}
		if resolved.Members != nil {
			for id, member := range resolved.Members {
				member.User = resolved.Users[id]
				disgo.Cache().CacheMember(member)
			}
		}
		if resolved.Roles != nil {
			for _, role := range resolved.Roles {
				disgo.Cache().CacheRole(role)
			}
		}
		/*if resolved.Channels != nil {
			for _, user := range resolved.Users {
				disgo.Cache().CacheChannel(user)
			}
		}*/
	}

	genericInteractionEvent := events.GenericInteractionEvent{
		Event: api.Event{
			Disgo: disgo,
		},
		Interaction: *interaction,
	}
	eventManager.Dispatch(genericInteractionEvent)

	if interaction.Data != nil {
		options := interaction.Data.Options
		var subCommandName *string
		var subCommandGroup *string
		if len(options) == 1 {
			option := interaction.Data.Options[0]
			if option.Type == api.CommandOptionTypeSubCommandGroup {
				subCommandGroup = &option.Name
				options = option.Options
				option = option.Options[0]
			}
			if option.Type == api.CommandOptionTypeSubCommand {
				subCommandName = &option.Name
				options = option.Options
			}
		}
		var newOptions []*events.Option
		for _, optionData := range options {
			newOptions = append(newOptions, &events.Option{
				Resolved: interaction.Data.Resolved,
				Name:     optionData.Name,
				Type:     optionData.Type,
				Value:    optionData.Value,
			})
		}
		eventManager.Dispatch(events.SlashCommandEvent{
			GenericInteractionEvent: genericInteractionEvent,
			CommandID:               interaction.Data.ID,
			Name:                    interaction.Data.Name,
			SubCommandName:          subCommandName,
			SubCommandGroup:         subCommandGroup,
			Options:                 newOptions,
			Replied:                 false,
		})
	}

}
