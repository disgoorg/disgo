package disgo

import (
	"github.com/DiscoOrg/disgo/api"
	"github.com/DiscoOrg/disgo/internal"
)

// New Initialises a new Disgo client
func New(token string, options api.Options) (api.Disgo, error) {
	return internal.New(token, options)
}

// NewBuilder creates an api.DisgoBuilder for the client
func NewBuilder(token string) api.DisgoBuilder {
	return internal.NewBuilder(token)
}