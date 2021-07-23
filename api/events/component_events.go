package events

import "github.com/DisgoOrg/disgo/api"

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
