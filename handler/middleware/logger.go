package middleware

import (
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/handler"
)

var Logger handler.Middleware = func(next handler.Handler) handler.Handler {
	return func(e *events.InteractionCreate) error {
		e.Client().Logger().Infof("handling interaction: %s\n", e.Interaction.ID())
		return next(e)
	}
}
