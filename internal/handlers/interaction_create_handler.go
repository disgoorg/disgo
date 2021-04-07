package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

// InteractionCreateHandler handles api.InteractionCreateGatewayEvent
type InteractionCreateHandler struct{}

// Event returns the raw gateway event Event
func (h InteractionCreateHandler) Event() api.GatewayEventName {
	return api.GatewayEventInteractionCreate
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
	handleInteraction(disgo, eventManager, nil, interaction)
}

func handleInteraction(disgo api.Disgo, eventManager api.EventManager, c chan interface{}, interaction *api.Interaction) {
	if interaction.Member != nil {
		interaction.Member.Disgo = disgo
		if interaction.Member.User != nil {
			interaction.Member.User.Disgo = disgo
		}
		disgo.Cache().CacheMember(interaction.Member)
	}
	if interaction.User != nil {
		interaction.User.Disgo = disgo
		disgo.Cache().CacheUser(interaction.User)
	}

	if interaction.Data != nil && interaction.Data.Resolved != nil {
		resolved := interaction.Data.Resolved
		if resolved.Users != nil {
			for _, user := range resolved.Users {
				user.Disgo = disgo
				disgo.Cache().CacheUser(user)
			}
		}
		if resolved.Members != nil {
			for id, member := range resolved.Members {
				member.User = resolved.Users[id]
				member.Disgo = disgo
				disgo.Cache().CacheMember(member)
			}
		}
		if resolved.Roles != nil {
			for _, role := range resolved.Roles {
				role.Disgo = disgo
				disgo.Cache().CacheRole(role)
			}
		}
		// TODO how do we cache partial channels?
		/*if resolved.Channels != nil {
			for _, channel := range resolved.Channels {
				channel.Disgo = disgo
				disgo.Cache().CacheChannel(channel)
			}
		}*/
	}

	genericInteractionEvent := events.GenericInteractionEvent{
		GenericEvent: events.NewEvent(disgo),
		Interaction:  *interaction,
	}
	eventManager.Dispatch(genericInteractionEvent)

	if interaction.Data != nil {
		options := interaction.Data.Options
		var subCommandName *string
		var subCommandGroupName *string
		if len(options) == 1 {
			option := interaction.Data.Options[0]
			if option.Type == api.OptionTypeSubCommandGroup {
				subCommandGroupName = &option.Name
				options = option.Options
				option = option.Options[0]
			}
			if option.Type == api.OptionTypeSubCommand {
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
			ResponseChannel:         c,
			FromWebhook:             c != nil,
			GenericInteractionEvent: genericInteractionEvent,
			CommandID:               interaction.Data.ID,
			CommandName:             interaction.Data.Name,
			SubCommandName:          subCommandName,
			SubCommandGroupName:     subCommandGroupName,
			Options:                 newOptions,
			Replied:                 false,
		})
	}
}
