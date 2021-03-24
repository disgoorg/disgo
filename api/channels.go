package api

import (
	//"time"
)

// ChannelType for interacting with discord's channels
type ChannelType int

// Channel constants
const (
	ChannelTypeText ChannelType = iota
	ChannelTypeDM
	ChannelTypeVoice
	ChannelTypeGroupDM
	ChannelTypeCategory
	ChannelTypeNews
	ChannelTypeStore
)

// Channel is a generic discord channel object
type Channel struct {
	Disgo            Disgo
	ID               Snowflake   `json:"id"`
	Type             ChannelType `json:"type"`
	LastMessageID    *Snowflake  `json:"last_message_id,omitempty"`
	Name             *string     `json:"name,omitempty"`
	GuildID          *Snowflake  `json:"guild_id,omitempty"`
	Position         *int        `json:"position,omitempty"`
	Topic            *string     `json:"topic,omitempty"`
	NSFW             *bool       `json:"nsfw,omitempty"`
	Bitrate          *int        `json:"bitrate,omitempty"`
	UserLimit        *int        `json:"user_limit,omitempty"`
	RateLimitPerUser *int        `json:"rate_limit_per_user,omitempty"`
	Recipients       []*User     `json:"recipients,omitempty"`
	Icon             *string     `json:"icon,omitempty"`
	OwnerID          *Snowflake  `json:"owner_id,omitempty"`
	ApplicationID    *Snowflake  `json:"application_id,omitempty"`
	ParentID         *Snowflake  `json:"parent_id,omitempty"`
	//LastPinTimestamp *time.Time  `json:"last_pin_timestamp,omitempty"`
}

// MessageChannel is used for sending messages to user
type MessageChannel struct {
	Channel
}

// SendMessage a Message to a TextChannel
func (c MessageChannel) SendMessage(content string) (*Message, error) {
	// Todo: embeds, attachments etc.
	return c.Disgo.RestClient().SendMessage(c.ID, Message{Content: &content})
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
}

// Guild returns the channel's Guild
func (c GuildChannel) Guild() *Guild {
	return c.Disgo.Cache().Guild(c.GuildID)
}

// CategoryChannel groups text & voice channels in servers together
type CategoryChannel struct {
	GuildChannel
}

// VoiceChannel adds methods specifically for interacting with discord's voice
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
