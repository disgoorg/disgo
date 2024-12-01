package rest

import (
	"github.com/disgoorg/disgo/discord"
)

type QueryParams interface {
	ToQueryValues() discord.QueryValues
}
