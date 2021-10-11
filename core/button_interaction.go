package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

// ButtonInteractionFilter used to filter ButtonInteraction(s) in a collectors.ButtonClickCollector
type ButtonInteractionFilter func(buttonInteraction *ButtonInteraction) bool

type ButtonInteraction struct {
	*ComponentInteraction
}

// Button returns the Button which issued this ButtonInteraction
func (i *ButtonInteraction) Button() discord.Button {
	// this should never be nil
	return *i.Message.ButtonByID(i.CustomID)
}

// UpdateButton updates the clicked Button with a new Button
func (i *ButtonInteraction) UpdateButton(button discord.Button, opts ...rest.RequestOpt) error {
	return i.UpdateComponent(button, opts...)
}
