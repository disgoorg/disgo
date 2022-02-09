package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/snowflake"
)

type ApplicationCommandInteractionFilter func(interaction *ApplicationCommandInteraction) bool

var _ Interaction = (*ApplicationCommandInteraction)(nil)

// ApplicationCommandInteraction represents a generic ApplicationCommandInteraction received from discord
type ApplicationCommandInteraction struct {
	CreateInteraction
	Data ApplicationCommandInteractionData
}

func (i ApplicationCommandInteraction) interaction() {}
func (i ApplicationCommandInteraction) Type() discord.InteractionType {
	return discord.InteractionTypeApplicationCommand
}

func (i ApplicationCommandInteraction) SlashCommandInteractionData() SlashCommandInteractionData {
	return i.Data.(SlashCommandInteractionData)
}

func (i ApplicationCommandInteraction) UserCommandInteractionData() UserCommandInteractionData {
	return i.Data.(UserCommandInteractionData)
}

func (i ApplicationCommandInteraction) MessageCommandInteractionData() MessageCommandInteractionData {
	return i.Data.(MessageCommandInteractionData)
}

func (i ApplicationCommandInteraction) CreateModal(modalCreate discord.ModalCreate, opts ...rest.RequestOpt) error {
	return i.Respond(discord.InteractionCallbackTypeModal, modalCreate, opts...)
}

type ApplicationCommandInteractionData interface {
	discord.ApplicationCommandInteractionData
}

type SlashCommandInteractionData struct {
	discord.SlashCommandInteractionData
	SubCommandName      *string
	SubCommandGroupName *string
	Resolved            *SlashCommandResolved
	Options             SlashCommandOptionsMap
}

// CommandPath returns the ApplicationCommand path
func (i SlashCommandInteractionData) CommandPath() string {
	path := i.CommandName
	if name := i.SubCommandName; name != nil {
		path += "/" + *name
	}
	if name := i.SubCommandGroupName; name != nil {
		path += "/" + *name
	}
	return path
}

// SlashCommandResolved contains resolved mention data for SlashCommand(s)
type SlashCommandResolved struct {
	Users    map[snowflake.Snowflake]*User
	Members  map[snowflake.Snowflake]*Member
	Roles    map[snowflake.Snowflake]*Role
	Channels map[snowflake.Snowflake]Channel
}

type UserCommandInteractionData struct {
	discord.UserCommandInteractionData
	Resolved *UserCommandResolved
}

func (i *UserCommandInteractionData) TargetUser() *User {
	return i.Resolved.Users[i.TargetID]
}

func (i *UserCommandInteractionData) TargetMember() *Member {
	return i.Resolved.Members[i.TargetID]
}

type UserCommandResolved struct {
	Users   map[snowflake.Snowflake]*User
	Members map[snowflake.Snowflake]*Member
}

type MessageCommandInteractionData struct {
	discord.MessageCommandInteractionData
	Resolved *MessageCommandResolved
}

func (i *MessageCommandInteractionData) TargetMessage() *Message {
	return i.Resolved.Messages[i.TargetID]
}

type MessageCommandResolved struct {
	Messages map[snowflake.Snowflake]*Message
}
