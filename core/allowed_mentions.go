package core

import "github.com/DisgoOrg/disgo/discord"

// DefaultMessageAllowedMentions gives you the default AllowedMentions for a Message
var DefaultMessageAllowedMentions = discord.AllowedMentions{
	Parse:       []discord.AllowedMentionType{discord.AllowedMentionTypeUsers, discord.AllowedMentionTypeRoles, discord.AllowedMentionTypeEveryone},
	Roles:       []discord.Snowflake{},
	Users:       []discord.Snowflake{},
	RepliedUser: true,
}
