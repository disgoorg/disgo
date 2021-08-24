package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

// GenericComponentEvent generic api.ComponentInteraction event
type GenericComponentEvent struct {
	*GenericInteractionEvent
	ComponentInteraction *core.ComponentInteraction
}

// DeferUpdate replies to the api.ButtonInteraction with api.InteractionResponseTypeDeferredUpdateMessage and cancels the loading state
func (e *GenericComponentEvent) DeferUpdate(opts ...rest.RequestOpt) error {
	return e.ComponentInteraction.DeferUpdate(opts...)
}

// Update replies to the api.ButtonInteraction with api.InteractionResponseTypeUpdateMessage & api.MessageUpdate which edits the original api.Message
func (e *GenericComponentEvent) Update(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) error {
	return e.ComponentInteraction.Update(messageUpdate, opts...)
}

// CustomID returns the customID from the called api.Component
func (e *GenericComponentEvent) CustomID() string {
	return e.ComponentInteraction.CustomID()
}

// ComponentType returns the api.ComponentType from the called api.Component
func (e *GenericComponentEvent) ComponentType() discord.ComponentType {
	return e.ComponentInteraction.ComponentType()
}

// Component returns the api.Component from the event
func (e *GenericComponentEvent) Component() core.Component {
	return e.ComponentInteraction.Component()
}

// Message returns the api.Message of a GenericComponentEvent
func (e *GenericComponentEvent) Message() *core.Message {
	return e.ComponentInteraction.Message
}

// ButtonClickEvent indicates that an api.Button was clicked
type ButtonClickEvent struct {
	*GenericComponentEvent
	ButtonInteraction *core.ButtonInteraction
}

// Button returns the api.Button that was clicked on a ButtonClickEvent
func (e *ButtonClickEvent) Button() *core.Button {
	return e.ButtonInteraction.Button()
}

// SelectMenuSubmitEvent indicates that an api.SelectMenu was submitted
type SelectMenuSubmitEvent struct {
	*GenericComponentEvent
	SelectMenuInteraction *core.SelectMenuInteraction
}

// SelectMenu returns the api.SelectMenu of a SelectMenuSubmitEvent
func (e *SelectMenuSubmitEvent) SelectMenu() *core.SelectMenu {
	return e.SelectMenuInteraction.SelectMenu()
}

// Values returns the submitted values from the api.SelectMenu
func (e *SelectMenuSubmitEvent) Values() []string {
	return e.SelectMenuInteraction.Values()
}

// SelectedOptions returns a slice of api.SelectOption(s) that were chosen in an api.SelectMenu
func (e *SelectMenuSubmitEvent) SelectedOptions() []discord.SelectOption {
	return e.SelectMenuInteraction.SelectedOptions()
}
