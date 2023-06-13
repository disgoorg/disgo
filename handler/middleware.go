package handler

import (
	"context"

	"github.com/disgoorg/disgo/events"
)

type (
	Handler func(ctx context.Context, e *events.InteractionCreate) error

	Middleware func(next Handler) Handler

	Middlewares []Middleware
)
