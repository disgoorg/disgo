package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type ModalSubmitInteractionFilter func(ModalSubmitInteraction *ModalSubmitInteraction) bool

var _ Interaction = (*ModalSubmitInteraction)(nil)

type ModalSubmitInteraction struct {
	CreateInteraction
	Data ModalSubmitInteractionData
}

func (i ModalSubmitInteraction) interaction() {}
func (i ModalSubmitInteraction) Type() discord.InteractionType {
	return discord.InteractionTypeModalSubmit
}

func (i ModalSubmitInteraction) UpdateMessage(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) error {
	return i.Respond(discord.InteractionCallbackTypeUpdateMessage, messageUpdate, opts...)
}

func (i ModalSubmitInteraction) DeferUpdateMessage(opts ...rest.RequestOpt) error {
	return i.Respond(discord.InteractionCallbackTypeDeferredUpdateMessage, nil, opts...)
}

type ModalSubmitInteractionData struct {
	discord.ModalSubmitInteractionData
	Components ModalComponentsMap
}

type ModalComponentsMap map[discord.CustomID]discord.InputComponent

func (m ModalComponentsMap) Get(customID discord.CustomID) discord.InputComponent {
	if component, ok := m[customID]; ok {
		return component
	}
	return nil
}

func (m ModalComponentsMap) TextComponent(customID discord.CustomID) *discord.TextInputComponent {
	component := m.Get(customID)
	if component == nil {
		return nil
	}
	if cmp, ok := component.(discord.TextInputComponent); ok {
		return &cmp
	}
	return nil
}

func (m ModalComponentsMap) Text(customID discord.CustomID) *string {
	component := m.TextComponent(customID)
	if component == nil {
		return nil
	}
	return &component.Value
}
