package events

import (
	"github.com/DisgoOrg/disgo/discord"
)

type InteractionEvent struct {
	*GenericEvent
	discord.Interaction
}

type ApplicationCommandInteractionEvent struct {
	*GenericEvent
	discord.ApplicationCommandInteraction
}

type ComponentInteractionEvent struct {
	*GenericEvent
	discord.ComponentInteraction
}

type AutocompleteInteractionEvent struct {
	*GenericEvent
	discord.AutocompleteInteraction
}

type ModalSubmitInteractionEvent struct {
	*GenericEvent
	discord.ModalSubmitInteraction
}
