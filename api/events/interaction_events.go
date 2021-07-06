package events

import (
	"github.com/DisgoOrg/disgo/api"
)

// GenericInteractionEvent generic api.Interaction event
type GenericInteractionEvent struct {
	*GenericEvent
	Interaction *api.Interaction
}

// Respond replies to the api.Interaction with the provided api.InteractionResponse
func (e GenericInteractionEvent) Respond(responseType api.InteractionResponseType, data interface{}) error {
	return e.Interaction.Respond(responseType, data)
}

// DeferReply replies to the api.CommandInteraction with api.InteractionResponseTypeDeferredChannelMessageWithSource and shows a loading state
func (e GenericInteractionEvent) DeferReply(ephemeral bool) error {
	return e.Interaction.DeferReply(ephemeral)
}

// Reply replies to the api.Interaction with api.InteractionResponseTypeDeferredChannelMessageWithSource & api.MessageCreate
func (e GenericInteractionEvent) Reply(messageCreate api.MessageCreate) error {
	return e.Interaction.Reply(messageCreate)
}

// EditOriginal edits the original api.InteractionResponse
func (e GenericInteractionEvent) EditOriginal(messageUpdate api.MessageUpdate) (*api.Message, error) {
	return e.Interaction.EditOriginal(messageUpdate)
}

// DeleteOriginal deletes the original api.InteractionResponse
func (e GenericInteractionEvent) DeleteOriginal() error {
	return e.Interaction.DeleteOriginal()
}

// SendFollowup used to send a followup api.MessageCreate to an api.Interaction
func (e GenericInteractionEvent) SendFollowup(messageCreate api.MessageCreate) (*api.Message, error) {
	return e.Interaction.SendFollowup(messageCreate)
}

// EditFollowup used to edit a followup api.Message from an api.Interaction
func (e GenericInteractionEvent) EditFollowup(messageID api.Snowflake, messageUpdate api.MessageUpdate) (*api.Message, error) {
	return e.Interaction.EditFollowup(messageID, messageUpdate)
}

// DeleteFollowup used to delete a followup api.Message from an api.Interaction
func (e GenericInteractionEvent) DeleteFollowup(messageID api.Snowflake) error {
	return e.Interaction.DeleteFollowup(messageID)
}

// CommandEvent indicates that a slash api.Command was ran
type CommandEvent struct {
	*GenericInteractionEvent
	CommandInteraction *api.CommandInteraction
}

// CommandID returns the ID of the api.Command which got used
func (e *CommandEvent) CommandID() api.Snowflake {
	return e.CommandInteraction.CommandID()
}

// CommandName the name of the api.Command which got used
func (e *CommandEvent) CommandName() string {
	return e.CommandInteraction.CommandName()
}

// SubCommandName the subcommand name of the api.Command which got used. May be nil
func (e *CommandEvent) SubCommandName() *string {
	return e.CommandInteraction.SubCommandName()
}

// SubCommandGroupName the subcommand group name of the api.Command which got used. May be nil
func (e *CommandEvent) SubCommandGroupName() *string {
	return e.CommandInteraction.SubCommandGroupName()
}

// CommandPath returns the api.Command path
func (e *CommandEvent) CommandPath() string {
	return e.CommandInteraction.CommandPath()
}

// Options returns the parsed api.Option which the api.Command got used with
func (e *CommandEvent) Options() []*api.Option {
	return e.CommandInteraction.Options()
}

// Option returns an Option by name
func (e *CommandEvent) Option(name string) *api.Option {
	options := e.OptionN(name)
	if len(options) == 0 {
		return nil
	}
	return options[0]
}

// OptionN returns Option(s) by name
func (e *CommandEvent) OptionN(name string) []*api.Option {
	options := make([]*api.Option, 0)
	for _, option := range e.Options() {
		if option.Name == name {
			options = append(options, option)
		}
	}
	return options
}

// OptionsT returns Option(s) by api.CommandOptionType
func (e *CommandEvent) OptionsT(optionType api.CommandOptionType) []*api.Option {
	options := make([]*api.Option, 0)
	for _, option := range e.Options() {
		if option.Type == optionType {
			options = append(options, option)
		}
	}
	return options
}

// GenericComponentEvent generic api.ComponentInteraction event
type GenericComponentEvent struct {
	*GenericInteractionEvent
	ComponentInteraction *api.ComponentInteraction
}

// DeferEdit replies to the api.ButtonInteraction with api.InteractionResponseTypeDeferredUpdateMessage and cancels the loading state
func (e *GenericComponentEvent) DeferEdit() error {
	return e.ComponentInteraction.DeferEdit()
}

// Edit replies to the api.ButtonInteraction with api.InteractionResponseTypeUpdateMessage & api.MessageUpdate which edits the original api.Message
func (e *GenericComponentEvent) Edit(messageUpdate api.MessageUpdate) error {
	return e.ComponentInteraction.Edit(messageUpdate)
}

// CustomID returns the customID from the called api.Component
func (e *GenericComponentEvent) CustomID() string {
	return e.ComponentInteraction.CustomID()
}

// ComponentType returns the api.ComponentType from the called api.Component
func (e *GenericComponentEvent) ComponentType() api.ComponentType {
	return e.ComponentInteraction.ComponentType()
}

// Component returns the api.Component from the event
func (e *GenericComponentEvent) Component() api.Component {
	return e.ComponentInteraction.Component()
}

// Message returns the api.Message of a GenericComponentEvent
func (e *GenericComponentEvent) Message() *api.Message {
	return e.ComponentInteraction.Message
}

// ButtonClickEvent indicates that a api.Button was clicked
type ButtonClickEvent struct {
	*GenericComponentEvent
	ButtonInteraction *api.ButtonInteraction
}

// Button returns the api.Button that was clicked on a ButtonClickEvent
func (e *ButtonClickEvent) Button() *api.Button {
	return e.ButtonInteraction.Button()
}

// SelectMenuSubmitEvent indicates that a api.SelectMenu was submitted
type SelectMenuSubmitEvent struct {
	*GenericComponentEvent
	SelectMenuInteraction *api.SelectMenuInteraction
}

// SelectMenu returns the api.SelectMenu of a SelectMenuSubmitEvent
func (e *SelectMenuSubmitEvent) SelectMenu() *api.SelectMenu {
	return e.SelectMenuInteraction.SelectMenu()
}

// Values returns the submitted values from the api.SelectMenu
func (e *SelectMenuSubmitEvent) Values() []string {
	return e.SelectMenuInteraction.Values()
}

// SelectedOptions returns a slice of api.SelectOption(s) that were chosen in an api.SelectMenu
func (e *SelectMenuSubmitEvent) SelectedOptions() []api.SelectOption {
	return e.SelectMenuInteraction.SelectedOptions()
}
