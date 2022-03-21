package events

import (
	"github.com/DisgoOrg/disgo/discord"
)

type InteractionEvent struct {
	*GenericEvent
	discord.Interaction
	Respond func(callbackType discord.InteractionCallbackType, data discord.InteractionCallbackData) error
}

type ApplicationCommandInteractionEvent struct {
	*GenericEvent
	discord.ApplicationCommandInteraction
	Respond func(callbackType discord.InteractionCallbackType, data discord.CommandInteractionCallbackData) error
}

type ComponentInteractionEvent struct {
	*GenericEvent
	discord.ComponentInteraction
	Respond func(callbackType discord.InteractionCallbackType, data discord.ComponentInteractionCallbackData) error
}

type AutocompleteInteractionEvent struct {
	*GenericEvent
	discord.AutocompleteInteraction
	Respond func(data discord.AutocompleteResult) error
}

type ModalSubmitInteractionEvent struct {
	*GenericEvent
	discord.ModalSubmitInteraction
	Respond func(callbackType discord.InteractionCallbackType, data discord.ModalInteractionCallbackData) error
}
