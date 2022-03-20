package discord

import (
	"github.com/DisgoOrg/snowflake"
)

type dmChannel struct {
	ID               snowflake.Snowflake  `json:"id"`
	Type             ChannelType          `json:"type"`
	LastMessageID    *snowflake.Snowflake `json:"last_message_id"`
	Recipients       []User               `json:"recipients"`
	LastPinTimestamp *Time                `json:"last_pin_timestamp"`
}

type guildTextChannel struct {
	ID                         snowflake.Snowflake   `json:"id"`
	Type                       ChannelType           `json:"type"`
	GuildID                    snowflake.Snowflake   `json:"guild_id"`
	Position                   int                   `json:"position"`
	PermissionOverwrites       []PermissionOverwrite `json:"permission_overwrites"`
	Name                       string                `json:"name"`
	Topic                      *string               `json:"topic"`
	NSFW                       bool                  `json:"nsfw"`
	LastMessageID              *snowflake.Snowflake  `json:"last_message_id"`
	RateLimitPerUser           int                   `json:"rate_limit_per_user"`
	ParentID                   *snowflake.Snowflake  `json:"parent_id"`
	LastPinTimestamp           *Time                 `json:"last_pin_timestamp"`
	DefaultAutoArchiveDuration AutoArchiveDuration   `json:"default_auto_archive_duration"`
}

type guildNewsChannel struct {
	ID                         snowflake.Snowflake   `json:"id"`
	Type                       ChannelType           `json:"type"`
	GuildID                    snowflake.Snowflake   `json:"guild_id"`
	Position                   int                   `json:"position"`
	PermissionOverwrites       []PermissionOverwrite `json:"permission_overwrites"`
	Name                       string                `json:"name"`
	Topic                      *string               `json:"topic"`
	NSFW                       bool                  `json:"nsfw"`
	RateLimitPerUser           int                   `json:"rate_limit_per_user"`
	ParentID                   *snowflake.Snowflake  `json:"parent_id"`
	LastMessageID              *snowflake.Snowflake  `json:"last_message_id"`
	LastPinTimestamp           *Time                 `json:"last_pin_timestamp"`
	DefaultAutoArchiveDuration AutoArchiveDuration   `json:"default_auto_archive_duration"`
}

type guildThread struct {
	ID               snowflake.Snowflake  `json:"id"`
	Type             ChannelType          `json:"type"`
	GuildID          snowflake.Snowflake  `json:"guild_id"`
	Name             string               `json:"name"`
	NSFW             bool                 `json:"nsfw"`
	LastMessageID    *snowflake.Snowflake `json:"last_message_id"`
	RateLimitPerUser int                  `json:"rate_limit_per_user"`
	OwnerID          snowflake.Snowflake  `json:"owner_id"`
	ParentID         snowflake.Snowflake  `json:"parent_id"`
	LastPinTimestamp *Time                `json:"last_pin_timestamp"`
	MessageCount     int                  `json:"message_count"`
	MemberCount      int                  `json:"member_count"`
	ThreadMetadata   ThreadMetadata       `json:"thread_metadata"`
}

type guildCategoryChannel struct {
	ID                   snowflake.Snowflake   `json:"id"`
	Type                 ChannelType           `json:"type"`
	GuildID              snowflake.Snowflake   `json:"guild_id"`
	Position             int                   `json:"position"`
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites"`
	Name                 string                `json:"name"`
}

type guildVoiceChannel struct {
	ID                   snowflake.Snowflake   `json:"id"`
	Type                 ChannelType           `json:"type"`
	GuildID              snowflake.Snowflake   `json:"guild_id"`
	Position             int                   `json:"position"`
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites"`
	Name                 string                `json:"name"`
	Bitrate              int                   `json:"bitrate"`
	UserLimit            int                   `json:"user_limit"`
	ParentID             *snowflake.Snowflake  `json:"parent_id"`
	RTCRegion            string                `json:"rtc_region"`
	VideoQualityMode     VideoQualityMode      `json:"video_quality_mode"`
}

type guildStageVoiceChannel struct {
	ID                   snowflake.Snowflake   `json:"id"`
	Type                 ChannelType           `json:"type"`
	GuildID              snowflake.Snowflake   `json:"guild_id"`
	Position             int                   `json:"position"`
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites"`
	Name                 string                `json:"name"`
	Bitrate              int                   `json:"bitrate,"`
	ParentID             *snowflake.Snowflake  `json:"parent_id"`
	RTCRegion            string                `json:"rtc_region"`
}
