package events

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
)

type InteractionResponderFunc func(callbackType discord.InteractionCallbackType, data discord.InteractionCallbackData, opts ...rest.RequestOpt) error

type InteractionEvent struct {
	*GenericEvent
	discord.Interaction
	Respond InteractionResponderFunc
}

type ApplicationCommandInteractionEvent struct {
	*GenericEvent
	discord.ApplicationCommandInteraction
	Respond InteractionResponderFunc
}

func (e *ApplicationCommandInteractionEvent) CreateMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionCallbackTypeCreateMessage, messageCreate, opts...)
}

func (e *ApplicationCommandInteractionEvent) DeferCreateMessage(ephemeral bool, opts ...rest.RequestOpt) error {
	var messageCreate discord.MessageCreate
	if ephemeral {
		messageCreate = discord.MessageCreate{Flags: discord.MessageFlagEphemeral}
	}
	return e.Respond(discord.InteractionCallbackTypeDeferredCreateMessage, messageCreate, opts...)
}

func (e *ApplicationCommandInteractionEvent) CreateModal(modalCreate discord.ModalCreate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionCallbackTypeModal, modalCreate, opts...)
}

type ComponentInteractionEvent struct {
	*GenericEvent
	discord.ComponentInteraction
	Respond InteractionResponderFunc
}

func (e *ComponentInteractionEvent) CreateMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionCallbackTypeCreateMessage, messageCreate, opts...)
}

func (e *ComponentInteractionEvent) DeferCreateMessage(ephemeral bool, opts ...rest.RequestOpt) error {
	var messageCreate discord.MessageCreate
	if ephemeral {
		messageCreate = discord.MessageCreate{Flags: discord.MessageFlagEphemeral}
	}
	return e.Respond(discord.InteractionCallbackTypeDeferredCreateMessage, messageCreate, opts...)
}

func (e *ComponentInteractionEvent) UpdateMessage(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionCallbackTypeUpdateMessage, messageUpdate, opts...)
}

func (e *ComponentInteractionEvent) DeferUpdateMessage(opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionCallbackTypeDeferredUpdateMessage, nil, opts...)
}

func (e *ComponentInteractionEvent) CreateModal(modalCreate discord.ModalCreate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionCallbackTypeModal, modalCreate, opts...)
}

type AutocompleteInteractionEvent struct {
	*GenericEvent
	discord.AutocompleteInteraction
	Respond InteractionResponderFunc
}

func (e *AutocompleteInteractionEvent) Result(choices []discord.AutocompleteChoice, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionCallbackTypeApplicationCommandAutocompleteResult, discord.AutocompleteResult{Choices: choices}, opts...)
}

type ModalSubmitInteractionEvent struct {
	*GenericEvent
	discord.ModalSubmitInteraction
	Respond InteractionResponderFunc
}

func (e *ModalSubmitInteractionEvent) CreateMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionCallbackTypeCreateMessage, messageCreate, opts...)
}

func (e *ModalSubmitInteractionEvent) DeferCreateMessage(ephemeral bool, opts ...rest.RequestOpt) error {
	var messageCreate discord.MessageCreate
	if ephemeral {
		messageCreate = discord.MessageCreate{Flags: discord.MessageFlagEphemeral}
	}
	return e.Respond(discord.InteractionCallbackTypeDeferredCreateMessage, messageCreate, opts...)
}

func (e *ModalSubmitInteractionEvent) UpdateMessage(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionCallbackTypeUpdateMessage, messageUpdate, opts...)
}

func (e *ModalSubmitInteractionEvent) DeferUpdateMessage(opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionCallbackTypeDeferredUpdateMessage, nil, opts...)
}
