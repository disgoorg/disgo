package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type ComponentInteractionFilter func(interaction *ComponentInteraction) bool

var _ Interaction = (*ComponentInteraction)(nil)

// ComponentInteraction represents a generic ComponentInteraction received from discord
type ComponentInteraction struct {
	CreateInteraction
	Data    ComponentInteractionData
	Message *Message
}

func (i ComponentInteraction) interaction() {}
func (i ComponentInteraction) Type() discord.InteractionType {
	return discord.InteractionTypeComponent
}

func (i ComponentInteraction) ButtonInteractionData() ButtonInteractionData {
	return i.Data.(ButtonInteractionData)
}

func (i ComponentInteraction) SelectMenuInteractionData() SelectMenuInteractionData {
	return i.Data.(SelectMenuInteractionData)
}

func (i ComponentInteraction) UpdateMessage(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) error {
	return i.Respond(discord.InteractionCallbackTypeUpdateMessage, messageUpdate, opts...)
}

func (i ComponentInteraction) DeferUpdateMessage(opts ...rest.RequestOpt) error {
	return i.Respond(discord.InteractionCallbackTypeDeferredUpdateMessage, nil, opts...)
}

func (i ComponentInteraction) UpdateComponent(customID discord.CustomID, component discord.InteractiveComponent, opts ...rest.RequestOpt) error {
	containerComponents := make([]discord.ContainerComponent, len(i.Message.Components))
	for ii := range i.Message.Components {
		switch container := containerComponents[ii].(type) {
		case discord.ActionRowComponent:
			containerComponents[ii] = container.UpdateComponent(customID, component)

		default:
			containerComponents[ii] = container
			continue
		}
	}

	return i.UpdateMessage(discord.NewMessageUpdateBuilder().SetContainerComponents(containerComponents...).Build(), opts...)
}

func (i ComponentInteraction) CreateModal(modalCreate discord.ModalCreate, opts ...rest.RequestOpt) error {
	return i.Respond(discord.InteractionCallbackTypeModal, modalCreate, opts...)
}

type ComponentInteractionData interface {
	discord.ComponentInteractionData
}

type ButtonInteractionData struct {
	discord.ButtonInteractionData
	interaction *ComponentInteraction
}

// UpdateButton updates the clicked ButtonComponent with a new ButtonComponent
func (d *ButtonInteractionData) UpdateButton(button discord.ButtonComponent, opts ...rest.RequestOpt) error {
	return d.interaction.UpdateComponent(d.CustomID, button, opts...)
}

// ButtonComponent returns the ButtonComponent which issued this ButtonInteraction
func (d *ButtonInteractionData) ButtonComponent() discord.ButtonComponent {
	// this should never be nil
	return *d.interaction.Message.ButtonByID(d.CustomID)
}

type SelectMenuInteractionData struct {
	discord.SelectMenuInteractionData
	interaction *ComponentInteraction
}

// SelectMenuComponent returns the SelectMenuComponent which issued this SelectMenuInteraction
func (d *SelectMenuInteractionData) SelectMenuComponent() discord.SelectMenuComponent {
	// this should never be nil
	return *d.interaction.Message.SelectMenuByID(d.CustomID)
}

// SelectedOptions returns the selected SelectMenuOption(s)
func (d *SelectMenuInteractionData) SelectedOptions() []discord.SelectMenuOption {
	options := make([]discord.SelectMenuOption, len(d.Values))
	for ii, option := range d.SelectMenuComponent().Options {
		for _, value := range d.Values {
			if value == option.Value {
				options[ii] = option
				break
			}
		}
	}
	return options
}
