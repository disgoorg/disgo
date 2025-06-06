package discord

import (
	"github.com/disgoorg/omit"
	"github.com/disgoorg/snowflake/v2"
)

type ChannelUpdate interface {
	channelUpdate()
}

type GuildChannelUpdate interface {
	ChannelUpdate
	guildChannelUpdate()
}

type GuildTextChannelUpdate struct {
	Name                          *string                `json:"name,omitempty"`
	Type                          *ChannelType           `json:"type,omitempty"`
	Position                      *int                   `json:"position,omitempty"`
	Topic                         *string                `json:"topic,omitempty"`
	NSFW                          *bool                  `json:"nsfw,omitempty"`
	RateLimitPerUser              *int                   `json:"rate_limit_per_user,omitempty"`
	PermissionOverwrites          *[]PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID                      *snowflake.ID          `json:"parent_id,omitempty"`
	DefaultAutoArchiveDuration    *AutoArchiveDuration   `json:"default_auto_archive_duration,omitempty"`
	DefaultThreadRateLimitPerUser *int                   `json:"default_thread_rate_limit_per_user,omitempty"`
}

func (GuildTextChannelUpdate) channelUpdate()      {}
func (GuildTextChannelUpdate) guildChannelUpdate() {}

type GuildVoiceChannelUpdate struct {
	Name                 *string                `json:"name,omitempty"`
	Position             *int                   `json:"position,omitempty"`
	RateLimitPerUser     *int                   `json:"rate_limit_per_user,omitempty"`
	Bitrate              *int                   `json:"bitrate,omitempty"`
	UserLimit            *int                   `json:"user_limit,omitempty"`
	PermissionOverwrites *[]PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID             *snowflake.ID          `json:"parent_id,omitempty"`
	RTCRegion            *string                `json:"rtc_region,omitempty"`
	NSFW                 *bool                  `json:"nsfw,omitempty"`
	VideoQualityMode     *VideoQualityMode      `json:"video_quality_mode,omitempty"`
}

func (GuildVoiceChannelUpdate) channelUpdate()      {}
func (GuildVoiceChannelUpdate) guildChannelUpdate() {}

type GuildCategoryChannelUpdate struct {
	Name                 *string                `json:"name,omitempty"`
	Position             *int                   `json:"position,omitempty"`
	PermissionOverwrites *[]PermissionOverwrite `json:"permission_overwrites,omitempty"`
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
	ParentID                   *snowflake.ID          `json:"parent_id,omitempty"`
	DefaultAutoArchiveDuration *int                   `json:"default_auto_archive_duration,omitempty"`
}

func (GuildNewsChannelUpdate) channelUpdate()      {}
func (GuildNewsChannelUpdate) guildChannelUpdate() {}

type GuildThreadUpdate struct {
	Name                *string              `json:"name,omitempty"`
	Archived            *bool                `json:"archived,omitempty"`
	AutoArchiveDuration *AutoArchiveDuration `json:"auto_archive_duration,omitempty"`
	Locked              *bool                `json:"locked,omitempty"`
	Invitable           *bool                `json:"invitable,omitempty"`
	RateLimitPerUser    *int                 `json:"rate_limit_per_user,omitempty"`
}

func (GuildThreadUpdate) channelUpdate()      {}
func (GuildThreadUpdate) guildChannelUpdate() {}

type GuildStageVoiceChannelUpdate struct {
	Name                 *string                `json:"name,omitempty"`
	Position             *int                   `json:"position,omitempty"`
	RateLimitPerUser     *int                   `json:"rate_limit_per_user,omitempty"`
	Bitrate              *int                   `json:"bitrate,omitempty"`
	UserLimit            *int                   `json:"user_limit,omitempty"`
	PermissionOverwrites *[]PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID             *snowflake.ID          `json:"parent_id,omitempty"`
	RTCRegion            *string                `json:"rtc_region,omitempty"`
	NSFW                 *bool                  `json:"nsfw,omitempty"`
	VideoQualityMode     *VideoQualityMode      `json:"video_quality_mode,omitempty"`
}

func (GuildStageVoiceChannelUpdate) channelUpdate()      {}
func (GuildStageVoiceChannelUpdate) guildChannelUpdate() {}

type GuildForumChannelUpdate struct {
	Name                          *string                          `json:"name,omitempty"`
	Position                      *int                             `json:"position,omitempty"`
	Topic                         *string                          `json:"topic,omitempty"`
	NSFW                          *bool                            `json:"nsfw,omitempty"`
	PermissionOverwrites          *[]PermissionOverwrite           `json:"permission_overwrites,omitempty"`
	ParentID                      *snowflake.ID                    `json:"parent_id,omitempty"`
	RateLimitPerUser              *int                             `json:"rate_limit_per_user"`
	AvailableTags                 *[]ChannelTag                    `json:"available_tags,omitempty"`
	Flags                         *ChannelFlags                    `json:"flags,omitempty"`
	DefaultReactionEmoji          omit.Omit[*DefaultReactionEmoji] `json:"default_reaction_emoji,omitzero"`
	DefaultThreadRateLimitPerUser *int                             `json:"default_thread_rate_limit_per_user,omitempty"`
	DefaultSortOrder              omit.Omit[*DefaultSortOrder]     `json:"default_sort_order,omitzero"`
	DefaultForumLayout            omit.Omit[*DefaultForumLayout]   `json:"default_forum_layout,omitzero"`
}

func (GuildForumChannelUpdate) channelUpdate()      {}
func (GuildForumChannelUpdate) guildChannelUpdate() {}

type GuildMediaChannelUpdate struct {
	Name                          *string                          `json:"name,omitempty"`
	Position                      *int                             `json:"position,omitempty"`
	Topic                         *string                          `json:"topic,omitempty"`
	NSFW                          *bool                            `json:"nsfw,omitempty"`
	PermissionOverwrites          *[]PermissionOverwrite           `json:"permission_overwrites,omitempty"`
	ParentID                      *snowflake.ID                    `json:"parent_id,omitempty"`
	RateLimitPerUser              *int                             `json:"rate_limit_per_user"`
	AvailableTags                 *[]ChannelTag                    `json:"available_tags,omitempty"`
	Flags                         *ChannelFlags                    `json:"flags,omitempty"`
	DefaultReactionEmoji          omit.Omit[*DefaultReactionEmoji] `json:"default_reaction_emoji,omitzero"`
	DefaultThreadRateLimitPerUser *int                             `json:"default_thread_rate_limit_per_user,omitempty"`
	DefaultSortOrder              omit.Omit[*DefaultSortOrder]     `json:"default_sort_order,omitzero"`
}

func (GuildMediaChannelUpdate) channelUpdate()      {}
func (GuildMediaChannelUpdate) guildChannelUpdate() {}

type GuildPostUpdate struct {
	Name                *string              `json:"name,omitempty"`
	Archived            *bool                `json:"archived,omitempty"`
	AutoArchiveDuration *AutoArchiveDuration `json:"auto_archive_duration,omitempty"`
	Locked              *bool                `json:"locked,omitempty"`
	Invitable           *bool                `json:"invitable,omitempty"`
	RateLimitPerUser    *int                 `json:"rate_limit_per_user,omitempty"`
	Flags               *ChannelFlags        `json:"flags,omitempty"`
	AppliedTags         *[]snowflake.ID      `json:"applied_tags,omitempty"`
}

func (GuildPostUpdate) channelUpdate()      {}
func (GuildPostUpdate) guildChannelUpdate() {}

type GuildChannelPositionUpdate struct {
	ID              snowflake.ID     `json:"id"`
	Position        omit.Omit[*int]  `json:"position,omitzero"`
	LockPermissions omit.Omit[*bool] `json:"lock_permissions,omitzero"`
	ParentID        *snowflake.ID    `json:"parent_id,omitempty"`
}
