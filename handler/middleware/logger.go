package middleware

import (
	"log/slog"

	"github.com/disgoorg/disgo/handler"
)

// Logger is a middleware that logs the interaction and its variables.
var Logger handler.Middleware = func(next handler.Handler) handler.Handler {
	return func(event *handler.InteractionEvent) error {
		event.Client().Logger.InfoContext(event.Ctx, "handling interaction", slog.Any("interaction", event.Interaction), slog.Any("vars", event.Vars))
		return next(event)
	}
}
