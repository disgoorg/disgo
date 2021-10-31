package events

import "github.com/DisgoOrg/disgo/core"

type InteractionCreateEvent struct {
	*GenericEvent
	core.Interaction
}

type ApplicationCommandInteractionCreateEvent struct {
	*GenericEvent
	core.ApplicationCommandInteraction
}

// SlashCommandEvent indicates that a slash core.ApplicationCommand was run
type SlashCommandEvent struct {
	*GenericEvent
	*core.SlashCommandInteraction
}

type UserCommandEvent struct {
	*GenericEvent
	*core.UserCommandInteraction
}

type MessageCommandEvent struct {
	*GenericEvent
	*core.MessageCommandInteraction
}

type ComponentInteractionCreateEvent struct {
	*GenericEvent
	core.ComponentInteraction
}

// ButtonClickEvent indicates that a core.ButtonComponent was clicked
type ButtonClickEvent struct {
	*GenericEvent
	*core.ButtonInteraction
}

// SelectMenuSubmitEvent indicates that a core.SelectMenuComponent was submitted
type SelectMenuSubmitEvent struct {
	*GenericEvent
	*core.SelectMenuInteraction
}

type AutocompleteEvent struct {
	*GenericEvent
	*core.AutocompleteInteraction
}
