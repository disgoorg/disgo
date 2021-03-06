package models

import "github.com/DiscoOrg/disgo/constants"

// Channel is a generic discord channel object
type Channel struct {
	ID   Snowflake             `json:"id"`
	Type constants.ChannelType `json:"type"`
}

/*
DMChannel is used for interacting in private messages with users
*/
type DMChannel struct {
	Channel
}

/*
GuildChannel is a generic type for all server channels
*/
type GuildChannel struct {
	Channel
}

/*
VoiceChannel adds methods specifically for interacting with discord's voice
*/
type VoiceChannel struct {
	GuildChannel
}

/*
TextChannel allows you to interact with discord's text channels
{
  "id": "41771983423143937",
  "guild_id": "41771983423143937",
  "name": "general",
  "type": 0,
  "position": 6,
  "permission_overwrites": [],
  "rate_limit_per_user": 2,
  "nsfw": true,
  "topic": "24/7 chat about how to gank Mike #2",
  "last_message_id": "155117677105512449",
  "parent_id": "399942396007890945"
}
*/
type TextChannel struct {
	GuildChannel
}

/*
StoreChannel allows you to interact with discord's store channels
*/
type StoreChannel struct {
	GuildChannel
}

/*
NewsChannel allows you to interact with discord's news channels
*/
type NewsChannel struct {
	TextChannel
}
