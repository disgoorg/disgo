package handler

type Handler func(e *InteractionEvent) error

type (
	Middleware func(next Handler) Handler

	Middlewares []Middleware
)
