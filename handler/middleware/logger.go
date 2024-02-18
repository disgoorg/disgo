package middleware

import (
	"log/slog"

	"github.com/disgoorg/disgo/handler"
)

var Logger handler.Middleware = func(next handler.Handler) handler.Handler {
	return func(event *handler.InteractionEvent) error {
		event.Client().Logger().Info("handling interaction", slog.Int64("interaction_id", int64(event.Interaction.ID())))
		return next(event)
	}
}
