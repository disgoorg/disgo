package disgo

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/internal"
)

// New Initialises a new Disgo client
//goland:noinspection GoUnusedExportedFunction
func New(token string, options api.Options) (api.Disgo, error) {
	return internal.New(token, options)
}

// NewBuilder creates an api.DisgoBuilder for the client
//goland:noinspection GoUnusedExportedFunction
func NewBuilder(token string) api.DisgoBuilder {
	return internal.NewBuilder(token)
}
