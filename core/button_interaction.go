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
	discord.ButtonInteraction
	*InteractionFields
	Message  *Message
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

func (i *ButtonInteraction) GetOriginal(opts ...rest.RequestOpt) (*Message, error) {
	return getOriginal(i.InteractionFields, opts...)
}

func (i *ButtonInteraction) UpdateOriginal(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*Message, error) {
	return updateOriginal(i.InteractionFields, messageUpdate, opts...)
}

func (i *ButtonInteraction) DeleteOriginal(opts ...rest.RequestOpt) error {
	return deleteOriginal(i.InteractionFields, opts...)
}

func (i *ButtonInteraction) GetFollowup(messageID discord.Snowflake, opts ...rest.RequestOpt) (*Message, error) {
	return getFollowup(i.InteractionFields, messageID, opts...)
}

func (i *ButtonInteraction) CreateFollowup(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*Message, error) {
	return createFollowup(i.InteractionFields, messageCreate, opts...)
}

func (i *ButtonInteraction) UpdateFollowup(messageID discord.Snowflake, messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*Message, error) {
	return updateFollowup(i.InteractionFields, messageID, messageUpdate, opts...)
}

func (i *ButtonInteraction) DeleteFollowup(messageID discord.Snowflake, opts ...rest.RequestOpt) error {
	return deleteFollowup(i.InteractionFields, messageID, opts...)
}

func (i *ButtonInteraction) Update(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) error {
	return update(i.InteractionFields, messageUpdate, opts...)
}

func (i *ButtonInteraction) DeferUpdate(opts ...rest.RequestOpt) error {
	return deferUpdate(i.InteractionFields, opts...)
}

// UpdateButton updates the clicked ButtonComponent with a new ButtonComponent
func (i *ButtonInteraction) UpdateButton(button discord.ButtonComponent, opts ...rest.RequestOpt) error {
	return updateComponent(i.InteractionFields, i.Message, i.Data.CustomID, button, opts...)
}

// ButtonComponent returns the ButtonComponent which issued this ButtonInteraction
func (i *ButtonInteraction) ButtonComponent() discord.ButtonComponent {
	// this should never be nil
	return *i.Message.ButtonByID(i.Data.CustomID)
}
