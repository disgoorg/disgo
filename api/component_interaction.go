package api

// ComponentInteraction is a specific Interaction when using Component(s)
type ComponentInteraction struct {
	*Interaction
	Message *Message                  `json:"message,omitempty"`
	Data    *ComponentInteractionData `json:"data,omitempty"`
}

// ComponentInteractionData is the Component data payload
type ComponentInteractionData struct {
	CustomID      string        `json:"custom_id"`
	ComponentType ComponentType `json:"component_type"`
}

// DeferEdit replies to the api.ComponentInteraction with api.InteractionResponseTypeDeferredUpdateMessage and cancels the loading state
func (i *ComponentInteraction) DeferEdit() error {
	return i.Respond(InteractionResponseTypeDeferredUpdateMessage, nil)
}

// Edit replies to the api.ComponentInteraction with api.InteractionResponseTypeUpdateMessage & api.MessageUpdate which edits the original api.Message
func (i *ComponentInteraction) Edit(messageUpdate MessageUpdate) error {
	return i.Respond(InteractionResponseTypeUpdateMessage, messageUpdate)
}

// Component returns the Component which issued this ComponentInteraction. nil for ephemeral Message(s)
func (i *ComponentInteraction) Component() Component {
	return i.Message.ComponentByID(i.Data.CustomID)
}
