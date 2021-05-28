package events

import (
	"github.com/DisgoOrg/disgo/api"
)

// GenericInteractionEvent generic api.Interaction event
type GenericInteractionEvent struct {
	GenericEvent
	Interaction *api.Interaction
}

// Reply replies to the api.Interaction with the provided api.InteractionResponse
func (e *GenericInteractionEvent) Reply(response *api.InteractionResponse) error {
	return e.Interaction.Reply(response)
}

// EditOriginal edits the original api.InteractionResponse
func (e *GenericInteractionEvent) EditOriginal(followupMessage *api.FollowupMessage) (*api.Message, error) {
	return e.Interaction.EditOriginal(followupMessage)
}

// DeleteOriginal deletes the original api.InteractionResponse
func (e *GenericInteractionEvent) DeleteOriginal() error {
	return e.Interaction.DeleteOriginal()
}

// SendFollowup used to send a api.FollowupMessage to an api.Interaction
func (e *GenericInteractionEvent) SendFollowup(followupMessage *api.FollowupMessage) (*api.Message, error) {
	return e.Interaction.SendFollowup(followupMessage)
}

// EditFollowup used to edit a api.FollowupMessage from an api.Interaction
func (e *GenericInteractionEvent) EditFollowup(messageID api.Snowflake, followupMessage *api.FollowupMessage) (*api.Message, error) {
	return e.Interaction.EditFollowup(messageID, followupMessage)
}

// DeleteFollowup used to delete a api.FollowupMessage from an api.Interaction
func (e *GenericInteractionEvent) DeleteFollowup(messageID api.Snowflake) error {
	return e.Interaction.DeleteFollowup(messageID)
}

// CommandEvent indicates that a slash api.Command was ran
type CommandEvent struct {
	GenericInteractionEvent
	CommandInteraction  *api.CommandInteraction
	CommandID           api.Snowflake
	CommandName         string
	SubCommandName      *string
	SubCommandGroupName *string
	Options             []*api.Option
}

// DeferReply replies to the api.CommandInteraction with api.InteractionResponseTypeDeferredChannelMessageWithSource and shows a loading state
func (e *CommandEvent) DeferReply(ephemeral bool) error {
	return e.CommandInteraction.DeferReply(ephemeral)
}

// ReplyCreate replies to the api.CommandInteraction with api.InteractionResponseTypeDeferredChannelMessageWithSource & api.InteractionResponseData
func (e *CommandEvent) ReplyCreate(data *api.InteractionResponseData) error {
	return e.CommandInteraction.ReplyCreate(data)
}

// CommandPath returns the api.Command path
func (e CommandEvent) CommandPath() string {
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
func (e CommandEvent) Option(name string) *api.Option {
	options := e.OptionN(name)
	if len(options) == 0 {
		return nil
	}
	return options[0]
}

// OptionN returns Option(s) by name
func (e CommandEvent) OptionN(name string) []*api.Option {
	options := make([]*api.Option, 0)
	for _, option := range e.Options {
		if option.Name == name {
			options = append(options, option)
		}
	}
	return options
}

// OptionsT returns Option(s) by api.CommandOptionType
func (e CommandEvent) OptionsT(optionType api.CommandOptionType) []*api.Option {
	options := make([]*api.Option, 0)
	for _, option := range e.Options {
		if option.Type == optionType {
			options = append(options, option)
		}
	}
	return options
}

// ButtonClickEvent indicates that a api.Button was clicked
type ButtonClickEvent struct {
	GenericInteractionEvent
	ButtonInteraction *api.ButtonInteraction
}

// DeferEdit replies to the api.ButtonInteraction with api.InteractionResponseTypeDeferredUpdateMessage and cancels the loading state
func (e *ButtonClickEvent) DeferEdit() error {
	return e.ButtonInteraction.DeferEdit()
}

// ReplyEdit replies to the api.ButtonInteraction with api.InteractionResponseTypeUpdateMessage & api.InteractionResponseData which edits the original api.Message
func (e *ButtonClickEvent) ReplyEdit(data *api.InteractionResponseData) error {
	return e.ButtonInteraction.ReplyEdit(data)
}

// CustomID returns the customID from the called api.Button
func (e *ButtonClickEvent) CustomID() string {
	return e.ButtonInteraction.Data.CustomID
}

// ComponentType returns the api.ComponentType from the called api.Button
func (e *ButtonClickEvent) ComponentType() string {
	return e.ButtonInteraction.Data.CustomID
}

// Message returns the api.Message the api.Button is called from
func (e *ButtonClickEvent) Message() *api.Message {
	return e.ButtonInteraction.Message
}
