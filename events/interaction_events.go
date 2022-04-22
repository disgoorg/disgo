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

func (e *InteractionEvent) Guild() (discord.Guild, bool) {
	if e.GuildID() != nil {
		return e.Client().Caches().Guilds().Get(*e.GuildID())
	}
	return discord.Guild{}, false
}

func (e *InteractionEvent) Channel() (discord.MessageChannel, bool) {
	return e.Client().Caches().Channels().GetMessageChannel(e.ChannelID())
}

func (e *InteractionEvent) DMChannel() (discord.DMChannel, bool) {
	return e.Client().Caches().Channels().GetDMChannel(e.ChannelID())
}

func (e *InteractionEvent) GuildChannel() (discord.GuildMessageChannel, bool) {
	return e.Client().Caches().Channels().GetGuildMessageChannel(e.ChannelID())
}

type ApplicationCommandInteractionEvent struct {
	*GenericEvent
	discord.ApplicationCommandInteraction
	Respond InteractionResponderFunc
}

func (e *ApplicationCommandInteractionEvent) Guild() (discord.Guild, bool) {
	if e.GuildID() != nil {
		return e.Client().Caches().Guilds().Get(*e.GuildID())
	}
	return discord.Guild{}, false
}

func (e *ApplicationCommandInteractionEvent) Channel() (discord.MessageChannel, bool) {
	return e.Client().Caches().Channels().GetMessageChannel(e.ChannelID())
}

func (e *ApplicationCommandInteractionEvent) DMChannel() (discord.DMChannel, bool) {
	return e.Client().Caches().Channels().GetDMChannel(e.ChannelID())
}

func (e *ApplicationCommandInteractionEvent) GuildChannel() (discord.GuildMessageChannel, bool) {
	return e.Client().Caches().Channels().GetGuildMessageChannel(e.ChannelID())
}

func (e *ApplicationCommandInteractionEvent) CreateMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionCallbackTypeCreateMessage, messageCreate, opts...)
}

func (e *ApplicationCommandInteractionEvent) DeferCreateMessage(ephemeral bool, opts ...rest.RequestOpt) error {
	var data discord.InteractionCallbackData
	if ephemeral {
		data = discord.MessageCreate{Flags: discord.MessageFlagEphemeral}
	}
	return e.Respond(discord.InteractionCallbackTypeDeferredCreateMessage, data, opts...)
}

func (e *ApplicationCommandInteractionEvent) CreateModal(modalCreate discord.ModalCreate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionCallbackTypeModal, modalCreate, opts...)
}

type ComponentInteractionEvent struct {
	*GenericEvent
	discord.ComponentInteraction
	Respond InteractionResponderFunc
}

func (e *ComponentInteractionEvent) Guild() (discord.Guild, bool) {
	if e.GuildID() != nil {
		return e.Client().Caches().Guilds().Get(*e.GuildID())
	}
	return discord.Guild{}, false
}

func (e *ComponentInteractionEvent) Channel() (discord.MessageChannel, bool) {
	return e.Client().Caches().Channels().GetMessageChannel(e.ChannelID())
}

func (e *ComponentInteractionEvent) DMChannel() (discord.DMChannel, bool) {
	return e.Client().Caches().Channels().GetDMChannel(e.ChannelID())
}

func (e *ComponentInteractionEvent) GuildChannel() (discord.GuildMessageChannel, bool) {
	return e.Client().Caches().Channels().GetGuildMessageChannel(e.ChannelID())
}

func (e *ComponentInteractionEvent) CreateMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionCallbackTypeCreateMessage, messageCreate, opts...)
}

func (e *ComponentInteractionEvent) DeferCreateMessage(ephemeral bool, opts ...rest.RequestOpt) error {
	var data discord.InteractionCallbackData
	if ephemeral {
		data = discord.MessageCreate{Flags: discord.MessageFlagEphemeral}
	}
	return e.Respond(discord.InteractionCallbackTypeDeferredCreateMessage, data, opts...)
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

func (e *AutocompleteInteractionEvent) Guild() (discord.Guild, bool) {
	if e.GuildID() != nil {
		return e.Client().Caches().Guilds().Get(*e.GuildID())
	}
	return discord.Guild{}, false
}

func (e *AutocompleteInteractionEvent) Channel() (discord.MessageChannel, bool) {
	return e.Client().Caches().Channels().GetMessageChannel(e.ChannelID())
}

func (e *AutocompleteInteractionEvent) DMChannel() (discord.DMChannel, bool) {
	return e.Client().Caches().Channels().GetDMChannel(e.ChannelID())
}

func (e *AutocompleteInteractionEvent) GuildChannel() (discord.GuildMessageChannel, bool) {
	return e.Client().Caches().Channels().GetGuildMessageChannel(e.ChannelID())
}

func (e *AutocompleteInteractionEvent) Result(choices []discord.AutocompleteChoice, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionCallbackTypeApplicationCommandAutocompleteResult, discord.AutocompleteResult{Choices: choices}, opts...)
}

type ModalSubmitInteractionEvent struct {
	*GenericEvent
	discord.ModalSubmitInteraction
	Respond InteractionResponderFunc
}

func (e *ModalSubmitInteractionEvent) Guild() (discord.Guild, bool) {
	if e.GuildID() != nil {
		return e.Client().Caches().Guilds().Get(*e.GuildID())
	}
	return discord.Guild{}, false
}

func (e *ModalSubmitInteractionEvent) Channel() (discord.MessageChannel, bool) {
	return e.Client().Caches().Channels().GetMessageChannel(e.ChannelID())
}

func (e *ModalSubmitInteractionEvent) DMChannel() (discord.DMChannel, bool) {
	return e.Client().Caches().Channels().GetDMChannel(e.ChannelID())
}

func (e *ModalSubmitInteractionEvent) GuildChannel() (discord.GuildMessageChannel, bool) {
	return e.Client().Caches().Channels().GetGuildMessageChannel(e.ChannelID())
}

func (e *ModalSubmitInteractionEvent) CreateMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionCallbackTypeCreateMessage, messageCreate, opts...)
}

func (e *ModalSubmitInteractionEvent) DeferCreateMessage(ephemeral bool, opts ...rest.RequestOpt) error {
	var data discord.InteractionCallbackData
	if ephemeral {
		data = discord.MessageCreate{Flags: discord.MessageFlagEphemeral}
	}
	return e.Respond(discord.InteractionCallbackTypeDeferredCreateMessage, data, opts...)
}

func (e *ModalSubmitInteractionEvent) UpdateMessage(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionCallbackTypeUpdateMessage, messageUpdate, opts...)
}

func (e *ModalSubmitInteractionEvent) DeferUpdateMessage(opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionCallbackTypeDeferredUpdateMessage, nil, opts...)
}
