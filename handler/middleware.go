package handler

import (
	"github.com/disgoorg/disgo/events"
)

type (
	Handler func(e *events.InteractionCreate)

	Middleware func(next Handler) Handler

	Middlewares []Middleware
)
