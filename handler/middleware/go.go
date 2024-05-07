package middleware

import (
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

// Go is a middleware that runs the next handler in a goroutine.
var Go = GoErr(func(event *handler.InteractionEvent, err error) {
	event.Client().Logger().Error("failed to handle interaction", slog.Any("err", err))
})

// GoDefer combines Go and Defer
func GoDefer(interactionType discord.InteractionType, updateMessage bool, ephemeral bool) handler.Middleware {
	return func(next handler.Handler) handler.Handler {
		return Go(Defer(interactionType, updateMessage, ephemeral)(next))
	}
}

// GoErr is a middleware that runs the next handler in a goroutine and lets you handle the error which may occur.
func GoErr(h handler.ErrorHandler) handler.Middleware {
	return func(next handler.Handler) handler.Handler {
		return func(event *handler.InteractionEvent) error {
			go func() {
				if err := next(event); err != nil {
					h(event, err)
				}
			}()
			return nil
		}
	}
}

// GoErrDefer combines GoErr and Defer
func GoErrDefer(h handler.ErrorHandler, interactionType discord.InteractionType, updateMessage bool, ephemeral bool) handler.Middleware {
	return func(next handler.Handler) handler.Handler {
		return GoErr(h)(Defer(interactionType, updateMessage, ephemeral)(next))
	}
}
