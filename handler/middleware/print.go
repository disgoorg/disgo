package middleware

import (
	"github.com/snekROmonoro/disgo/events"
	"github.com/snekROmonoro/disgo/handler"
)

func Print(content string) handler.Middleware {
	return func(next handler.Handler) handler.Handler {
		return func(event *events.InteractionCreate) error {
			println(content)
			return next(event)
		}
	}
}
