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

// SlashCommandEvent indicates that a slash discord.ApplicationCommand was run
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

// ButtonClickEvent indicates that a discord.ButtonComponent was clicked
type ButtonClickEvent struct {
	*GenericEvent
	*core.ButtonInteraction
}

// SelectMenuSubmitEvent indicates that a discord.SelectMenuComponent was submitted
type SelectMenuSubmitEvent struct {
	*GenericEvent
	*core.SelectMenuInteraction
}

type AutocompleteEvent struct {
	*GenericEvent
	*core.AutocompleteInteraction
}

type ModalSubmitEvent struct {
	*GenericEvent
	*core.ModalSubmitInteraction
}