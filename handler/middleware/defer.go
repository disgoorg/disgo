package middleware

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

// Defer is a middleware that defers the specified interaction type.
// If updateMessage is true, it will respond with discord.InteractionResponseTypeDeferredUpdateMessage instead of discord.InteractionResponseTypeDeferredCreateMessage.
// If ephemeral is true, it will respond with discord.MessageFlagEphemeral flag in case of a discord.InteractionResponseTypeDeferredCreateMessage.
// Note: You can use this middleware multiple times with different interaction types.
// Note: You can use this middleware in combination with the Go middleware to defer & run in a goroutine.
func Defer(interactionType discord.InteractionType, updateMessage bool, ephemeral bool) handler.Middleware {
	return func(next handler.Handler) handler.Handler {
		return func(event *handler.InteractionEvent) error {
			if event.Type() == interactionType {
				responseType := discord.InteractionResponseTypeDeferredCreateMessage
				if updateMessage {
					responseType = discord.InteractionResponseTypeDeferredUpdateMessage
				}

				var data discord.InteractionResponseData
				if ephemeral && !updateMessage {
					data = discord.MessageCreate{
						Flags: discord.MessageFlagEphemeral,
					}
				}
				if err := event.Respond(responseType, data); err != nil {
					return err
				}
			}
			return next(event)
		}
	}
}
