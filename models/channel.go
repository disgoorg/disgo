package models

import "github.com/DiscoOrg/disgo/constants"

type Channel struct {
	ID Snowflake `json:"id"`
	Type constants.ChannelType `json:"type"`
}

/*

 */
type DMChannel struct {
	Channel
}

/*

 */
type GuildChannel struct {
	Channel
}

/*

*/
type VoiceChannel struct {
	GuildChannel
}

/*
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

 */
type StoreChannel struct {
	GuildChannel
}

/*

 */
type NewsChannel struct {
	TextChannel
}