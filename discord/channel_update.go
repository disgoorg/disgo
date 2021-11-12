package discord

type ChannelUpdate interface {
	channelUpdate()
}

type GuildChannelUpdate interface {
	ChannelUpdate
	guildChannelUpdate()
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
func (_ GuildTextChannelUpdate) guildChannelUpdate() {}

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
func (_ GuildVoiceChannelUpdate) guildChannelUpdate() {}

type GroupDMChannelUpdate struct {
	Name *string   `json:"name,omitempty"`
	Icon *NullIcon `json:"icon,omitempty"`
}

func (_ GroupDMChannelUpdate) channelUpdate() {}

type GuildCategoryChannelUpdate struct {
}

func (_ GuildCategoryChannelUpdate) channelUpdate() {}
func (_ GuildCategoryChannelUpdate) guildChannelUpdate() {}


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
func (_ GuildNewsChannelUpdate) guildChannelUpdate() {}

type GuildStoreChannelUpdate struct {
}

func (_ GuildStoreChannelUpdate) channelUpdate() {}
func (_ GuildStoreChannelUpdate) guildChannelUpdate() {}

type GuildNewsThreadUpdate struct {
}

func (_ GuildNewsThreadUpdate) channelUpdate() {}
func (_ GuildNewsThreadUpdate) guildChannelUpdate() {}

type GuildPublicThreadUpdate struct {
}

func (_ GuildPublicThreadUpdate) channelUpdate() {}
func (_ GuildPublicThreadUpdate) guildChannelUpdate() {}

type GuildPrivateThreadUpdate struct {
}

func (_ GuildPrivateThreadUpdate) channelUpdate() {}
func (_ GuildPrivateThreadUpdate) guildChannelUpdate() {}

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
func (_ GuildStageVoiceChannelUpdate) guildChannelUpdate() {}
