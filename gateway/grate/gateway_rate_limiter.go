package grate

import (
	"context"

	"github.com/DisgoOrg/log"
)

type Limiter interface {
	Logger() log.Logger
	Close(ctx context.Context) error
	Config() Config
	Wait(ctx context.Context) error
	Unlock()
}
