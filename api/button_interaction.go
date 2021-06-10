package api

// ButtonInteraction is a specific Interaction when CLicked on Button(s)
type ButtonInteraction struct {
	*Interaction
	Message *Message               `json:"message,omitempty"`
	Data    *ButtonInteractionData `json:"data,omitempty"`
}

// DeferEdit replies to the api.ButtonInteraction with api.InteractionResponseTypeDeferredUpdateMessage and cancels the loading state
func (i *ButtonInteraction) DeferEdit() error {
	return i.Respond(InteractionResponseTypeDeferredUpdateMessage, nil)
}

// Edit replies to the api.ButtonInteraction with api.InteractionResponseTypeUpdateMessage & api.MessageCreate which edits the original api.Message
func (i *ButtonInteraction) Edit(messageCreate MessageCreate) error {
	return i.Respond(InteractionResponseTypeUpdateMessage, messageCreate)
}

// ButtonInteractionData is the command data payload
type ButtonInteractionData struct {
	CustomID      string        `json:"custom_id"`
	ComponentType ComponentType `json:"component_type"`
}
