package disgo

import (
	"github.com/DiscoOrg/disgo/api"
	"github.com/DiscoOrg/disgo/internal"
)

func New(token string, options api.Options) api.Disgo {
	return internal.New(token, options)
}

func NewBuilder(token string) api.DisgoBuilder {
	return internal.NewBuilder(token)
}