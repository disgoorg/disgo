package handler

import (
	"errors"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

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
		if part != parts[i] {
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
