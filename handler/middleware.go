package handler

import (
	"github.com/disgoorg/disgo/events"
)

var EmptyMiddleware = func(next Handler) Handler {
	return next
}

type (
	Handler func(event *events.InteractionCreate)

	Middleware func(next Handler) Handler

	Middlewares []Middleware
)
