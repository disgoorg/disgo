// Package handler provides a way to handle interactions like application commands, autocomplete, buttons, select menus & modals with a simple interface.
//
// The handler package is inspired by the go-chi/chi http router.
// Each interaction has a path which is either the command name (starting with /) or the custom id. According to this path all interactions are routed to the correct handler.
// Slash Commands can have subcommands, which are nested paths. For example /test/subcommand1 or /test/subcommandgroup/subcommand.
//
// The handler also supports variables in its path which is especially useful for subcommands, components and modals.
// Vars are defined by curly braces like {variable} and can be accessed in the handler via the Vars map.
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
	"context"
	"errors"
	"slices"
	"strings"

	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/rest"
)

// SyncCommands sets the given commands for the given guilds or globally if guildIDs is empty. It will return on the first error for multiple guilds.
func SyncCommands(client *bot.Client, commands []discord.ApplicationCommandCreate, guildIDs []snowflake.ID, opts ...rest.RequestOpt) error {
	if len(guildIDs) == 0 {
		_, err := client.Rest.SetGlobalCommands(client.ApplicationID, commands, opts...)
		return err
	}
	for _, guildID := range guildIDs {
		_, err := client.Rest.SetGuildCommands(client.ApplicationID, guildID, commands, opts...)
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
	t2      []int
}

func (h *handlerHolder[T]) Match(path string, t discord.InteractionType, t2 int) bool {
	if h.t > 0 && (h.t != t || (len(h.t2) > 0 && !slices.Contains(h.t2, t2))) {
		return false
	}
	parts := splitPath(path)
	patternParts := splitPath(h.pattern)
	if len(parts) < len(patternParts) {
		return false
	}

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

func (h *handlerHolder[T]) Handle(path string, event *InteractionEvent) error {
	parseVariables(path, h.pattern, event.Vars)

	// Re-wrap event.Respond with the typed event's Ctx so that reassigning
	// e.Ctx inside a handler is picked up. This double-wraps the OnEvent
	// wrapper but the last WithCtx wins so the behavior is correct.

	switch handler := any(h.handler).(type) {
	case InteractionHandler:
		return handler(event)
	case CommandHandler, SlashCommandHandler, UserCommandHandler, MessageCommandHandler, EntryPointCommandHandler:
		e := &CommandEvent{
			ApplicationCommandInteractionCreate: &events.ApplicationCommandInteractionCreate{
				GenericEvent:                  event.GenericEvent,
				ApplicationCommandInteraction: event.Interaction.(discord.ApplicationCommandInteraction),
			},
			Vars: event.Vars,
			Ctx:  event.Ctx,
		}
		e.Respond = wrapRespond(event.Respond, &e.Ctx)
		switch handler := any(h.handler).(type) {
		case CommandHandler:
			return handler(e)
		case SlashCommandHandler:
			return handler(e.SlashCommandInteractionData(), e)
		case UserCommandHandler:
			return handler(e.UserCommandInteractionData(), e)
		case MessageCommandHandler:
			return handler(e.MessageCommandInteractionData(), e)
		case EntryPointCommandHandler:
			return handler(e.EntryPointCommandInteractionData(), e)
		}
	case AutocompleteHandler:
		e := &AutocompleteEvent{
			AutocompleteInteractionCreate: &events.AutocompleteInteractionCreate{
				GenericEvent:            event.GenericEvent,
				AutocompleteInteraction: event.Interaction.(discord.AutocompleteInteraction),
			},
			Vars: event.Vars,
			Ctx:  event.Ctx,
		}
		e.Respond = wrapRespond(event.Respond, &e.Ctx)
		return handler(e)
	case ComponentHandler, ButtonComponentHandler, SelectMenuComponentHandler:
		e := &ComponentEvent{
			ComponentInteractionCreate: &events.ComponentInteractionCreate{
				GenericEvent:         event.GenericEvent,
				ComponentInteraction: event.Interaction.(discord.ComponentInteraction),
			},
			Vars: event.Vars,
			Ctx:  event.Ctx,
		}
		e.Respond = wrapRespond(event.Respond, &e.Ctx)
		switch handler := any(h.handler).(type) {
		case ComponentHandler:
			return handler(e)
		case ButtonComponentHandler:
			return handler(e.ButtonInteractionData(), e)
		case SelectMenuComponentHandler:
			return handler(e.SelectMenuInteractionData(), e)
		}
	case ModalHandler:
		e := &ModalEvent{
			ModalSubmitInteractionCreate: &events.ModalSubmitInteractionCreate{
				GenericEvent:           event.GenericEvent,
				ModalSubmitInteraction: event.Interaction.(discord.ModalSubmitInteraction),
			},
			Vars: event.Vars,
			Ctx:  event.Ctx,
		}
		e.Respond = wrapRespond(event.Respond, &e.Ctx)
		return handler(e)
	}
	return errors.New("unknown handler type")
}

func wrapRespond(respond events.InteractionResponderFunc, ctx *context.Context) events.InteractionResponderFunc {
	if respond == nil {
		return nil
	}
	return func(responseType discord.InteractionResponseType, data discord.InteractionResponseData, opts ...rest.RequestOpt) error {
		return respond(responseType, data, append([]rest.RequestOpt{rest.WithCtx(*ctx)}, opts...)...)
	}
}
