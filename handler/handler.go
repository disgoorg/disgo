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
	"errors"
	"slices"
	"strings"

	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/rest"
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
	t2      []int
}

func (h *handlerHolder[T]) Match(path string, t discord.InteractionType, t2 int) bool {
	if h.t != t || (len(h.t2) > 0 && !slices.Contains(h.t2, t2)) {
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

func (h *handlerHolder[T]) Handle(path string, event *InteractionEvent) error {
	parseVariables(path, h.pattern, event.Vars)

	switch handler := any(h.handler).(type) {
	case InteractionHandler:
		return handler(event)
	case CommandHandler:
		return handler(&CommandEvent{
			ApplicationCommandInteractionCreate: &events.ApplicationCommandInteractionCreate{
				GenericEvent:                  event.GenericEvent,
				ApplicationCommandInteraction: event.Interaction.(discord.ApplicationCommandInteraction),
				Respond:                       event.Respond,
			},
			Vars: event.Vars,
			Ctx:  event.Ctx,
		})
	case SlashCommandHandler:
		commandInteraction := event.Interaction.(discord.ApplicationCommandInteraction)
		return handler(commandInteraction.Data.(discord.SlashCommandInteractionData), &CommandEvent{
			ApplicationCommandInteractionCreate: &events.ApplicationCommandInteractionCreate{
				GenericEvent:                  event.GenericEvent,
				ApplicationCommandInteraction: commandInteraction,
				Respond:                       event.Respond,
			},
			Vars: event.Vars,
			Ctx:  event.Ctx,
		})
	case UserCommandHandler:
		commandInteraction := event.Interaction.(discord.ApplicationCommandInteraction)
		return handler(commandInteraction.Data.(discord.UserCommandInteractionData), &CommandEvent{
			ApplicationCommandInteractionCreate: &events.ApplicationCommandInteractionCreate{
				GenericEvent:                  event.GenericEvent,
				ApplicationCommandInteraction: commandInteraction,
				Respond:                       event.Respond,
			},
			Vars: event.Vars,
			Ctx:  event.Ctx,
		})
	case MessageCommandHandler:
		commandInteraction := event.Interaction.(discord.ApplicationCommandInteraction)
		return handler(commandInteraction.Data.(discord.MessageCommandInteractionData), &CommandEvent{
			ApplicationCommandInteractionCreate: &events.ApplicationCommandInteractionCreate{
				GenericEvent:                  event.GenericEvent,
				ApplicationCommandInteraction: commandInteraction,
				Respond:                       event.Respond,
			},
			Vars: event.Vars,
			Ctx:  event.Ctx,
		})
	case EntryPointCommandHandler:
		commandInteraction := event.Interaction.(discord.ApplicationCommandInteraction)
		return handler(commandInteraction.Data.(discord.EntryPointCommandInteractionData), &CommandEvent{
			ApplicationCommandInteractionCreate: &events.ApplicationCommandInteractionCreate{
				GenericEvent:                  event.GenericEvent,
				ApplicationCommandInteraction: commandInteraction,
				Respond:                       event.Respond,
			},
			Vars: event.Vars,
			Ctx:  event.Ctx,
		})
	case AutocompleteHandler:
		return handler(&AutocompleteEvent{
			AutocompleteInteractionCreate: &events.AutocompleteInteractionCreate{
				GenericEvent:            event.GenericEvent,
				AutocompleteInteraction: event.Interaction.(discord.AutocompleteInteraction),
				Respond:                 event.Respond,
			},
			Vars: event.Vars,
			Ctx:  event.Ctx,
		})
	case ComponentHandler:
		return handler(&ComponentEvent{
			ComponentInteractionCreate: &events.ComponentInteractionCreate{
				GenericEvent:         event.GenericEvent,
				ComponentInteraction: event.Interaction.(discord.ComponentInteraction),
				Respond:              event.Respond,
			},
			Vars: event.Vars,
			Ctx:  event.Ctx,
		})
	case ButtonComponentHandler:
		componentInteraction := event.Interaction.(discord.ComponentInteraction)
		return handler(componentInteraction.Data.(discord.ButtonInteractionData), &ComponentEvent{
			ComponentInteractionCreate: &events.ComponentInteractionCreate{
				GenericEvent:         event.GenericEvent,
				ComponentInteraction: componentInteraction,
				Respond:              event.Respond,
			},
			Vars: event.Vars,
			Ctx:  event.Ctx,
		})
	case SelectMenuComponentHandler:
		componentInteraction := event.Interaction.(discord.ComponentInteraction)
		return handler(componentInteraction.Data.(discord.SelectMenuInteractionData), &ComponentEvent{
			ComponentInteractionCreate: &events.ComponentInteractionCreate{
				GenericEvent:         event.GenericEvent,
				ComponentInteraction: componentInteraction,
				Respond:              event.Respond,
			},
			Vars: event.Vars,
			Ctx:  event.Ctx,
		})
	case ModalHandler:
		return handler(&ModalEvent{
			ModalSubmitInteractionCreate: &events.ModalSubmitInteractionCreate{
				GenericEvent:           event.GenericEvent,
				ModalSubmitInteraction: event.Interaction.(discord.ModalSubmitInteraction),
				Respond:                event.Respond,
			},
			Vars: event.Vars,
			Ctx:  event.Ctx,
		})
	}
	return errors.New("unknown handler type")
}
