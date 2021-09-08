package core

type ButtonInteraction struct {
	*ComponentInteraction
}

// Button returns the Button which issued this ButtonInteraction
func (i *ButtonInteraction) Button() Button {
	// this should never be nil
	return *i.Message.ButtonByID(i.CustomID)
}
