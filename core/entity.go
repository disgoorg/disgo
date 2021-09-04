package core

import "github.com/DisgoOrg/disgo/discord"

type Entity interface {
	ID() discord.Snowflake
	Copy() Entity
	Update(entity Entity)
	Disgo() Disgo
}
