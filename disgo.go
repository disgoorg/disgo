package disgo

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/endpoints"
	"github.com/DisgoOrg/disgo/internal"
)

// New Initialises a new Disgo client
func New(token endpoints.Token, options api.Options) (api.Disgo, error) {
	return internal.New(token, options)
}

// NewBuilder creates an api.DisgoBuilder for the client
func NewBuilder(token endpoints.Token) api.DisgoBuilder {
	return internal.NewBuilder(token)
}
