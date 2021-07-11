package api

// ButtonInteraction is a specific Interaction when CLicked on Button(s)
type ButtonInteraction struct {
	*ComponentInteraction
	Data *ButtonInteractionData `json:"data,omitempty"`
}

// Button returns the Button which issued this ButtonInteraction. nil for ephemeral Message(s)
func (i *ButtonInteraction) Button() *Button {
	return i.Message.ButtonByID(i.CustomID())
}

// ButtonInteractionData is the Button data payload
type ButtonInteractionData struct {
	*ComponentInteractionData
}
