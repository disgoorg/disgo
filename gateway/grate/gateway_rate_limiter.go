package grate

import (
	"context"

	"github.com/disgoorg/log"
)

type Limiter interface {
	Logger() log.Logger
	Close(ctx context.Context)
	Reset()

	Wait(ctx context.Context) error
	Unlock()
}
