package handler

// Handler is a function that handles an interaction event.
type Handler func(e *InteractionEvent) error

type (
	// Middleware is a function that wraps a handler to intercept and short-circuit the interactions.
	Middleware func(next Handler) Handler

	// Middlewares is a list of middlewares.
	Middlewares []Middleware
)
