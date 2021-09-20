package core

import "github.com/DisgoOrg/disgo/rest"

type ButtonInteraction struct {
	*ComponentInteraction
}

// Button returns the Button which issued this ButtonInteraction
func (i *ButtonInteraction) Button() Button {
	// this should never be nil
	return *i.Message.ButtonByID(i.CustomID)
}

// UpdateButton updates the clicked Button with a new Button
func (i *ButtonInteraction) UpdateButton(button Button, opts ...rest.RequestOpt) rest.Error {
	return i.UpdateComponent(button, opts...)
}
