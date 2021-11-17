package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type MessageCommandInteractionFilter func(messageCommandInteraction *MessageCommandInteraction) bool

type MessageCommandInteraction struct {
	*InteractionFields
	CommandID   discord.Snowflake
	CommandName string
	Resolved    *MessageCommandResolved
	TargetID    discord.Snowflake
}

func (i *MessageCommandInteraction) InteractionType() discord.InteractionType {
	return discord.InteractionTypeApplicationCommand
}

func (i *MessageCommandInteraction) ApplicationCommandType() discord.ApplicationCommandType {
	return discord.ApplicationCommandTypeMessage
}

func (i *MessageCommandInteraction) Respond(callbackType discord.InteractionCallbackType, callbackData discord.InteractionCallbackData, opts ...rest.RequestOpt) error {
	return respond(i.InteractionFields, callbackType, callbackData, opts...)
}

func (i *MessageCommandInteraction) Create(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) error {
	return create(i.InteractionFields, messageCreate, opts...)
}

func (i *MessageCommandInteraction) DeferCreate(ephemeral bool, opts ...rest.RequestOpt) error {
	return deferCreate(i.InteractionFields, ephemeral, opts...)
}

func (i *MessageCommandInteraction) GetOriginal(opts ...rest.RequestOpt) (*Message, error) {
	return getOriginal(i.InteractionFields, opts...)
}

func (i *MessageCommandInteraction) UpdateOriginal(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*Message, error) {
	return updateOriginal(i.InteractionFields, messageUpdate, opts...)
}

func (i *MessageCommandInteraction) DeleteOriginal(opts ...rest.RequestOpt) error {
	return deleteOriginal(i.InteractionFields, opts...)
}

func (i *MessageCommandInteraction) GetFollowup(messageID discord.Snowflake, opts ...rest.RequestOpt) (*Message, error) {
	return getFollowup(i.InteractionFields, messageID, opts...)
}

func (i *MessageCommandInteraction) CreateFollowup(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*Message, error) {
	return createFollowup(i.InteractionFields, messageCreate, opts...)
}

func (i *MessageCommandInteraction) UpdateFollowup(messageID discord.Snowflake, messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*Message, error) {
	return updateFollowup(i.InteractionFields, messageID, messageUpdate, opts...)
}

func (i *MessageCommandInteraction) DeleteFollowup(messageID discord.Snowflake, opts ...rest.RequestOpt) error {
	return deleteFollowup(i.InteractionFields, messageID, opts...)
}

func (i *MessageCommandInteraction) TargetMessage() *Message {
	return i.Resolved.Messages[i.TargetID]
}

type MessageCommandResolved struct {
	Messages map[discord.Snowflake]*Message
}
