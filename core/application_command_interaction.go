package core

import "github.com/DisgoOrg/disgo/discord"

type ApplicationCommandInteraction struct {
	*Interaction
	ApplicationCommandInteractionData
}

// ApplicationCommandInteractionData is the command data payload
type ApplicationCommandInteractionData struct {
	CommandName string
	Resolved    *Resolved
}

// Resolved contains resolved mention data
type Resolved struct {
	Users    map[discord.Snowflake]*User
	Members  map[discord.Snowflake]*Member
	Roles    map[discord.Snowflake]*Role
	Channels map[discord.Snowflake]*Channel
	Messages map[discord.Snowflake]*Message
}
