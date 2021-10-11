package core

import "github.com/DisgoOrg/disgo/discord"

type ContextCommandInteraction struct {
	*ApplicationCommandInteraction
	CreateInteractionResponses
	TargetID discord.Snowflake
}

type MessageCommandInteractionFilter func(messageCommandInteraction *MessageCommandInteraction) bool

type MessageCommandInteraction struct {
	*ContextCommandInteraction
}

func (i *MessageCommandInteraction) TargetMessage() *Message {
	return i.Resolved.Messages[i.TargetID]
}

type UserCommandInteractionFilter func(userCommandInteraction *UserCommandInteraction) bool

type UserCommandInteraction struct {
	*ContextCommandInteraction
}

func (i *UserCommandInteraction) TargetUser() *User {
	return i.Resolved.Users[i.TargetID]
}

func (i *UserCommandInteraction) TargetMember() *Member {
	return i.Resolved.Members[i.TargetID]
}
