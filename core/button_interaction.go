package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

// ButtonInteractionFilter used to filter ButtonInteraction(s) in a collectors.ButtonClickCollector
type ButtonInteractionFilter func(buttonInteraction *ButtonInteraction) bool

var _ Interaction = (*ButtonInteraction)(nil)
var _ ComponentInteraction = (*ButtonInteraction)(nil)

type ButtonInteraction struct {
	*InteractionFields
	Message  *Message
	CustomID string
}

func (i *ButtonInteraction) InteractionType() discord.InteractionType {
	return discord.InteractionTypeComponent
}

func (i *ButtonInteraction) ComponentType() discord.ComponentType {
	return discord.ComponentTypeButton
}

func (i *ButtonInteraction) Respond(callbackType discord.InteractionCallbackType, callbackData discord.InteractionCallbackData, opts ...rest.RequestOpt) error {
	return respond(i.InteractionFields, callbackType, callbackData, opts...)
}

func (i *ButtonInteraction) Create(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) error {
	return create(i.InteractionFields, messageCreate, opts...)
}

func (i *ButtonInteraction) DeferCreate(ephemeral bool, opts ...rest.RequestOpt) error {
	return deferCreate(i.InteractionFields, ephemeral, opts...)
}

func (i *ButtonInteraction) Update(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) error {
	return update(i.InteractionFields, messageUpdate, opts...)
}

func (i *ButtonInteraction) DeferUpdate(opts ...rest.RequestOpt) error {
	return deferUpdate(i.InteractionFields, opts...)
}

// UpdateButton updates the clicked Button with a new Button
func (i *ButtonInteraction) UpdateButton(button discord.Button, opts ...rest.RequestOpt) error {
	return updateComponent(i.InteractionFields, i.Message, i.CustomID, button, opts...)
}

// Button returns the Button which issued this ButtonInteraction
func (i *ButtonInteraction) Button() discord.Button {
	// this should never be nil
	return *i.Message.ButtonByID(i.CustomID)
}
