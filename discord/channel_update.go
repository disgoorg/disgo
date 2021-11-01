package discord

import "github.com/DisgoOrg/disgo/json"

type ChannelUpdate interface {
	json.Marshaler
	Type() ChannelType
}

type GuildChannelUpdate interface {
	ChannelUpdate
	guildChannelCreate()
}

type GuildTextChannelUpdate struct {
	Name                       *string                `json:"name,omitempty"`
	Position                   *int                   `json:"position,omitempty"`
	Topic                      *string                `json:"topic,omitempty"`
	NSFW                       *bool                  `json:"nsfw,omitempty"`
	RateLimitPerUser           *int                   `json:"rate_limit_per_user,omitempty"`
	PermissionOverwrites       *[]PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID                   *Snowflake             `json:"parent_id,omitempty"`
	DefaultAutoArchiveDuration *int                   `json:"default_auto_archive_duration"`
}

type DMChannelUpdate struct {
	Name *string `json:"name,omitempty"`
	Icon *OptionalIcon   `json:"icon,omitempty"`
}

type GuildVoiceChannelUpdate struct {
	Name                       *string                `json:"name,omitempty"`
	Type                       *ChannelType           `json:"type,omitempty"`
	Position                   *int                   `json:"position,omitempty"`
	Topic                      *string                `json:"topic,omitempty"`
	RateLimitPerUser           *int                   `json:"rate_limit_per_user,omitempty"`
	Bitrate                    *int                   `json:"bitrate,omitempty"`
	UserLimit                  *int                   `json:"user_limit,omitempty"`
	PermissionOverwrites       *[]PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID                   *Snowflake             `json:"parent_id,omitempty"`
	RTCRegion                  *string                `json:"rtc_region"`
	VideoQualityMode           *VideoQualityMode      `json:"video_quality_mode"`
	DefaultAutoArchiveDuration *int                   `json:"default_auto_archive_duration"`
}

type GuildCategoryChannelUpdate struct {

}

type GuildNewsChannelUpdate struct {
}

type GuildStoreChannelUpdate struct {
}

type GuildNewsThreadUpdate struct {
}

type GuildPublicThreadUpdate struct {
}

type GuildPrivateThreadUpdate struct {
}

type GuildStageChannelUpdate struct {
}
