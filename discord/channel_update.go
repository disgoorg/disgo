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

func (GuildTextChannelUpdate) channelUpdate()      {}
func (GuildTextChannelUpdate) guildChannelUpdate() {}

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

func (GuildVoiceChannelUpdate) channelUpdate()      {}
func (GuildVoiceChannelUpdate) guildChannelUpdate() {}

type GroupDMChannelUpdate struct {
	Name *string   `json:"name,omitempty"`
	Icon *NullIcon `json:"icon,omitempty"`
}

func (GroupDMChannelUpdate) channelUpdate() {}

type GuildCategoryChannelUpdate struct {
}

func (GuildCategoryChannelUpdate) channelUpdate()      {}
func (GuildCategoryChannelUpdate) guildChannelUpdate() {}

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

func (GuildNewsChannelUpdate) channelUpdate()      {}
func (GuildNewsChannelUpdate) guildChannelUpdate() {}

type GuildStoreChannelUpdate struct {
}

func (GuildStoreChannelUpdate) channelUpdate()      {}
func (GuildStoreChannelUpdate) guildChannelUpdate() {}

type GuildNewsThreadUpdate struct {
}

func (GuildNewsThreadUpdate) channelUpdate()      {}
func (GuildNewsThreadUpdate) guildChannelUpdate() {}

type GuildPublicThreadUpdate struct {
}

func (GuildPublicThreadUpdate) channelUpdate()      {}
func (GuildPublicThreadUpdate) guildChannelUpdate() {}

type GuildPrivateThreadUpdate struct {
}

func (GuildPrivateThreadUpdate) channelUpdate()      {}
func (GuildPrivateThreadUpdate) guildChannelUpdate() {}

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

func (GuildStageVoiceChannelUpdate) channelUpdate()      {}
func (GuildStageVoiceChannelUpdate) guildChannelUpdate() {}
