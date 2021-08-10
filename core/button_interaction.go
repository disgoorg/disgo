package core

type ButtonInteraction struct {
	*ComponentInteraction
	Data *ButtonInteractionData
}

// Button returns the Button which issued this ButtonInteraction. nil for ephemeral Message(s)
func (i *ButtonInteraction) Button() *Button {
	return i.Message.ButtonByID(i.CustomID())
}

type ButtonInteractionData struct {
	*ComponentInteractionData
}
