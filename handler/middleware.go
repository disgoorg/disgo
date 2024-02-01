package handler

import (
	"github.com/snekROmonoro/disgo/events"
)

type (
	Handler func(e *events.InteractionCreate) error

	Middleware func(next Handler) Handler

	Middlewares []Middleware
)
