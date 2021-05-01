package events

import (
	"errors"

	"github.com/DisgoOrg/disgo/api"
)

// GenericInteractionEvent generic api.Interaction event
type GenericInteractionEvent struct {
	GenericEvent
	*api.Interaction
	ResponseChannel chan *api.InteractionResponse
	FromWebhook     bool
	Replied         bool
}

// Reply replies to the api.Interaction with the provided api.InteractionResponse
func (e *GenericInteractionEvent) Reply(response *api.InteractionResponse) error {
	if e.Replied {
		return errors.New("you already replied to this interaction")
	}
	e.Replied = true

	if e.FromWebhook {
		e.ResponseChannel <- response
		return nil
	}

	return e.Disgo().RestClient().SendInteractionResponse(e.Interaction.ID, e.Interaction.Token, response)
}

// Guild returns the api.Guild from the api.Cache
func (e GenericInteractionEvent) Guild() *api.Guild {
	if e.Interaction.GuildID == nil {
		return nil
	}
	return e.Disgo().Cache().Guild(*e.Interaction.GuildID)
}

// DMChannel returns the api.DMChannel from the api.Cache
func (e GenericInteractionEvent) DMChannel() *api.DMChannel {
	if e.Interaction.ChannelID == nil {
		return nil
	}
	return e.Disgo().Cache().DMChannel(*e.Interaction.ChannelID)
}

// MessageChannel returns the api.MessageChannel from the api.Cache
func (e GenericInteractionEvent) MessageChannel() *api.MessageChannel {
	if e.Interaction.ChannelID == nil {
		return nil
	}
	return e.Disgo().Cache().MessageChannel(*e.Interaction.ChannelID)
}

// TextChannel returns the api.TextChannel from the api.Cache
func (e GenericInteractionEvent) TextChannel() *api.TextChannel {
	if e.Interaction.ChannelID == nil {
		return nil
	}
	return e.Disgo().Cache().TextChannel(*e.Interaction.ChannelID)
}

// GuildChannel returns the api.GuildChannel from the api.Cache
func (e GenericInteractionEvent) GuildChannel() *api.GuildChannel {
	if e.Interaction.ChannelID == nil {
		return nil
	}
	return e.Disgo().Cache().GuildChannel(*e.Interaction.ChannelID)
}

// SlashCommandEvent indicates that a slash api.Command was ran in a api.Guild
type SlashCommandEvent struct {
	GenericInteractionEvent
	CommandID           api.Snowflake
	CommandName         string
	SubCommandName      *string
	SubCommandGroupName *string
	Options             []*api.Option
	Resolved            *api.Resolved
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
func (e *SlashCommandEvent) Acknowledge() error {
	return e.Reply(&api.InteractionResponse{Type: api.InteractionResponseTypeDeferredChannelMessageWithSource})
}

// EditOriginal edits the original api.InteractionResponse
func (e *SlashCommandEvent) EditOriginal(followupMessage *api.FollowupMessage) (*api.Message, error) {
	return e.Disgo().RestClient().EditInteractionResponse(e.Disgo().ApplicationID(), e.Interaction.Token, followupMessage)
}

// DeleteOriginal deletes the original api.InteractionResponse
func (e *SlashCommandEvent) DeleteOriginal() error {
	return e.Disgo().RestClient().DeleteInteractionResponse(e.Disgo().ApplicationID(), e.Interaction.Token)
}

// SendFollowup used to send a api.FollowupMessage to an api.Interaction
func (e *SlashCommandEvent) SendFollowup(followupMessage *api.FollowupMessage) (*api.Message, error) {
	return e.Disgo().RestClient().SendFollowupMessage(e.Disgo().ApplicationID(), e.Interaction.Token, followupMessage)
}

// EditFollowup used to edit a api.FollowupMessage from an api.Interaction
func (e *SlashCommandEvent) EditFollowup(messageID api.Snowflake, followupMessage *api.FollowupMessage) (*api.Message, error) {
	return e.Disgo().RestClient().EditFollowupMessage(e.Disgo().ApplicationID(), e.Interaction.Token, messageID, followupMessage)
}

// DeleteFollowup used to delete a api.FollowupMessage from an api.Interaction
func (e *SlashCommandEvent) DeleteFollowup(messageID api.Snowflake) error {
	return e.Disgo().RestClient().DeleteFollowupMessage(e.Disgo().ApplicationID(), e.Interaction.Token, messageID)
}

type ButtonClickEvent struct {
	GenericInteractionEvent
	CustomID      string
	ComponentType api.ComponentType
	Message       *api.Message
}
