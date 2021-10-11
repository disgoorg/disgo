package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type ComponentInteraction struct {
	*Interaction
	CreateInteractionResponses
	Message       *Message
	CustomID      string
	ComponentType discord.ComponentType
}

// DeferUpdate replies to the ComponentInteraction with discord.InteractionCallbackTypeDeferredUpdateMessage and cancels the loading state
func (i *ComponentInteraction) DeferUpdate(opts ...rest.RequestOpt) error {
	return i.Respond(discord.InteractionCallbackTypeDeferredUpdateMessage, nil, opts...)
}

// Update replies to the ComponentInteraction with discord.InteractionCallbackTypeUpdateMessage & discord.MessageUpdate which edits the original Message
func (i *ComponentInteraction) Update(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) error {
	return i.Respond(discord.InteractionCallbackTypeUpdateMessage, messageUpdate, opts...)
}

func (i *ComponentInteraction) UpdateComponent(component Component, opts ...rest.RequestOpt) error {
	actionRows := i.Message.ActionRows()
	for _, actionRow := range actionRows {
		actionRow = actionRow.SetComponent(i.CustomID, component)
	}

	messageUpdate := NewMessageUpdateBuilder().SetActionRows(actionRows...).Build()
	if i.Acknowledged {
		_, err := i.UpdateOriginal(messageUpdate, opts...)
		return err
	}
	return i.Update(messageUpdate, opts...)
}

// Component returns the Component which issued this ComponentInteraction
func (i *ComponentInteraction) Component() Component {
	// this should never be nil
	return i.Message.ComponentByID(i.CustomID)
}

// ButtonInteractionFilter used to filter ButtonInteraction(s) in a collectors.ButtonClickCollector
type ButtonInteractionFilter func(buttonInteraction *ButtonInteraction) bool

type ButtonInteraction struct {
	*ComponentInteraction
}

// Button returns the Button which issued this ButtonInteraction
func (i *ButtonInteraction) Button() Button {
	// this should never be nil
	return *i.Message.ButtonByID(i.CustomID)
}

// UpdateButton updates the clicked Button with a new Button
func (i *ButtonInteraction) UpdateButton(button Button, opts ...rest.RequestOpt) error {
	return i.UpdateComponent(button, opts...)
}

// SelectMenuInteractionFilter used to filter SelectMenuInteraction(s) in a collectors.SelectMenuSubmitCollector
type SelectMenuInteractionFilter func(selectMenuInteraction *SelectMenuInteraction) bool

type SelectMenuInteraction struct {
	*ComponentInteraction
	Values []string
}

// SelectMenu returns the SelectMenu which issued this SelectMenuInteraction
func (i *SelectMenuInteraction) SelectMenu() SelectMenu {
	// this should never be nil
	return *i.Message.SelectMenuByID(i.CustomID)
}

// UpdateSelectMenu updates the used SelectMenu with a new SelectMenu
func (i *SelectMenuInteraction) UpdateSelectMenu(selectMenu SelectMenu, opts ...rest.RequestOpt) error {
	return i.UpdateComponent(selectMenu, opts...)
}

// SelectedOptions returns the selected SelectMenuOption(s)
func (i *SelectMenuInteraction) SelectedOptions() []SelectMenuOption {
	options := make([]SelectMenuOption, len(i.Values))
	for ii, option := range i.SelectMenu().Options {
		for _, value := range i.Values {
			if value == option.Value {
				options[ii] = option
				break
			}
		}
	}
	return options
}
