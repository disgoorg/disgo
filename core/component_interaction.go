package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type ComponentInteraction struct {
	*Interaction
	ComponentInteractionData
	Message *Message
}

// DeferUpdate replies to the ComponentInteraction with discord.InteractionResponseTypeDeferredUpdateMessage and cancels the loading state
func (i *ComponentInteraction) DeferUpdate(opts ...rest.RequestOpt) error {
	return i.Respond(discord.InteractionResponseTypeDeferredUpdateMessage, nil, opts...)
}

// Update replies to the ComponentInteraction with discord.InteractionResponseTypeUpdateMessage & discord.MessageUpdate which edits the original Message
func (i *ComponentInteraction) Update(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) error {
	return i.Respond(discord.InteractionResponseTypeUpdateMessage, messageUpdate, opts...)
}

// Component returns the Component which issued this ComponentInteraction
func (i *ComponentInteraction) Component() Component {
	// this should never be nil
	return i.Message.ComponentByID(i.CustomID)
}

type ComponentInteractionData struct {
	CustomID      string
	ComponentType discord.ComponentType
}
