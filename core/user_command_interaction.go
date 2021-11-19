package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type UserCommandInteractionFilter func(userCommandInteraction *UserCommandInteraction) bool

type UserCommandInteraction struct {
	discord.UserCommandInteraction
	*InteractionFields
	Data UserCommandInteractionData
}

type UserCommandInteractionData struct {
	discord.UserCommandInteractionData
	Resolved    *UserCommandResolved
}

func (i *UserCommandInteraction) InteractionType() discord.InteractionType {
	return discord.InteractionTypeApplicationCommand
}

func (i *UserCommandInteraction) ApplicationCommandType() discord.ApplicationCommandType {
	return discord.ApplicationCommandTypeUser
}

func (i *UserCommandInteraction) Respond(callbackType discord.InteractionCallbackType, callbackData discord.InteractionCallbackData, opts ...rest.RequestOpt) error {
	return respond(i.InteractionFields, callbackType, callbackData, opts...)
}

func (i *UserCommandInteraction) Create(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) error {
	return create(i.InteractionFields, messageCreate, opts...)
}

func (i *UserCommandInteraction) DeferCreate(ephemeral bool, opts ...rest.RequestOpt) error {
	return deferCreate(i.InteractionFields, ephemeral, opts...)
}

func (i *UserCommandInteraction) GetOriginal(opts ...rest.RequestOpt) (*Message, error) {
	return getOriginal(i.InteractionFields, opts...)
}

func (i *UserCommandInteraction) UpdateOriginal(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*Message, error) {
	return updateOriginal(i.InteractionFields, messageUpdate, opts...)
}

func (i *UserCommandInteraction) DeleteOriginal(opts ...rest.RequestOpt) error {
	return deleteOriginal(i.InteractionFields, opts...)
}

func (i *UserCommandInteraction) GetFollowup(messageID discord.Snowflake, opts ...rest.RequestOpt) (*Message, error) {
	return getFollowup(i.InteractionFields, messageID, opts...)
}

func (i *UserCommandInteraction) CreateFollowup(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*Message, error) {
	return createFollowup(i.InteractionFields, messageCreate, opts...)
}

func (i *UserCommandInteraction) UpdateFollowup(messageID discord.Snowflake, messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*Message, error) {
	return updateFollowup(i.InteractionFields, messageID, messageUpdate, opts...)
}

func (i *UserCommandInteraction) DeleteFollowup(messageID discord.Snowflake, opts ...rest.RequestOpt) error {
	return deleteFollowup(i.InteractionFields, messageID, opts...)
}

func (i *UserCommandInteraction) TargetUser() *User {
	return i.Data.Resolved.Users[i.Data.TargetID]
}

func (i *UserCommandInteraction) TargetMember() *Member {
	return i.Data.Resolved.Members[i.Data.TargetID]
}

type UserCommandResolved struct {
	Users   map[discord.Snowflake]*User
	Members map[discord.Snowflake]*Member
}
