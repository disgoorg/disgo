package api

// ComponentInteraction is a specific Interaction when using Component(s)
type ComponentInteraction struct {
	*Interaction
	Message *Message                  `json:"message,omitempty"`
	Data    *ComponentInteractionData `json:"data,omitempty"`
}

// DeferEdit replies to the api.ComponentInteraction with api.InteractionResponseTypeDeferredUpdateMessage and cancels the loading state
func (i *ComponentInteraction) DeferEdit() error {
	return i.Respond(InteractionResponseTypeDeferredUpdateMessage, nil)
}

// Edit replies to the api.ComponentInteraction with api.InteractionResponseTypeUpdateMessage & api.MessageUpdate which edits the original api.Message
func (i *ComponentInteraction) Edit(messageUpdate MessageUpdate) error {
	return i.Respond(InteractionResponseTypeUpdateMessage, messageUpdate)
}

// CustomID returns the Custom ID of the ComponentInteraction
func (i *ComponentInteraction) CustomID() string {
	return i.Data.CustomID
}

// ComponentType returns the ComponentType of a Component
func (i *ComponentInteraction) ComponentType() ComponentType {
	return i.Data.ComponentType
}

// Component returns the Component which issued this ComponentInteraction. nil for ephemeral Message(s)
func (i *ComponentInteraction) Component() Component {
	if i.Message.IsEphemeral() {
		return nil
	}
	return i.Message.ComponentByID(i.CustomID())
}

// ComponentInteractionData is the Component data payload
type ComponentInteractionData struct {
	CustomID      string        `json:"custom_id"`
	ComponentType ComponentType `json:"component_type"`
}
