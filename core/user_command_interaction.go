package core

import "github.com/DisgoOrg/disgo/discord"

type UserCommandInteractionFilter func(userCommandInteraction *UserCommandInteraction) bool

type UserCommandInteraction struct {
	discord.UserCommandInteraction
	InteractionData
	UserCommandInteractionData
}

func (i *UserCommandInteraction) TargetUser() *User {
	return i.Resolved.Users[i.TargetID]
}

func (i *UserCommandInteraction) TargetMember() *Member {
	return i.Resolved.Members[i.TargetID]
}

type UserCommandInteractionData struct {
	CommandID           discord.Snowflake
	CommandName         string
	SubCommandName      *string
	SubCommandGroupName *string
	Resolved            UserCommandResolved
	TargetID            discord.Snowflake
}

type UserCommandResolved struct {
	Users   map[discord.Snowflake]*User
	Members map[discord.Snowflake]*Member
}
