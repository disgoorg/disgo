package events

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
)

// InteractionResponderFunc is a function that can be used to respond to a discord.Interaction.
type InteractionResponderFunc func(responseType discord.InteractionResponseType, data discord.InteractionResponseData, opts ...rest.RequestOpt) error

// InteractionCreate indicates that a new interaction has been created.
type InteractionCreate struct {
	*GenericEvent
	discord.Interaction
	Respond InteractionResponderFunc
}

// Guild returns the guild that the interaction happened in if it happened in a guild.
// If the interaction happened in a DM, it returns nil.
// This only returns cached guilds.
func (e *InteractionCreate) Guild() (discord.Guild, bool) {
	if e.GuildID() != nil {
		return e.Client().Caches().Guilds().Get(*e.GuildID())
	}
	return discord.Guild{}, false
}

// Channel returns the discord.MessageChannel that the interaction happened in.
// This only returns cached channels.
func (e *InteractionCreate) Channel() (discord.MessageChannel, bool) {
	return e.Client().Caches().Channels().GetMessageChannel(e.ChannelID())
}

// DMChannel returns the discord.DMChannel that the interaction happened in.
// If the interaction happened in a guild, it returns nil.
// This only returns cached channels.
func (e *InteractionCreate) DMChannel() (discord.DMChannel, bool) {
	return e.Client().Caches().Channels().GetDMChannel(e.ChannelID())
}

// GuildChannel returns the discord.GuildMessageChannel that the interaction happened in.
// If the interaction happened in a dm, it returns nil.
// This only returns cached channels.
func (e *InteractionCreate) GuildChannel() (discord.GuildMessageChannel, bool) {
	return e.Client().Caches().Channels().GetGuildMessageChannel(e.ChannelID())
}

// ApplicationCommandInteractionCreate is the base struct for all application command interaction create events.
type ApplicationCommandInteractionCreate struct {
	*GenericEvent
	discord.ApplicationCommandInteraction
	Respond InteractionResponderFunc
}

// Guild returns the guild that the interaction happened in if it happened in a guild.
// If the interaction happened in a DM, it returns nil.
// This only returns cached guilds.
func (e *ApplicationCommandInteractionCreate) Guild() (discord.Guild, bool) {
	if e.GuildID() != nil {
		return e.Client().Caches().Guilds().Get(*e.GuildID())
	}
	return discord.Guild{}, false
}

// Channel returns the discord.MessageChannel that the interaction happened in.
// This only returns cached channels.
func (e *ApplicationCommandInteractionCreate) Channel() (discord.MessageChannel, bool) {
	return e.Client().Caches().Channels().GetMessageChannel(e.ChannelID())
}

// DMChannel returns the discord.DMChannel that the interaction happened in.
// If the interaction happened in a guild, it returns nil.
// This only returns cached channels.
func (e *ApplicationCommandInteractionCreate) DMChannel() (discord.DMChannel, bool) {
	return e.Client().Caches().Channels().GetDMChannel(e.ChannelID())
}

// GuildChannel returns the discord.GuildMessageChannel that the interaction happened in.
// If the interaction happened in a dm, it returns nil.
// This only returns cached channels.
func (e *ApplicationCommandInteractionCreate) GuildChannel() (discord.GuildMessageChannel, bool) {
	return e.Client().Caches().Channels().GetGuildMessageChannel(e.ChannelID())
}

// ComponentInteractionCreate indicates that a new component interaction has been created.
type ComponentInteractionCreate struct {
	*GenericEvent
	discord.ComponentInteraction
	Respond InteractionResponderFunc
}

// Guild returns the guild that the interaction happened in if it happened in a guild.
// If the interaction happened in a DM, it returns nil.
// This only returns cached guilds.
func (e *ComponentInteractionCreate) Guild() (discord.Guild, bool) {
	if e.GuildID() != nil {
		return e.Client().Caches().Guilds().Get(*e.GuildID())
	}
	return discord.Guild{}, false
}

// Channel returns the discord.MessageChannel that the interaction happened in.
// This only returns cached channels.
func (e *ComponentInteractionCreate) Channel() (discord.MessageChannel, bool) {
	return e.Client().Caches().Channels().GetMessageChannel(e.ChannelID())
}

// DMChannel returns the discord.DMChannel that the interaction happened in.
// If the interaction happened in a guild, it returns nil.
// This only returns cached channels.
func (e *ComponentInteractionCreate) DMChannel() (discord.DMChannel, bool) {
	return e.Client().Caches().Channels().GetDMChannel(e.ChannelID())
}

// GuildChannel returns the discord.GuildMessageChannel that the interaction happened in.
// If the interaction happened in a dm, it returns nil.
// This only returns cached channels.
func (e *ComponentInteractionCreate) GuildChannel() (discord.GuildMessageChannel, bool) {
	return e.Client().Caches().Channels().GetGuildMessageChannel(e.ChannelID())
}

// AutocompleteInteractionCreate indicates that a new autocomplete interaction has been created.
type AutocompleteInteractionCreate struct {
	*GenericEvent
	discord.AutocompleteInteraction
	Respond InteractionResponderFunc
}

// Guild returns the guild that the interaction happened in if it happened in a guild.
// If the interaction happened in a DM, it returns nil.
// This only returns cached guilds.
func (e *AutocompleteInteractionCreate) Guild() (discord.Guild, bool) {
	if e.GuildID() != nil {
		return e.Client().Caches().Guilds().Get(*e.GuildID())
	}
	return discord.Guild{}, false
}

// Channel returns the discord.MessageChannel that the interaction happened in.
// This only returns cached channels.
func (e *AutocompleteInteractionCreate) Channel() (discord.MessageChannel, bool) {
	return e.Client().Caches().Channels().GetMessageChannel(e.ChannelID())
}

// DMChannel returns the discord.DMChannel that the interaction happened in.
// If the interaction happened in a guild, it returns nil.
// This only returns cached channels.
func (e *AutocompleteInteractionCreate) DMChannel() (discord.DMChannel, bool) {
	return e.Client().Caches().Channels().GetDMChannel(e.ChannelID())
}

// GuildChannel returns the discord.GuildMessageChannel that the interaction happened in.
// If the interaction happened in a dm, it returns nil.
// This only returns cached channels.
func (e *AutocompleteInteractionCreate) GuildChannel() (discord.GuildMessageChannel, bool) {
	return e.Client().Caches().Channels().GetGuildMessageChannel(e.ChannelID())
}

// ModalSubmitInteractionCreate indicates that a new modal submit interaction has been created.
type ModalSubmitInteractionCreate struct {
	*GenericEvent
	discord.ModalSubmitInteraction
	Respond InteractionResponderFunc
}

// Guild returns the guild that the interaction happened in if it happened in a guild.
// If the interaction happened in a DM, it returns nil.
// This only returns cached guilds.
func (e *ModalSubmitInteractionCreate) Guild() (discord.Guild, bool) {
	if e.GuildID() != nil {
		return e.Client().Caches().Guilds().Get(*e.GuildID())
	}
	return discord.Guild{}, false
}

// Channel returns the discord.MessageChannel that the interaction happened in.
// This only returns cached channels.
func (e *ModalSubmitInteractionCreate) Channel() (discord.MessageChannel, bool) {
	return e.Client().Caches().Channels().GetMessageChannel(e.ChannelID())
}

// DMChannel returns the discord.DMChannel that the interaction happened in.
// If the interaction happened in a guild, it returns nil.
// This only returns cached channels.
func (e *ModalSubmitInteractionCreate) DMChannel() (discord.DMChannel, bool) {
	return e.Client().Caches().Channels().GetDMChannel(e.ChannelID())
}

// GuildChannel returns the discord.GuildMessageChannel that the interaction happened in.
// If the interaction happened in a dm, it returns nil.
// This only returns cached channels.
func (e *ModalSubmitInteractionCreate) GuildChannel() (discord.GuildMessageChannel, bool) {
	return e.Client().Caches().Channels().GetGuildMessageChannel(e.ChannelID())
}
