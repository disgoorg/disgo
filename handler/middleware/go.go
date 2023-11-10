package middleware

import (
	"log/slog"

	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/handler"
)

// Go is a middleware that runs the next handler in a goroutine
var Go handler.Middleware = func(next handler.Handler) handler.Handler {
	return func(e *events.InteractionCreate) error {
		go func() {
			if err := next(e); err != nil {
				e.Client().Logger().Error("failed to handle interaction", slog.String("err", err.Error()))
			}
		}()
		return nil
	}
}
