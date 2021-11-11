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
	channel()
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

func (c *GuildTextChannel) UnmarshalJSON(data []byte) error {
	type guildTextChannel GuildTextChannel
	var v struct {
		PermissionOverwrites []UnmarshalPermissionOverwrite `json:"permission_overwrites"`
		guildTextChannel
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*c = GuildTextChannel(v.guildTextChannel)
	c.PermissionOverwrites = parsePermissionOverwrites(v.PermissionOverwrites)
	return nil
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

func (_ GuildTextChannel) Type() ChannelType {
	return ChannelTypeGuildText
}
func (_ GuildTextChannel) channel()             {}
func (_ GuildTextChannel) messageChannel()      {}
func (_ GuildTextChannel) guildMessageChannel() {}
func (_ GuildTextChannel) guildChannel()        {}

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

func (_ DMChannel) Type() ChannelType {
	return ChannelTypeGuildText
}
func (_ DMChannel) channel()        {}
func (_ DMChannel) messageChannel() {}

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

func (c *GuildVoiceChannel) UnmarshalJSON(data []byte) error {
	type guildVoiceChannel GuildVoiceChannel
	var v struct {
		PermissionOverwrites []UnmarshalPermissionOverwrite `json:"permission_overwrites"`
		guildVoiceChannel
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*c = GuildVoiceChannel(v.guildVoiceChannel)
	c.PermissionOverwrites = parsePermissionOverwrites(v.PermissionOverwrites)
	return nil
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

func (_ GuildVoiceChannel) Type() ChannelType {
	return ChannelTypeGuildText
}
func (_ GuildVoiceChannel) channel()      {}
func (_ GuildVoiceChannel) guildChannel() {}
func (_ GuildVoiceChannel) audioChannel() {}

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

func (_ GroupDMChannel) Type() ChannelType {
	return ChannelTypeGuildText
}
func (_ GroupDMChannel) channel()        {}
func (_ GroupDMChannel) messageChannel() {}

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

func (c *GuildCategoryChannel) UnmarshalJSON(data []byte) error {
	type guildCategoryChannel GuildCategoryChannel
	var v struct {
		PermissionOverwrites []UnmarshalPermissionOverwrite `json:"permission_overwrites"`
		guildCategoryChannel
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*c = GuildCategoryChannel(v.guildCategoryChannel)
	c.PermissionOverwrites = parsePermissionOverwrites(v.PermissionOverwrites)
	return nil
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

func (_ GuildCategoryChannel) Type() ChannelType {
	return ChannelTypeGuildCategory
}
func (_ GuildCategoryChannel) channel()      {}
func (_ GuildCategoryChannel) guildChannel() {}

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

func (c *GuildNewsChannel) UnmarshalJSON(data []byte) error {
	type guildNewsChannel GuildNewsChannel
	var v struct {
		PermissionOverwrites []UnmarshalPermissionOverwrite `json:"permission_overwrites"`
		guildNewsChannel
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*c = GuildNewsChannel(v.guildNewsChannel)
	c.PermissionOverwrites = parsePermissionOverwrites(v.PermissionOverwrites)
	return nil
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

func (_ GuildNewsChannel) Type() ChannelType {
	return ChannelTypeGuildNews
}
func (_ GuildNewsChannel) channel()        {}
func (_ GuildNewsChannel) guildChannel()   {}
func (_ GuildNewsChannel) messageChannel() {}

type GuildStoreChannel struct {
	ID                     Snowflake             `json:"id"`
	GuildID                Snowflake             `json:"guild_id"`
	Position               int                   `json:"position"`
	PermissionOverwrites   []PermissionOverwrite `json:"permission_overwrites"`
	Name                   string                `json:"name"`
	NSFW                   bool                  `json:"nsfw,omitempty"`
	ParentID               *Snowflake            `json:"parent_id"`
	InteractionPermissions Permissions           `json:"permissions,omitempty"`
}

func (c *GuildStoreChannel) UnmarshalJSON(data []byte) error {
	type guildStoreChannel GuildStoreChannel
	var v struct {
		PermissionOverwrites []UnmarshalPermissionOverwrite `json:"permission_overwrites"`
		guildStoreChannel
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*c = GuildStoreChannel(v.guildStoreChannel)
	c.PermissionOverwrites = parsePermissionOverwrites(v.PermissionOverwrites)
	return nil
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

func (_ GuildStoreChannel) Type() ChannelType {
	return ChannelTypeGuildStore
}
func (_ GuildStoreChannel) channel()      {}
func (_ GuildStoreChannel) guildChannel() {}

type GuildNewsThread struct {
	ID               Snowflake      `json:"id"`
	GuildID          Snowflake      `json:"guild_id"`
	Name             string         `json:"name"`
	LastMessageID    Snowflake      `json:"last_message_id"`
	LastPinTimestamp *Time          `json:"last_pin_timestamp"`
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

func (_ GuildNewsThread) Type() ChannelType {
	return ChannelTypeGuildNewsThread
}
func (_ GuildNewsThread) channel()        {}
func (_ GuildNewsThread) guildChannel()   {}
func (_ GuildNewsThread) messageChannel() {}
func (_ GuildNewsThread) guildThread()    {}

type GuildPublicThread struct {
	ID               Snowflake      `json:"id"`
	GuildID          Snowflake      `json:"guild_id"`
	Name             string         `json:"name"`
	LastMessageID    Snowflake      `json:"last_message_id"`
	LastPinTimestamp *Time          `json:"last_pin_timestamp"`
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

func (_ GuildPublicThread) Type() ChannelType {
	return ChannelTypeGuildPublicThread
}
func (_ GuildPublicThread) channel()        {}
func (_ GuildPublicThread) guildChannel()   {}
func (_ GuildPublicThread) messageChannel() {}
func (_ GuildPublicThread) guildThread()    {}

type GuildPrivateThread struct {
	ID               Snowflake      `json:"id"`
	GuildID          Snowflake      `json:"guild_id"`
	Name             string         `json:"name"`
	LastMessageID    Snowflake      `json:"last_message_id"`
	LastPinTimestamp *Time          `json:"last_pin_timestamp"`
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

func (_ GuildPrivateThread) Type() ChannelType {
	return ChannelTypeGuildPrivateThread
}
func (_ GuildPrivateThread) channel()        {}
func (_ GuildPrivateThread) guildChannel()   {}
func (_ GuildPrivateThread) messageChannel() {}
func (_ GuildPrivateThread) guildThread()    {}

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

func (c *GuildStageVoiceChannel) UnmarshalJSON(data []byte) error {
	type guildStageVoiceChannel GuildStageVoiceChannel
	var v struct {
		PermissionOverwrites []UnmarshalPermissionOverwrite `json:"permission_overwrites"`
		guildStageVoiceChannel
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*c = GuildStageVoiceChannel(v.guildStageVoiceChannel)
	c.PermissionOverwrites = parsePermissionOverwrites(v.PermissionOverwrites)
	return nil
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

func (_ GuildStageVoiceChannel) Type() ChannelType {
	return ChannelTypeGuildStageVoice
}
func (_ GuildStageVoiceChannel) channel()      {}
func (_ GuildStageVoiceChannel) guildChannel() {}
func (_ GuildStageVoiceChannel) audioChannel() {}

// VideoQualityMode https://com/developers/docs/resources/channel#channel-object-video-quality-modes
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
// This will be nil for every ChannelType except ChannelTypeGroupDM
func (c *PartialChannel) GetIconURL(size int) *string {
	return FormatAssetURL(route.ChannelIcon, c.ID, c.Icon, size)
}

func parsePermissionOverwrites(unmarshalOverwrites []UnmarshalPermissionOverwrite) []PermissionOverwrite {
	overwrites := make([]PermissionOverwrite, len(unmarshalOverwrites))
	for i := range unmarshalOverwrites {
		overwrites[i] = unmarshalOverwrites[i].PermissionOverwrite
	}
	return overwrites
}

func ChannelID(channel Channel) Snowflake {
	if channel == nil {
		return ""
	}
	switch ch := channel.(type) {
	case GuildTextChannel:
		return ch.ID

	case DMChannel:
		return ch.ID

	case GuildVoiceChannel:
		return ch.ID

	case GroupDMChannel:
		return ch.ID

	case GuildCategoryChannel:
		return ch.ID

	case GuildNewsChannel:
		return ch.ID

	case GuildStoreChannel:
		return ch.ID

	case GuildNewsThread:
		return ch.ID

	case GuildPrivateThread:
		return ch.ID

	case GuildPublicThread:
		return ch.ID

	case GuildStageVoiceChannel:
		return ch.ID

	default:
		panic("unknown channel type")
	}
}

func GuildID(channel GuildChannel) Snowflake {
	if channel == nil {
		return ""
	}
	switch ch := channel.(type) {
	case GuildTextChannel:
		return ch.GuildID

	case GuildVoiceChannel:
		return ch.GuildID

	case GuildCategoryChannel:
		return ch.GuildID

	case GuildNewsChannel:
		return ch.GuildID

	case GuildStoreChannel:
		return ch.GuildID

	case GuildNewsThread:
		return ch.GuildID

	case GuildPrivateThread:
		return ch.GuildID

	case GuildPublicThread:
		return ch.GuildID

	case GuildStageVoiceChannel:
		return ch.GuildID

	default:
		panic("unknown channel type")
	}
}

func LastPinTimestamp(channel MessageChannel) *Time {
	if channel == nil {
		return nil
	}
	switch ch := channel.(type) {
	case GuildTextChannel:
		return ch.LastPinTimestamp

	case DMChannel:
		return ch.LastPinTimestamp

	case GroupDMChannel:
		return ch.LastPinTimestamp

	case GuildNewsChannel:
		return ch.LastPinTimestamp

	case GuildNewsThread:
		return ch.LastPinTimestamp

	case GuildPrivateThread:
		return ch.LastPinTimestamp

	case GuildPublicThread:
		return ch.LastPinTimestamp

	default:
		panic("unknown channel type")
	}
}
