package events

import (
	"errors"

	"github.com/DisgoOrg/disgo/api"
)

// GenericInteractionEvent generic api.Interaction event
type GenericInteractionEvent struct {
	GenericEvent
	Interaction     *api.Interaction
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

// EditOriginal edits the original api.InteractionResponse
func (e *GenericInteractionEvent) EditOriginal(followupMessage *api.FollowupMessage) (*api.Message, error) {
	return e.Disgo().RestClient().EditInteractionResponse(e.Disgo().ApplicationID(), e.Interaction.Token, followupMessage)
}

// DeleteOriginal deletes the original api.InteractionResponse
func (e *GenericInteractionEvent) DeleteOriginal() error {
	return e.Disgo().RestClient().DeleteInteractionResponse(e.Disgo().ApplicationID(), e.Interaction.Token)
}

// SendFollowup used to send a api.FollowupMessage to an api.Interaction
func (e *GenericInteractionEvent) SendFollowup(followupMessage *api.FollowupMessage) (*api.Message, error) {
	return e.Disgo().RestClient().SendFollowupMessage(e.Disgo().ApplicationID(), e.Interaction.Token, followupMessage)
}

// EditFollowup used to edit a api.FollowupMessage from an api.Interaction
func (e *GenericInteractionEvent) EditFollowup(messageID api.Snowflake, followupMessage *api.FollowupMessage) (*api.Message, error) {
	return e.Disgo().RestClient().EditFollowupMessage(e.Disgo().ApplicationID(), e.Interaction.Token, messageID, followupMessage)
}

// DeleteFollowup used to delete a api.FollowupMessage from an api.Interaction
func (e *GenericInteractionEvent) DeleteFollowup(messageID api.Snowflake) error {
	return e.Disgo().RestClient().DeleteFollowupMessage(e.Disgo().ApplicationID(), e.Interaction.Token, messageID)
}

// SlashCommandEvent indicates that a slash api.Command was ran
type SlashCommandEvent struct {
	GenericInteractionEvent
	SlashCommandInteraction *api.SlashCommandInteraction
	CommandID               api.Snowflake
	CommandName             string
	SubCommandName          *string
	SubCommandGroupName     *string
	Options                 []*api.Option
}

// DeferReply replies to the api.SlashCommandInteraction with api.InteractionResponseTypeDeferredChannelMessageWithSource and shows a loading state
func (e *SlashCommandEvent) DeferReply(ephemeral bool) error {
	var data *api.InteractionResponseData
	if ephemeral {
		data = &api.InteractionResponseData{Flags: api.MessageFlagEphemeral}
	}
	return e.Reply(&api.InteractionResponse{Type: api.InteractionResponseTypeDeferredChannelMessageWithSource, Data: data})
}

// ReplyCreate replies to the api.SlashCommandInteraction with api.InteractionResponseTypeDeferredChannelMessageWithSource & api.InteractionResponseData
func (e *SlashCommandEvent) ReplyCreate(data *api.InteractionResponseData) error {
	return e.Reply(&api.InteractionResponse{Type: api.InteractionResponseTypeChannelMessageWithSource, Data: data})
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

// ButtonClickEvent indicates that a api.Button was clicked
type ButtonClickEvent struct {
	GenericInteractionEvent
	ButtonInteraction *api.ButtonInteraction
}

// DeferEdit replies to the api.ButtonInteraction with api.InteractionResponseTypeDeferredUpdateMessage and cancels the loading state
func (e *ButtonClickEvent) DeferEdit() error {
	return e.Reply(&api.InteractionResponse{Type: api.InteractionResponseTypeDeferredUpdateMessage})
}

// ReplyEdit replies to the api.ButtonInteraction with api.InteractionResponseTypeUpdateMessage & api.InteractionResponseData which edits the original api.Message
func (e *ButtonClickEvent) ReplyEdit(data *api.InteractionResponseData) error {
	return e.Reply(&api.InteractionResponse{Type: api.InteractionResponseTypeUpdateMessage, Data: data})
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
