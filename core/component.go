package core

import "github.com/DisgoOrg/disgo/discord"

// Component is a general interface each Component needs to implement
type Component interface {
	Type() discord.ComponentType
}
