package events

import "github.com/DisgoOrg/disgo/core"

type InteractionEvent struct {
	*GenericEvent
	core.Interaction
}

type ApplicationCommandInteractionEvent struct {
	*GenericEvent
	*core.ApplicationCommandInteraction
}

type ComponentInteractionEvent struct {
	*GenericEvent
	*core.ComponentInteraction
}

type AutocompleteInteractionEvent struct {
	*GenericEvent
	*core.AutocompleteInteraction
}

type ModalSubmitEvent struct {
	*GenericEvent
	*core.ModalSubmitInteraction
}