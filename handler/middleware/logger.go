package middleware

import (
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/handler"
)

var Logger handler.Middleware = func(next handler.Handler) handler.Handler {
	return func(event *events.InteractionCreate) {
		event.Client().Logger().Infof("handling interaction: %s\n", event.Interaction.ID())
		next(event)
	}
}
