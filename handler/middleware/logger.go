package middleware

import (
	"log/slog"

	"github.com/snekROmonoro/disgo/events"
	"github.com/snekROmonoro/disgo/handler"
)

var Logger handler.Middleware = func(next handler.Handler) handler.Handler {
	return func(e *events.InteractionCreate) error {
		e.Client().Logger().Info("handling interaction", slog.Int64("interaction_id", int64(e.Interaction.ID())))
		return next(e)
	}
}
