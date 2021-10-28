package core

import "github.com/DisgoOrg/disgo/discord"

type UserCommandInteractionFilter func(userCommandInteraction *UserCommandInteraction) bool

type UserCommandInteraction struct {
	*InteractionFields
	CommandID   discord.Snowflake
	CommandName string
	Resolved    *UserCommandResolved
	TargetID    discord.Snowflake
}

func (i *UserCommandInteraction) TargetUser() *User {
	return i.Resolved.Users[i.TargetID]
}

func (i *UserCommandInteraction) TargetMember() *Member {
	return i.Resolved.Members[i.TargetID]
}

type UserCommandResolved struct {
	Users   map[discord.Snowflake]*User
	Members map[discord.Snowflake]*Member
}
