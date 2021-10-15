package events

import (
	"github.com/DisgoOrg/disgo/core"
)

type AutocompleteEvent struct {
	*GenericEvent
	*core.AutocompleteInteraction
}

func (e *AutocompleteEvent) Bot() *core.Bot {
	return e.bot
}

// SlashCommandEvent indicates that a slash core.ApplicationCommand was run
type SlashCommandEvent struct {
	*GenericEvent
	*core.SlashCommandInteraction
}

func (e *SlashCommandEvent) Bot() *core.Bot {
	return e.bot
}

type UserCommandEvent struct {
	*GenericEvent
	*core.UserCommandInteraction
}

func (e *UserCommandEvent) Bot() *core.Bot {
	return e.bot
}

type MessageCommandEvent struct {
	*GenericEvent
	*core.MessageCommandInteraction
}

func (e *MessageCommandEvent) Bot() *core.Bot {
	return e.bot
}
