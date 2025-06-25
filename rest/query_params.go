package rest

import (
	"github.com/disgoorg/disgo/discord"
)

// QueryParams serves as a generic interface for implementations of rest endpoint query parameters.
type QueryParams interface {
	// ToQueryValues transforms fields from the QueryParams interface implementations into discord.QueryValues.
	ToQueryValues() discord.QueryValues
}
