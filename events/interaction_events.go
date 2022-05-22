package events

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
)

type InteractionResponderFunc func(callbackType discord.InteractionCallbackType, data discord.InteractionCallbackData, opts ...rest.RequestOpt) error

type InteractionCreate struct {
	*GenericEvent
	discord.Interaction
	Respond InteractionResponderFunc
}

func (e *InteractionCreate) Guild() (discord.Guild, bool) {
	if e.GuildID() != nil {
		return e.Client().Caches().Guilds().Get(*e.GuildID())
	}
	return discord.Guild{}, false
}

func (e *InteractionCreate) Channel() (discord.MessageChannel, bool) {
	return e.Client().Caches().Channels().GetMessageChannel(e.ChannelID())
}

func (e *InteractionCreate) DMChannel() (discord.DMChannel, bool) {
	return e.Client().Caches().Channels().GetDMChannel(e.ChannelID())
}

func (e *InteractionCreate) GuildChannel() (discord.GuildMessageChannel, bool) {
	return e.Client().Caches().Channels().GetGuildMessageChannel(e.ChannelID())
}

type ApplicationCommandInteractionCreate struct {
	*GenericEvent
	discord.ApplicationCommandInteraction
	Respond InteractionResponderFunc
}

func (e *ApplicationCommandInteractionCreate) Guild() (discord.Guild, bool) {
	if e.GuildID() != nil {
		return e.Client().Caches().Guilds().Get(*e.GuildID())
	}
	return discord.Guild{}, false
}

func (e *ApplicationCommandInteractionCreate) Channel() (discord.MessageChannel, bool) {
	return e.Client().Caches().Channels().GetMessageChannel(e.ChannelID())
}

func (e *ApplicationCommandInteractionCreate) DMChannel() (discord.DMChannel, bool) {
	return e.Client().Caches().Channels().GetDMChannel(e.ChannelID())
}

func (e *ApplicationCommandInteractionCreate) GuildChannel() (discord.GuildMessageChannel, bool) {
	return e.Client().Caches().Channels().GetGuildMessageChannel(e.ChannelID())
}

func (e *ApplicationCommandInteractionCreate) CreateMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionCallbackTypeCreateMessage, messageCreate, opts...)
}

func (e *ApplicationCommandInteractionCreate) DeferCreateMessage(ephemeral bool, opts ...rest.RequestOpt) error {
	var data discord.InteractionCallbackData
	if ephemeral {
		data = discord.MessageCreate{Flags: discord.MessageFlagEphemeral}
	}
	return e.Respond(discord.InteractionCallbackTypeDeferredCreateMessage, data, opts...)
}

func (e *ApplicationCommandInteractionCreate) CreateModal(modalCreate discord.ModalCreate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionCallbackTypeModal, modalCreate, opts...)
}

type ComponentInteractionCreate struct {
	*GenericEvent
	discord.ComponentInteraction
	Respond InteractionResponderFunc
}

func (e *ComponentInteractionCreate) Guild() (discord.Guild, bool) {
	if e.GuildID() != nil {
		return e.Client().Caches().Guilds().Get(*e.GuildID())
	}
	return discord.Guild{}, false
}

func (e *ComponentInteractionCreate) Channel() (discord.MessageChannel, bool) {
	return e.Client().Caches().Channels().GetMessageChannel(e.ChannelID())
}

func (e *ComponentInteractionCreate) DMChannel() (discord.DMChannel, bool) {
	return e.Client().Caches().Channels().GetDMChannel(e.ChannelID())
}

func (e *ComponentInteractionCreate) GuildChannel() (discord.GuildMessageChannel, bool) {
	return e.Client().Caches().Channels().GetGuildMessageChannel(e.ChannelID())
}

func (e *ComponentInteractionCreate) CreateMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionCallbackTypeCreateMessage, messageCreate, opts...)
}

func (e *ComponentInteractionCreate) DeferCreateMessage(ephemeral bool, opts ...rest.RequestOpt) error {
	var data discord.InteractionCallbackData
	if ephemeral {
		data = discord.MessageCreate{Flags: discord.MessageFlagEphemeral}
	}
	return e.Respond(discord.InteractionCallbackTypeDeferredCreateMessage, data, opts...)
}

func (e *ComponentInteractionCreate) UpdateMessage(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionCallbackTypeUpdateMessage, messageUpdate, opts...)
}

func (e *ComponentInteractionCreate) DeferUpdateMessage(opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionCallbackTypeDeferredUpdateMessage, nil, opts...)
}

func (e *ComponentInteractionCreate) CreateModal(modalCreate discord.ModalCreate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionCallbackTypeModal, modalCreate, opts...)
}

type AutocompleteInteractionCreate struct {
	*GenericEvent
	discord.AutocompleteInteraction
	Respond InteractionResponderFunc
}

func (e *AutocompleteInteractionCreate) Guild() (discord.Guild, bool) {
	if e.GuildID() != nil {
		return e.Client().Caches().Guilds().Get(*e.GuildID())
	}
	return discord.Guild{}, false
}

func (e *AutocompleteInteractionCreate) Channel() (discord.MessageChannel, bool) {
	return e.Client().Caches().Channels().GetMessageChannel(e.ChannelID())
}

func (e *AutocompleteInteractionCreate) DMChannel() (discord.DMChannel, bool) {
	return e.Client().Caches().Channels().GetDMChannel(e.ChannelID())
}

func (e *AutocompleteInteractionCreate) GuildChannel() (discord.GuildMessageChannel, bool) {
	return e.Client().Caches().Channels().GetGuildMessageChannel(e.ChannelID())
}

func (e *AutocompleteInteractionCreate) Result(choices []discord.AutocompleteChoice, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionCallbackTypeApplicationCommandAutocompleteResult, discord.AutocompleteResult{Choices: choices}, opts...)
}

type ModalSubmitInteractionCreate struct {
	*GenericEvent
	discord.ModalSubmitInteraction
	Respond InteractionResponderFunc
}

func (e *ModalSubmitInteractionCreate) Guild() (discord.Guild, bool) {
	if e.GuildID() != nil {
		return e.Client().Caches().Guilds().Get(*e.GuildID())
	}
	return discord.Guild{}, false
}

func (e *ModalSubmitInteractionCreate) Channel() (discord.MessageChannel, bool) {
	return e.Client().Caches().Channels().GetMessageChannel(e.ChannelID())
}

func (e *ModalSubmitInteractionCreate) DMChannel() (discord.DMChannel, bool) {
	return e.Client().Caches().Channels().GetDMChannel(e.ChannelID())
}

func (e *ModalSubmitInteractionCreate) GuildChannel() (discord.GuildMessageChannel, bool) {
	return e.Client().Caches().Channels().GetGuildMessageChannel(e.ChannelID())
}

func (e *ModalSubmitInteractionCreate) CreateMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionCallbackTypeCreateMessage, messageCreate, opts...)
}

func (e *ModalSubmitInteractionCreate) DeferCreateMessage(ephemeral bool, opts ...rest.RequestOpt) error {
	var data discord.InteractionCallbackData
	if ephemeral {
		data = discord.MessageCreate{Flags: discord.MessageFlagEphemeral}
	}
	return e.Respond(discord.InteractionCallbackTypeDeferredCreateMessage, data, opts...)
}

func (e *ModalSubmitInteractionCreate) UpdateMessage(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionCallbackTypeUpdateMessage, messageUpdate, opts...)
}

func (e *ModalSubmitInteractionCreate) DeferUpdateMessage(opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionCallbackTypeDeferredUpdateMessage, nil, opts...)
}
