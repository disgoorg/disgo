package discord

type ChannelUpdate interface {
	Type() ChannelType
	channelUpdate()
}

type GuildChannelUpdate interface {
	ChannelUpdate
	guildChannelCreate()
}

type GuildTextChannelUpdate struct {
	Name                       *string                `json:"name,omitempty"`
	Type                       *ChannelType           `json:"type,omitempty"`
	Position                   *int                   `json:"position,omitempty"`
	Topic                      *string                `json:"topic,omitempty"`
	NSFW                       *bool                  `json:"nsfw,omitempty"`
	RateLimitPerUser           *int                   `json:"rate_limit_per_user,omitempty"`
	PermissionOverwrites       *[]PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID                   *Snowflake             `json:"parent_id,omitempty"`
	DefaultAutoArchiveDuration *int                   `json:"default_auto_archive_duration,omitempty"`
}

func (_ GuildTextChannelUpdate) channelUpdate() {}

type GuildVoiceChannelUpdate struct {
	Name                       *string                `json:"name,omitempty"`
	Position                   *int                   `json:"position,omitempty"`
	RateLimitPerUser           *int                   `json:"rate_limit_per_user,omitempty"`
	Bitrate                    *int                   `json:"bitrate,omitempty"`
	UserLimit                  *int                   `json:"user_limit,omitempty"`
	PermissionOverwrites       *[]PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID                   *Snowflake             `json:"parent_id,omitempty"`
	RTCRegion                  *string                `json:"rtc_region"`
	VideoQualityMode           *VideoQualityMode      `json:"video_quality_mode"`
	DefaultAutoArchiveDuration *int                   `json:"default_auto_archive_duration"`
}

func (_ GuildVoiceChannelUpdate) channelUpdate() {}

type GroupDMChannelUpdate struct {
	Name *string       `json:"name,omitempty"`
	Icon *OptionalIcon `json:"icon,omitempty"`
}

func (_ GroupDMChannelUpdate) channelUpdate() {}

type GuildCategoryChannelUpdate struct {
}

func (_ GuildCategoryChannelUpdate) channelUpdate() {}

type GuildNewsChannelUpdate struct {
	Name                       *string                `json:"name,omitempty"`
	Type                       *ChannelType           `json:"type,omitempty"`
	Position                   *int                   `json:"position,omitempty"`
	Topic                      *string                `json:"topic,omitempty"`
	RateLimitPerUser           *int                   `json:"rate_limit_per_user,omitempty"`
	PermissionOverwrites       *[]PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID                   *Snowflake             `json:"parent_id,omitempty"`
	DefaultAutoArchiveDuration *int                   `json:"default_auto_archive_duration"`
}

func (_ GuildNewsChannelUpdate) channelUpdate() {}

type GuildStoreChannelUpdate struct {
}

func (_ GuildStoreChannelUpdate) channelUpdate() {}

type GuildNewsThreadUpdate struct {
}

func (_ GuildNewsThreadUpdate) channelUpdate() {}

type GuildPublicThreadUpdate struct {
}

func (_ GuildPublicThreadUpdate) channelUpdate() {}

type GuildPrivateThreadUpdate struct {
}

func (_ GuildPrivateThreadUpdate) channelUpdate() {}

type GuildStageVoiceChannelUpdate struct {
	Name                 *string                `json:"name,omitempty"`
	Position             *int                   `json:"position,omitempty"`
	Topic                *string                `json:"topic,omitempty"`
	Bitrate              *int                   `json:"bitrate,omitempty"`
	UserLimit            *int                   `json:"user_limit,omitempty"`
	PermissionOverwrites *[]PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID             *Snowflake             `json:"parent_id,omitempty"`
	RTCRegion            *string                `json:"rtc_region"`
	VideoQualityMode     *VideoQualityMode      `json:"video_quality_mode"`
}

func (_ GuildStageVoiceChannelUpdate) channelUpdate() {}

func (_ GuildStageVoiceChannelUpdate) Type() ChannelType {
	return ChannelTypeGuildStageVoice
}
