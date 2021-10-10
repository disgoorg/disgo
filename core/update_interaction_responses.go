package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type UpdateInteractionResponses struct {
	*Interaction
}

// DeferUpdate replies to the ComponentInteraction with discord.InteractionCallbackTypeDeferredUpdateMessage and cancels the loading state
func (i *UpdateInteractionResponses) DeferUpdate(opts ...rest.RequestOpt) error {
	return i.Respond(discord.InteractionCallbackTypeDeferredUpdateMessage, nil, opts...)
}

// Update replies to the ComponentInteraction with discord.InteractionCallbackTypeUpdateMessage & discord.MessageUpdate which edits the original Message
func (i *UpdateInteractionResponses) Update(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) error {
	return i.Respond(discord.InteractionCallbackTypeUpdateMessage, messageUpdate, opts...)
}