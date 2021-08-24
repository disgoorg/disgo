package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type ComponentInteraction struct {
	*Interaction
	Message *Message
	Data    *ComponentInteractionData
}

// DeferUpdate replies to the ComponentInteraction with discord.InteractionResponseTypeDeferredUpdateMessage and cancels the loading state
func (i *ComponentInteraction) DeferUpdate(opts ...rest.RequestOpt) error {
	return i.Respond(discord.InteractionResponseTypeDeferredUpdateMessage, nil, opts...)
}

// Update replies to the ComponentInteraction with discord.InteractionResponseTypeUpdateMessage & MessageUpdate which edits the original Message
func (i *ComponentInteraction) Update(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) error {
	return i.Respond(discord.InteractionResponseTypeUpdateMessage, messageUpdate, opts...)
}

// CustomID returns the Custom ID of the ComponentInteraction
func (i *ComponentInteraction) CustomID() string {
	return i.Data.CustomID
}

// ComponentType returns the ComponentType of a Component
func (i *ComponentInteraction) ComponentType() discord.ComponentType {
	return i.Data.ComponentType
}

// Component returns the Component which issued this ComponentInteraction. nil for ephemeral Message(s)
func (i *ComponentInteraction) Component() Component {
	return i.Message.ComponentByID(i.CustomID())
}

type ComponentInteractionData struct {
	*InteractionData
}
