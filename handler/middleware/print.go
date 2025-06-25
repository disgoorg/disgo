package middleware

import (
	"github.com/disgoorg/disgo/handler"
)

// Print is a middleware that prints the specified message with the specified arguments.
func Print(msg string, args ...any) handler.Middleware {
	return func(next handler.Handler) handler.Handler {
		return func(event *handler.InteractionEvent) error {
			event.Client().Logger.InfoContext(event.Ctx, msg, args...)
			return next(event)
		}
	}
}
