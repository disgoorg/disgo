// Package handler provides a way to handle interactions like application commands, autocomplete, buttons, select menus & modals with a simple interface.
//
// The handler package is inspired by the go-chi/chi http router.
// Each interaction has a path which is either the command name (starting with /) or the custom id. According to this path all interactions are routed to the correct handler.
// Slash Commands can have subcommands, which are nested paths. For example /test/subcommand1 or /test/subcommandgroup/subcommand.
//
// The handler also supports variables in its path which is especially useful for subcommands, components and modals.
// Variables are defined by curly braces like {variable} and can be accessed in the handler via the Variables map.
//
// You can also register middlewares, which are executed before the handler is called. Middlewares can be used to check permissions, validate input or do other things.
// Middlewares can also be attached to sub-routers, which is useful if you want to have a middleware for all subcommands of a command as an example.
// A middleware does not care which interaction type it is, it is just executed before the handler and has the following signature:
// type Middleware func(next func(e *events.InteractionCreate)) func(e *events.InteractionCreate)
//
// The handler iterates over all routes until it finds the fist matching route. If no route matches, the handler will call the NotFoundHandler.
// The NotFoundHandler can be set via the `NotFound` method on the *Mux. If no NotFoundHandler is set nothing will happen.

package handler

import (
	"errors"
	"strings"

	"github.com/snekROmonoro/snowflake"

	"github.com/snekROmonoro/disgo/bot"
	"github.com/snekROmonoro/disgo/discord"
	"github.com/snekROmonoro/disgo/events"
	"github.com/snekROmonoro/disgo/rest"
)

// SyncCommands sets the given commands for the given guilds or globally if no guildIDs are empty. It will return on the first error for multiple guilds.
func SyncCommands(client bot.Client, commands []discord.ApplicationCommandCreate, guildIDs []snowflake.ID, opts ...rest.RequestOpt) error {
	if len(guildIDs) == 0 {
		_, err := client.Rest().SetGlobalCommands(client.ApplicationID(), commands, opts...)
		return err
	}
	for _, guildID := range guildIDs {
		_, err := client.Rest().SetGuildCommands(client.ApplicationID(), guildID, commands, opts...)
		if err != nil {
			return err
		}
	}
	return nil
}

type handlerHolder[T any] struct {
	pattern string
	handler T
	t       discord.InteractionType
}

func (h *handlerHolder[T]) Match(path string, t discord.InteractionType) bool {
	if h.t != t {
		return false
	}
	parts := splitPath(path)
	patternParts := splitPath(h.pattern)

	for i, part := range patternParts {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			continue
		}
		if len(parts) <= i || part != parts[i] {
			return false
		}
	}

	return true
}

func (h *handlerHolder[T]) Handle(path string, variables map[string]string, event *events.InteractionCreate) error {
	parseVariables(path, h.pattern, variables)

	switch handler := any(h.handler).(type) {
	case CommandHandler:
		return handler(&CommandEvent{
			ApplicationCommandInteractionCreate: &events.ApplicationCommandInteractionCreate{
				GenericEvent:                  event.GenericEvent,
				ApplicationCommandInteraction: event.Interaction.(discord.ApplicationCommandInteraction),
				Respond:                       event.Respond,
			},
			Variables: variables,
		})
	case AutocompleteHandler:
		return handler(&AutocompleteEvent{
			AutocompleteInteractionCreate: &events.AutocompleteInteractionCreate{
				GenericEvent:            event.GenericEvent,
				AutocompleteInteraction: event.Interaction.(discord.AutocompleteInteraction),
				Respond:                 event.Respond,
			},
			Variables: variables,
		})
	case ComponentHandler:
		return handler(&ComponentEvent{
			ComponentInteractionCreate: &events.ComponentInteractionCreate{
				GenericEvent:         event.GenericEvent,
				ComponentInteraction: event.Interaction.(discord.ComponentInteraction),
				Respond:              event.Respond,
			},
			Variables: variables,
		})
	case ModalHandler:
		return handler(&ModalEvent{
			ModalSubmitInteractionCreate: &events.ModalSubmitInteractionCreate{
				GenericEvent:           event.GenericEvent,
				ModalSubmitInteraction: event.Interaction.(discord.ModalSubmitInteraction),
				Respond:                event.Respond,
			},
			Variables: variables,
		})
	}
	return errors.New("unknown handler type")
}
