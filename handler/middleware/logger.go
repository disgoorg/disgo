package middleware

import (
	"context"
	"log/slog"

	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/handler"
)

var Logger handler.Middleware = func(next handler.Handler) handler.Handler {
	return func(ctx context.Context, e *events.InteractionCreate) error {
		e.Client().Logger().Info("handling interaction", slog.Int64("interaction_id", int64(e.Interaction.ID())))
		return next(ctx, e)
	}
}
