package entities

import (
	"errors"
	"fmt"

	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type Channel struct {
	discord.Channel
}

// Mention returns the channel as a string
func (c Channel) Mention() string {
	return fmt.Sprintf("<#%s>", c.ID)
}

// MessageChannel is used for sending Message(s) to User(s)
type MessageChannel struct {
	discord.Channel
	Disgo core.Disgo
}

// SendMessage sends a Message to a TextChannel
func (c MessageChannel) SendMessage(message discord.MessageCreate) (*discord.Message, rest.Error) {
	return c.Disgo.RestServices().ChannelsService().CreateMessage(c.ID, message)
}

// EditMessage edits a Message in this TextChannel
func (c MessageChannel) EditMessage(messageID discord.Snowflake, message discord.MessageUpdate) (*discord.Message, rest.Error) {
	return c.Disgo.RestServices().ChannelsService().UpdateMessage(c.ID, messageID, message)
}

// DeleteMessage allows you to edit an existing Message sent by you
func (c MessageChannel) DeleteMessage(messageID discord.Snowflake) rest.Error {
	return c.Disgo.RestServices().ChannelsService().DeleteMessage(c.ID, messageID)
}

// BulkDeleteMessages allows you bulk delete Message(s)
func (c MessageChannel) BulkDeleteMessages(messageIDs ...discord.Snowflake) rest.Error {
	return c.Disgo.RestServices().ChannelsService().BulkDeleteMessages(c.ID, messageIDs...)
}

// CrosspostMessage crossposts an existing Message
func (c MessageChannel) CrosspostMessage(messageID discord.Snowflake) (*discord.Message, rest.Error) {
	if c.Type != discord.ChannelTypeNews {
		return nil, rest.NewError(nil, errors.New("channel type is not NEWS"))
	}
	return c.Disgo.RestServices().ChannelsService().CrosspostMessage(c.ID, messageID)
}

// DMChannel is used for interacting in private Message(s) with users
type DMChannel struct {
	MessageChannel
}

// CreateDMChannel is the payload used to create a DMChannel
type CreateDMChannel struct {
	RecipientID discord.Snowflake `json:"recipient_id"`
}

// GuildChannel is a generic type for all server channels
type GuildChannel struct {
	Channel
	Disgo core.Disgo
}

// Guild returns the channel's Guild
func (c GuildChannel) Guild() *core.Guild {
	if c.GuildID == nil {
		return nil
	}
	return c.Disgo.Cache().GuildCache().Get(*c.GuildID)
}

// Category groups text & voice channels in servers together
type Category struct {
	GuildChannel
}

// VoiceChannel adds methods specifically for interacting with discord's voice
type VoiceChannel struct {
	GuildChannel
}

// Connect sends an GatewayCommand to connect to this VoiceChannel
func (c *VoiceChannel) Connect() error {
	return c.Disgo.AudioController().Connect(*c.GuildID, c.ID)
}

// TextChannel allows you to interact with discord's text channels
type TextChannel struct {
	GuildChannel
	MessageChannel
}

// StoreChannel allows you to interact with discord's store channels
type StoreChannel struct {
	GuildChannel
}
