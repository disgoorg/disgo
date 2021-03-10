package api

import (
	"github.com/chebyrash/promise"
)

// ChannelType for interacting with discord's channels
type ChannelType int

// Channel constants
const (
	GuildTextChannel ChannelType = iota
	DMTextChannel
	GuildVoiceChannel
	GroupDMChannel
	GuildCategoryChannel
	GuildNewsChannel
	GuildStoreChannel
)

// Channel is a generic discord channel object
type Channel struct {
	Disgo         Disgo
	ID            Snowflake             `json:"id"`
	Type          ChannelType `json:"type"`
	LastMessageID Snowflake             `json:"last_message_id"`
}

//  MessageChannel is used for sending messages to user
type MessageChannel struct {
	Channel
}

func (c MessageChannel) SendMessage(content string) *promise.Promise {
	return c.Disgo.RestClient().SendMessage(c.ID, Message{Content: content})
}

// DMChannel is used for interacting in private messages with users
type DMChannel struct {
	MessageChannel
	Users []User `json:"recipients"`
}

// GuildChannel is a generic type for all server channels
type GuildChannel struct {
	Channel
	GuildID Snowflake `json:"guild_id"`
	Guild   Guild
}

// CategoryChannel groups text & voice channels in servers together
type CategoryChannel struct {
	GuildChannel
}

//  VoiceChannel adds methods specifically for interacting with discord's voice
type VoiceChannel struct {
	GuildChannel
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

// NewsChannel allows you to interact with discord's news channels
type NewsChannel struct {
	TextChannel
}
