package core

import "github.com/DisgoOrg/disgo/discord"

// DefaultAllowedMentions returns the default discord.AllowedMentions for a Message
var DefaultAllowedMentions = discord.AllowedMentions{
	Parse:       []discord.AllowedMentionType{discord.AllowedMentionTypeUsers, discord.AllowedMentionTypeRoles, discord.AllowedMentionTypeEveryone},
	Roles:       []discord.Snowflake{},
	Users:       []discord.Snowflake{},
	RepliedUser: true,
}
