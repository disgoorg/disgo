package events

import (
	"errors"

	"github.com/DisgoOrg/disgo/api"
)

// GenericInteractionEvent generic api.Interaction event
type GenericInteractionEvent struct {
	GenericEvent
	Interaction *api.Interaction
}

// SlashCommandEvent indicates that a slash api.Command was ran in a api.Guild
type SlashCommandEvent struct {
	GenericInteractionEvent
	ResponseChannel     chan *api.InteractionResponse
	FromWebhook         bool
	CommandID           api.Snowflake
	CommandName         string
	SubCommandName      *string
	SubCommandGroupName *string
	Options             []*api.Option
	Replied             bool
}

// CommandPath returns the api.Command path
func (e SlashCommandEvent) CommandPath() string {
	path := e.CommandName
	if e.SubCommandName != nil {
		path += "/" + *e.SubCommandName
	}
	if e.SubCommandGroupName != nil {
		path += "/" + *e.SubCommandGroupName
	}
	return path
}

// Option returns an Option by name
func (e SlashCommandEvent) Option(name string) *api.Option {
	options := e.OptionN(name)
	if len(options) == 0 {
		return nil
	}
	return options[0]
}

// OptionN returns Option(s) by name
func (e SlashCommandEvent) OptionN(name string) []*api.Option {
	options := make([]*api.Option, 0)
	for _, option := range e.Options {
		if option.Name == name {
			options = append(options, option)
		}
	}
	return options
}

// OptionsT returns Option(s) by api.CommandOptionType
func (e SlashCommandEvent) OptionsT(optionType api.CommandOptionType) []*api.Option {
	options := make([]*api.Option, 0)
	for _, option := range e.Options {
		if option.Type == optionType {
			options = append(options, option)
		}
	}
	return options
}

// Acknowledge replies to the api.Interaction with api.InteractionResponseTypeDeferredChannelMessageWithSource
func (e *SlashCommandEvent) Acknowledge(ephemeral bool) error {
	var data *api.InteractionResponseData
	if ephemeral {
		data = &api.InteractionResponseData{
			Flags: api.MessageFlagEphemeral,
		}
	}
	return e.Reply(&api.InteractionResponse{
		Type: api.InteractionResponseTypeDeferredChannelMessageWithSource,
		Data: data,
	})
}

// Reply replies to the api.Interaction with the provided api.InteractionResponse
func (e *SlashCommandEvent) Reply(response *api.InteractionResponse) error {
	if e.Replied {
		return errors.New("you already replied to this interaction")
	}
	e.Replied = true

	if e.FromWebhook {
		e.ResponseChannel <- response
		return nil
	}

	return e.Interaction.Disgo.RestClient().SendInteractionResponse(e.Interaction.ID, e.Interaction.Token, response)
}
