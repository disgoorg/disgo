package core

import "github.com/DisgoOrg/disgo/discord"

type ComponentInteraction struct {
	*Interaction
	Message *Message
	Data    *ComponentInteractionData
}

// DeferEdit replies to the ComponentInteraction with InteractionResponseTypeDeferredUpdateMessage and cancels the loading state
func (i *ComponentInteraction) DeferEdit() error {
	return i.Respond(discord.InteractionResponseTypeDeferredUpdateMessage, nil)
}

// Edit replies to the ComponentInteraction with InteractionResponseTypeUpdateMessage & MessageUpdate which edits the original Message
func (i *ComponentInteraction) Edit(messageUpdate discord.MessageUpdate) error {
	return i.Respond(discord.InteractionResponseTypeUpdateMessage, messageUpdate)
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
	if i.Message.IsEphemeral() {
		return nil
	}
	return i.Message.ComponentByID(i.CustomID())
}

type ComponentInteractionData struct {
	*InteractionData
}
