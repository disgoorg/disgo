package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

type ModalSubmitInteractionFilter func(ModalSubmitInteraction *ModalSubmitInteraction) bool

var _ Interaction = (*ModalSubmitInteraction)(nil)

type ModalSubmitInteraction struct {
	*ReplyInteraction
	Data discord.ModalSubmitInteractionData
}

func (i ModalSubmitInteraction) interaction() {}
func (i ModalSubmitInteraction) Type() discord.InteractionType {
	return discord.InteractionTypeModalSubmit
}
