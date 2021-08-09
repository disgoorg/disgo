package core

import "github.com/DisgoOrg/disgo/discord"

type Role struct {
	discord.Role
	Disgo Disgo
}
