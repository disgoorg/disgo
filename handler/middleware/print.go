package middleware

import (
	"context"

	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/handler"
)

func Print(content string) handler.Middleware {
	return func(next handler.Handler) handler.Handler {
		return func(ctx context.Context, event *events.InteractionCreate) error {
			println(content)
			return next(ctx, event)
		}
	}
}
