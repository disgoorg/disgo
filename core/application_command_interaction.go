package core

import "github.com/DisgoOrg/disgo/discord"

type ApplicationCommandInteraction struct {
	*Interaction
	Data *ApplicationCommandInteractionData `json:"data,omitempty"`
}

// CommandType returns the type of ApplicationCommand which was used
func (i *ApplicationCommandInteraction) CommandType() discord.ApplicationCommandType {
	return i.Data.CommandType
}

// CommandID returns the ID of the ApplicationCommand which was used
func (i *ApplicationCommandInteraction) CommandID() discord.Snowflake {
	return i.Data.ID
}

// CommandName the name of the ApplicationCommand which was used
func (i *ApplicationCommandInteraction) CommandName() string {
	return i.Data.Name
}

// Resolved returns all Resolved mentions from this ApplicationCommand
func (i *ApplicationCommandInteraction) Resolved() *Resolved {
	return i.Data.Resolved
}

// ApplicationCommandInteractionData is the command data payload
type ApplicationCommandInteractionData struct {
	*InteractionData
	Resolved *Resolved
}

// Resolved contains resolved mention data
type Resolved struct {
	discord.Resolved
	Users    map[discord.Snowflake]*User
	Members  map[discord.Snowflake]*Member
	Roles    map[discord.Snowflake]*Role
	Channels map[discord.Snowflake]Channel
	Messages map[discord.Snowflake]*Message
}
