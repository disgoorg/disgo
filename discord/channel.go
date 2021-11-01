package discord

import (
	"fmt"

	"github.com/DisgoOrg/disgo/json"
	"github.com/DisgoOrg/disgo/rest/route"
)

// ChannelType for interacting with discord's channels
type ChannelType int

// Channel constants
//goland:noinspection GoUnusedConst
const (
	ChannelTypeGuildText ChannelType = iota
	ChannelTypeDM
	ChannelTypeGuildVoice
	ChannelTypeGroupDM
	ChannelTypeGuildCategory
	ChannelTypeGuildNews
	ChannelTypeGuildStore
	_
	_
	_
	ChannelTypeGuildNewsThread
	ChannelTypeGuildPublicThread
	ChannelTypeGuildPrivateThread
	ChannelTypeGuildStageVoice
)

var channels = map[ChannelType]func() Channel{
	ChannelTypeGuildText: func() Channel {
		return &GuildTextChannel{}
	},
	ChannelTypeDM: func() Channel {
		return &DMChannel{}
	},
	ChannelTypeGuildVoice: func() Channel {
		return &GuildVoiceChannel{}
	},
	ChannelTypeGroupDM: func() Channel {
		return &GroupDMChannel{}
	},
	ChannelTypeGuildCategory: func() Channel {
		return &GuildCategoryChannel{}
	},
	ChannelTypeGuildNews: func() Channel {
		return &GuildNewsChannel{}
	},
	ChannelTypeGuildStore: func() Channel {
		return &GuildStoreChannel{}
	},
	ChannelTypeGuildNewsThread: func() Channel {
		return &GuildNewsThread{}
	},
	ChannelTypeGuildPublicThread: func() Channel {
		return &GuildPublicThread{}
	},
	ChannelTypeGuildPrivateThread: func() Channel {
		return &GuildPrivateThread{}
	},
	ChannelTypeGuildStageVoice: func() Channel {
		return &GuildStageVoiceChannel{}
	},
}

type Channel interface {
	json.Marshaler
	Type() ChannelType
}

type GuildChannel interface {
	Channel
	guildChannel()
}

type MessageChannel interface {
	Channel
	messageChannel()
}

type GuildMessageChannel interface {
	GuildChannel
	MessageChannel
}

type GuildThread interface {
	GuildMessageChannel
	guildThread()
}

type AudioChannel interface {
	Channel
	audioChannel()
}

type UnmarshalChannel struct {
	Channel
}

func (u *UnmarshalChannel) UnmarshalJSON(data []byte) error {
	var cType struct {
		Type ChannelType `json:"type"`
	}

	if err := json.Unmarshal(data, &cType); err != nil {
		return err
	}

	fn, ok := channels[cType.Type]
	if !ok {
		return fmt.Errorf("unkown channel with type %d received", cType.Type)
	}

	v := fn()

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	u.Channel = v
	return nil
}

var (
	_ MessageChannel      = (*GuildTextChannel)(nil)
	_ GuildMessageChannel = (*GuildTextChannel)(nil)
	_ GuildChannel        = (*GuildTextChannel)(nil)
)

type GuildTextChannel struct {
	ID                         Snowflake             `json:"id"`
	GuildID                    Snowflake             `json:"guild_id,omitempty"`
	Position                   int                   `json:"position,omitempty"`
	PermissionOverwrites       []PermissionOverwrite `json:"permission_overwrites"`
	Name                       string                `json:"name,omitempty"`
	Topic                      *string               `json:"topic,omitempty"`
	NSFW                       bool                  `json:"nsfw,omitempty"`
	LastMessageID              *Snowflake            `json:"last_message_id,omitempty"`
	RateLimitPerUser           int                   `json:"rate_limit_per_user,omitempty"`
	ParentID                   *Snowflake            `json:"parent_id,omitempty"`
	LastPinTimestamp           *Time                 `json:"last_pin_timestamp,omitempty"`
	DefaultAutoArchiveDuration AutoArchiveDuration   `json:"default_auto_archive_duration"`
	InteractionPermissions     Permissions           `json:"permissions,omitempty"`
}

func (c GuildTextChannel) MarshalJSON() ([]byte, error) {
	type guildTextChannel GuildTextChannel
	return json.Marshal(struct {
		Type ChannelType `json:"type"`
		guildTextChannel
	}{
		Type:             c.Type(),
		guildTextChannel: guildTextChannel(c),
	})
}

func (_ *GuildTextChannel) Type() ChannelType {
	return ChannelTypeGuildText
}
func (_ *GuildTextChannel) messageChannel()      {}
func (_ *GuildTextChannel) guildMessageChannel() {}
func (_ *GuildTextChannel) guildChannel()        {}

var (
	_ MessageChannel = (*DMChannel)(nil)
)

type DMChannel struct {
	ID               Snowflake  `json:"id"`
	Name             string     `json:"name,omitempty"`
	LastMessageID    *Snowflake `json:"last_message_id,omitempty"`
	Recipients       []User     `json:"recipients,omitempty"`
	LastPinTimestamp *Time      `json:"last_pin_timestamp,omitempty"`
}

func (c DMChannel) MarshalJSON() ([]byte, error) {
	type dmChannel DMChannel
	return json.Marshal(struct {
		Type ChannelType `json:"type"`
		dmChannel
	}{
		Type:      c.Type(),
		dmChannel: dmChannel(c),
	})
}

func (_ *DMChannel) Type() ChannelType {
	return ChannelTypeGuildText
}
func (_ *DMChannel) messageChannel() {}

var (
	_ GuildChannel = (*GuildVoiceChannel)(nil)
	_ AudioChannel = (*GuildVoiceChannel)(nil)
)

type GuildVoiceChannel struct {
	ID                     Snowflake             `json:"id"`
	GuildID                Snowflake             `json:"guild_id,omitempty"`
	Position               int                   `json:"position,omitempty"`
	PermissionOverwrites   []PermissionOverwrite `json:"permission_overwrites"`
	Name                   string                `json:"name,omitempty"`
	NSFW                   bool                  `json:"nsfw,omitempty"`
	Topic                  *string               `json:"topic,omitempty"`
	Bitrate                int                   `json:"bitrate,omitempty"`
	UserLimit              int                   `json:"user_limit,omitempty"`
	ParentID               *Snowflake            `json:"parent_id,omitempty"`
	RTCRegion              string                `json:"rtc_region"`
	VideoQualityMode       VideoQualityMode      `json:"video_quality_mode"`
	InteractionPermissions Permissions           `json:"permissions,omitempty"`
}

func (c GuildVoiceChannel) MarshalJSON() ([]byte, error) {
	type guildVoiceChannel GuildVoiceChannel
	return json.Marshal(struct {
		Type ChannelType `json:"type"`
		guildVoiceChannel
	}{
		Type:              c.Type(),
		guildVoiceChannel: guildVoiceChannel(c),
	})
}

func (_ *GuildVoiceChannel) Type() ChannelType {
	return ChannelTypeGuildText
}
func (_ *GuildVoiceChannel) guildChannel() {}
func (_ *GuildVoiceChannel) audioChannel() {}

var (
	_ MessageChannel = (*GroupDMChannel)(nil)
)

type GroupDMChannel struct {
	ID               Snowflake  `json:"id"`
	Name             string     `json:"name,omitempty"`
	LastMessageID    *Snowflake `json:"last_message_id,omitempty"`
	Recipients       []User     `json:"recipients,omitempty"`
	Icon             *string    `json:"icon,omitempty"`
	OwnerID          Snowflake  `json:"owner_id,omitempty"`
	ApplicationID    Snowflake  `json:"application_id,omitempty"`
	LastPinTimestamp *Time      `json:"last_pin_timestamp,omitempty"`
}

func (c GroupDMChannel) MarshalJSON() ([]byte, error) {
	type groupDMChannel GroupDMChannel
	return json.Marshal(struct {
		Type ChannelType `json:"type"`
		groupDMChannel
	}{
		Type:           c.Type(),
		groupDMChannel: groupDMChannel(c),
	})
}

func (_ *GroupDMChannel) Type() ChannelType {
	return ChannelTypeGuildText
}
func (_ *GroupDMChannel) messageChannel() {}

var (
	_ GuildChannel = (*GuildCategoryChannel)(nil)
)

type GuildCategoryChannel struct {
	ID                     Snowflake             `json:"id"`
	GuildID                Snowflake             `json:"guild_id"`
	Position               int                   `json:"position"`
	PermissionOverwrites   []PermissionOverwrite `json:"permission_overwrites"`
	Name                   string                `json:"name"`
	NSFW                   bool                  `json:"nsfw"`
	InteractionPermissions Permissions           `json:"permissions,omitempty"`
}

func (c GuildCategoryChannel) MarshalJSON() ([]byte, error) {
	type guildCategoryChannel GuildCategoryChannel
	return json.Marshal(struct {
		Type ChannelType `json:"type"`
		guildCategoryChannel
	}{
		Type:                 c.Type(),
		guildCategoryChannel: guildCategoryChannel(c),
	})
}

func (_ *GuildCategoryChannel) Type() ChannelType {
	return ChannelTypeGuildCategory
}
func (_ *GuildCategoryChannel) guildChannel() {}

var (
	_ GuildChannel        = (*GuildNewsChannel)(nil)
	_ GuildMessageChannel = (*GuildNewsChannel)(nil)
)

type GuildNewsChannel struct {
	ID                         Snowflake             `json:"id"`
	GuildID                    Snowflake             `json:"guild_id,omitempty"`
	Position                   int                   `json:"position,omitempty"`
	PermissionOverwrites       []PermissionOverwrite `json:"permission_overwrites"`
	Name                       string                `json:"name,omitempty"`
	Topic                      *string               `json:"topic,omitempty"`
	NSFW                       bool                  `json:"nsfw,omitempty"`
	LastMessageID              *Snowflake            `json:"last_message_id,omitempty"`
	RateLimitPerUser           int                   `json:"rate_limit_per_user,omitempty"`
	ParentID                   *Snowflake            `json:"parent_id,omitempty"`
	LastPinTimestamp           *Time                 `json:"last_pin_timestamp,omitempty"`
	DefaultAutoArchiveDuration AutoArchiveDuration   `json:"default_auto_archive_duration"`
	InteractionPermissions     Permissions           `json:"permissions,omitempty"`
}

func (c GuildNewsChannel) MarshalJSON() ([]byte, error) {
	type guildNewsChannel GuildNewsChannel
	return json.Marshal(struct {
		Type ChannelType `json:"type"`
		guildNewsChannel
	}{
		Type:             c.Type(),
		guildNewsChannel: guildNewsChannel(c),
	})
}

func (_ *GuildNewsChannel) Type() ChannelType {
	return ChannelTypeGuildNews
}
func (_ *GuildNewsChannel) guildChannel()   {}
func (_ *GuildNewsChannel) messageChannel() {}

type GuildStoreChannel struct {
	ID                     Snowflake             `json:"id"`
	GuildID                Snowflake             `json:"guild_id"`
	Position               int                   `json:"position"`
	PermissionOverwrites   []PermissionOverwrite `json:"permission_overwrites"`
	Name                   string                `json:"name"`
	NSFW                   bool                  `json:"nsfw,omitempty"`
	InteractionPermissions Permissions           `json:"permissions,omitempty"`
}

func (c GuildStoreChannel) MarshalJSON() ([]byte, error) {
	type guildStoreChannel GuildStoreChannel
	return json.Marshal(struct {
		Type ChannelType `json:"type"`
		guildStoreChannel
	}{
		Type:              c.Type(),
		guildStoreChannel: guildStoreChannel(c),
	})
}

func (_ *GuildStoreChannel) Type() ChannelType {
	return ChannelTypeGuildStore
}
func (_ *GuildStoreChannel) guildChannel() {}

type GuildNewsThread struct {
	ID               Snowflake      `json:"id"`
	GuildID          Snowflake      `json:"guild_id"`
	Name             string         `json:"name"`
	LastMessageID    Snowflake      `json:"last_message_id"`
	LastPinTimestamp Time           `json:"last_pin_timestamp"`
	RateLimitPerUser int            `json:"rate_limit_per_user"`
	OwnerID          Snowflake      `json:"owner_id"`
	ParentID         Snowflake      `json:"parent_id"`
	MessageCount     int            `json:"message_count"`
	MemberCount      int            `json:"member_count"`
	ThreadMetadata   ThreadMetadata `json:"thread_metadata"`
}

func (c GuildNewsThread) MarshalJSON() ([]byte, error) {
	type guildNewsThread GuildNewsThread
	return json.Marshal(struct {
		Type ChannelType `json:"type"`
		guildNewsThread
	}{
		Type:            c.Type(),
		guildNewsThread: guildNewsThread(c),
	})
}

func (_ *GuildNewsThread) Type() ChannelType {
	return ChannelTypeGuildNewsThread
}
func (_ *GuildNewsThread) guildChannel()   {}
func (_ *GuildNewsThread) messageChannel() {}
func (_ *GuildNewsThread) guildThread()    {}

type GuildPublicThread struct {
	ID               Snowflake      `json:"id"`
	GuildID          Snowflake      `json:"guild_id"`
	Name             string         `json:"name"`
	LastMessageID    Snowflake      `json:"last_message_id"`
	LastPinTimestamp Time           `json:"last_pin_timestamp"`
	RateLimitPerUser int            `json:"rate_limit_per_user"`
	OwnerID          Snowflake      `json:"owner_id"`
	ParentID         Snowflake      `json:"parent_id"`
	MessageCount     int            `json:"message_count"`
	MemberCount      int            `json:"member_count"`
	ThreadMetadata   ThreadMetadata `json:"thread_metadata"`
}

func (c GuildPublicThread) MarshalJSON() ([]byte, error) {
	type guildTextChannel GuildPublicThread
	return json.Marshal(struct {
		Type ChannelType `json:"type"`
		guildTextChannel
	}{
		Type:             c.Type(),
		guildTextChannel: guildTextChannel(c),
	})
}

func (_ *GuildPublicThread) Type() ChannelType {
	return ChannelTypeGuildPublicThread
}
func (_ *GuildPublicThread) guildChannel()   {}
func (_ *GuildPublicThread) messageChannel() {}
func (_ *GuildPublicThread) guildThread()    {}

type GuildPrivateThread struct {
	ID               Snowflake      `json:"id"`
	GuildID          Snowflake      `json:"guild_id"`
	Name             string         `json:"name"`
	LastMessageID    Snowflake      `json:"last_message_id"`
	LastPinTimestamp Time           `json:"last_pin_timestamp"`
	RateLimitPerUser int            `json:"rate_limit_per_user"`
	OwnerID          Snowflake      `json:"owner_id"`
	ParentID         Snowflake      `json:"parent_id"`
	MessageCount     int            `json:"message_count"`
	MemberCount      int            `json:"member_count"`
	ThreadMetadata   ThreadMetadata `json:"thread_metadata"`
}

func (c GuildPrivateThread) MarshalJSON() ([]byte, error) {
	type guildPrivateThread GuildPrivateThread
	return json.Marshal(struct {
		Type ChannelType `json:"type"`
		guildPrivateThread
	}{
		Type:               c.Type(),
		guildPrivateThread: guildPrivateThread(c),
	})
}

func (_ *GuildPrivateThread) Type() ChannelType {
	return ChannelTypeGuildPrivateThread
}
func (_ *GuildPrivateThread) guildChannel()   {}
func (_ *GuildPrivateThread) messageChannel() {}
func (_ *GuildPrivateThread) guildThread()    {}

type GuildStageVoiceChannel struct {
	ID                     Snowflake             `json:"id"`
	GuildID                Snowflake             `json:"guild_id,omitempty"`
	Position               int                   `json:"position,omitempty"`
	PermissionOverwrites   []PermissionOverwrite `json:"permission_overwrites"`
	Name                   string                `json:"name,omitempty"`
	Topic                  *string               `json:"topic,omitempty"`
	Bitrate                int                   `json:"bitrate,omitempty"`
	UserLimit              int                   `json:"user_limit,omitempty"`
	ParentID               *Snowflake            `json:"parent_id,omitempty"`
	RTCRegion              string                `json:"rtc_region"`
	VideoQualityMode       VideoQualityMode      `json:"video_quality_mode"`
	InteractionPermissions Permissions           `json:"permissions,omitempty"`
}

func (c GuildStageVoiceChannel) MarshalJSON() ([]byte, error) {
	type guildStageVoiceChannel GuildStageVoiceChannel
	return json.Marshal(struct {
		Type ChannelType `json:"type"`
		guildStageVoiceChannel
	}{
		Type:                   c.Type(),
		guildStageVoiceChannel: guildStageVoiceChannel(c),
	})
}

func (_ *GuildStageVoiceChannel) Type() ChannelType {
	return ChannelTypeGuildStageVoice
}
func (_ *GuildStageVoiceChannel) guildChannel() {}
func (_ *GuildStageVoiceChannel) audioChannel() {}
func (_ *GuildStageVoiceChannel) guildThread()  {}

// VideoQualityMode https://discord.com/developers/docs/resources/channel#channel-object-video-quality-modes
type VideoQualityMode int

//goland:noinspection GoUnusedConst
const (
	VideoQualityModeAuto = iota + 1
	VideoQualityModeFull
)

type ThreadMetadata struct {
	Archived            bool                `json:"archived"`
	AutoArchiveDuration AutoArchiveDuration `json:"auto_archive_duration"`
	ArchiveTimestamp    Time                `json:"archive_timestamp"`
	Locked              bool                `json:"locked"`
	Invitable           bool                `json:"invitable"`
}

type AutoArchiveDuration int

//goland:noinspection GoUnusedConst
const (
	AutoArchiveDuration1h  AutoArchiveDuration = 60
	AutoArchiveDuration24h AutoArchiveDuration = 1440
	AutoArchiveDuration3d  AutoArchiveDuration = 4320
	AutoArchiveDuration1w  AutoArchiveDuration = 10080
)

// PartialChannel contains basic info about a Channel
type PartialChannel struct {
	ID   Snowflake   `json:"id"`
	Type ChannelType `json:"type"`
	Name string      `json:"name"`
	Icon *string     `json:"icon,omitempty"`
}

// GetIconURL returns the Icon URL of this channel.
// This will be nil for every discord.ChannelType except discord.ChannelTypeGroupDM
func (c *PartialChannel) GetIconURL(size int) *string {
	return FormatAssetURL(route.ChannelIcon, c.ID, c.Icon, size)
}
