package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type ModalSubmitInteractionFilter func(ModalSubmitInteraction *ModalSubmitInteraction) bool

var _ Interaction = (*ModalSubmitInteraction)(nil)

type ModalSubmitInteraction struct {
	CreateInteraction
	Data discord.ModalSubmitInteractionData
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
