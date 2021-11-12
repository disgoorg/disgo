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

func (t ChannelType) String() string {
	switch t {
	case ChannelTypeGuildText:
		return "GuildTextChannel"

	case ChannelTypeDM:
		return "DMChannel"

	case ChannelTypeGuildVoice:
		return "GuildVoiceChannel"

	case ChannelTypeGroupDM:
		return "GroupDMChannel"

	case ChannelTypeGuildCategory:
		return "GuildCategoryChannel"

	case ChannelTypeGuildNews:
		return "GuildNewsChannel"

	case ChannelTypeGuildStore:
		return "GuildStoreChannel"

	case ChannelTypeGuildNewsThread:
		return "GuildNewsThread"

	case ChannelTypeGuildPublicThread:
		return "GuildPublicThread"

	case ChannelTypeGuildPrivateThread:
		return "GuildPrivateThread"

	case ChannelTypeGuildStageVoice:
		return "GuildStageVoiceChannel"

	default:
		return "unknown"
	}
}

type Channel interface {
	json.Marshaler
	fmt.Stringer
	Type() ChannelType
	ID() Snowflake
	Name() string
	channel()
}

type GuildChannel interface {
	Channel
	Mentionable
	GuildID() Snowflake
	guildChannel()
}

type MessageChannel interface {
	Channel
	messageChannel()
}

type BaseGuildMessageChannel interface {
	GuildChannel
	MessageChannel
	baseGuildMessageChannel()
}

type GuildMessageChannel interface {
	BaseGuildMessageChannel
	guildMessageChannel()
}

type GuildThread interface {
	BaseGuildMessageChannel
	guildThread()
}

type GuildAudioChannel interface {
	GuildChannel
	guildAudioChannel()
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

	var (
		channel Channel
		err     error
	)

	switch cType.Type {
	case ChannelTypeGuildText:
		var v GuildTextChannel
		err = json.Unmarshal(data, &v)
		channel = v

	case ChannelTypeDM:
		var v DMChannel
		err = json.Unmarshal(data, &v)
		channel = v

	case ChannelTypeGuildVoice:
		var v GuildVoiceChannel
		err = json.Unmarshal(data, &v)
		channel = v

	case ChannelTypeGroupDM:
		var v GroupDMChannel
		err = json.Unmarshal(data, &v)
		channel = v

	case ChannelTypeGuildCategory:
		var v GuildCategoryChannel
		err = json.Unmarshal(data, &v)
		channel = v

	case ChannelTypeGuildNews:
		var v GuildNewsChannel
		err = json.Unmarshal(data, &v)
		channel = v

	case ChannelTypeGuildStore:
		var v GuildStoreChannel
		err = json.Unmarshal(data, &v)
		channel = v

	case ChannelTypeGuildNewsThread:
		var v GuildNewsThread
		err = json.Unmarshal(data, &v)
		channel = v

	case ChannelTypeGuildPublicThread:
		var v GuildPublicThread
		err = json.Unmarshal(data, &v)
		channel = v

	case ChannelTypeGuildPrivateThread:
		var v GuildPrivateThread
		err = json.Unmarshal(data, &v)
		channel = v

	case ChannelTypeGuildStageVoice:
		var v GuildStageVoiceChannel
		err = json.Unmarshal(data, &v)
		channel = v

	default:
		err = fmt.Errorf("unkown channel with type %d received", cType.Type)
	}

	if err != nil {
		return err
	}

	u.Channel = channel
	return nil
}

var (
	_ Channel                 = (*GuildTextChannel)(nil)
	_ GuildChannel            = (*GuildTextChannel)(nil)
	_ MessageChannel          = (*GuildTextChannel)(nil)
	_ BaseGuildMessageChannel = (*GuildTextChannel)(nil)
	_ GuildMessageChannel     = (*GuildTextChannel)(nil)
)

type GuildTextChannel struct {
	ChannelID                   Snowflake             `json:"id"`
	ChannelGuildID              Snowflake             `json:"guild_id,omitempty"`
	Position                    int                   `json:"position,omitempty"`
	ChannelPermissionOverwrites []PermissionOverwrite `json:"permission_overwrites"`
	ChannelName                 string                `json:"name,omitempty"`
	Topic                       *string               `json:"topic,omitempty"`
	NSFW                        bool                  `json:"nsfw,omitempty"`
	LastMessageID               *Snowflake            `json:"last_message_id,omitempty"`
	RateLimitPerUser            int                   `json:"rate_limit_per_user,omitempty"`
	ParentID                    *Snowflake            `json:"parent_id,omitempty"`
	LastPinTimestamp            *Time                 `json:"last_pin_timestamp,omitempty"`
	DefaultAutoArchiveDuration  AutoArchiveDuration   `json:"default_auto_archive_duration"`
	InteractionPermissions      Permissions           `json:"permissions,omitempty"`
}

func (c *GuildTextChannel) UnmarshalJSON(data []byte) error {
	type guildTextChannel GuildTextChannel
	var v struct {
		ChannelPermissionOverwrites []UnmarshalPermissionOverwrite `json:"permission_overwrites"`
		guildTextChannel
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*c = GuildTextChannel(v.guildTextChannel)
	c.ChannelPermissionOverwrites = parsePermissionOverwrites(v.ChannelPermissionOverwrites)
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

func (c GuildTextChannel) String() string {
	return channelString(c)
}

func (c GuildTextChannel) Mention() string {
	return channelMention(c.ID())
}

func (_ GuildTextChannel) Type() ChannelType {
	return ChannelTypeGuildText
}

func (c GuildTextChannel) Name() string {
	return c.ChannelName
}

func (c GuildTextChannel) ID() Snowflake {
	return c.ChannelID
}

func (c GuildTextChannel) GuildID() Snowflake {
	return c.ChannelGuildID
}

func (_ GuildTextChannel) channel()                 {}
func (_ GuildTextChannel) guildChannel()            {}
func (_ GuildTextChannel) messageChannel()          {}
func (_ GuildTextChannel) baseGuildMessageChannel() {}
func (_ GuildTextChannel) guildMessageChannel()     {}

var (
	_ Channel        = (*DMChannel)(nil)
	_ MessageChannel = (*DMChannel)(nil)
)

type DMChannel struct {
	ChannelID        Snowflake  `json:"id"`
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

func (c DMChannel) String() string {
	return channelString(c)
}

func (_ DMChannel) Type() ChannelType {
	return ChannelTypeGuildText
}

func (c DMChannel) ID() Snowflake {
	return c.ChannelID
}

func (c DMChannel) Name() string {
	return ""
}

func (_ DMChannel) channel()        {}
func (_ DMChannel) messageChannel() {}

var (
	_ Channel           = (*GuildVoiceChannel)(nil)
	_ GuildChannel      = (*GuildVoiceChannel)(nil)
	_ GuildAudioChannel = (*GuildVoiceChannel)(nil)
)

type GuildVoiceChannel struct {
	ChannelID                   Snowflake             `json:"id"`
	ChannelGuildID              Snowflake             `json:"guild_id,omitempty"`
	Position                    int                   `json:"position,omitempty"`
	ChannelPermissionOverwrites []PermissionOverwrite `json:"permission_overwrites"`
	ChannelName                 string                `json:"name,omitempty"`
	NSFW                        bool                  `json:"nsfw,omitempty"`
	Topic                       *string               `json:"topic,omitempty"`
	Bitrate                     int                   `json:"bitrate,omitempty"`
	UserLimit                   int                   `json:"user_limit,omitempty"`
	ParentID                    *Snowflake            `json:"parent_id,omitempty"`
	RTCRegion                   string                `json:"rtc_region"`
	VideoQualityMode            VideoQualityMode      `json:"video_quality_mode"`
	InteractionPermissions      Permissions           `json:"permissions,omitempty"`
}

func (c *GuildVoiceChannel) UnmarshalJSON(data []byte) error {
	type guildVoiceChannel GuildVoiceChannel
	var v struct {
		ChannelPermissionOverwrites []UnmarshalPermissionOverwrite `json:"permission_overwrites"`
		guildVoiceChannel
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*c = GuildVoiceChannel(v.guildVoiceChannel)
	c.ChannelPermissionOverwrites = parsePermissionOverwrites(v.ChannelPermissionOverwrites)
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

func (c GuildVoiceChannel) String() string {
	return channelString(c)
}

func (c GuildVoiceChannel) Mention() string {
	return channelMention(c.ID())
}

func (_ GuildVoiceChannel) Type() ChannelType {
	return ChannelTypeGuildText
}

func (c GuildVoiceChannel) ID() Snowflake {
	return c.ChannelID
}

func (c GuildVoiceChannel) Name() string {
	return c.ChannelName
}

func (c GuildVoiceChannel) GuildID() Snowflake {
	return c.ChannelGuildID
}

func (_ GuildVoiceChannel) channel()           {}
func (_ GuildVoiceChannel) guildChannel()      {}
func (_ GuildVoiceChannel) guildAudioChannel() {}

var (
	_ Channel        = (*GroupDMChannel)(nil)
	_ MessageChannel = (*GroupDMChannel)(nil)
)

type GroupDMChannel struct {
	ChannelID        Snowflake  `json:"id"`
	ChannelName      string     `json:"name,omitempty"`
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

func (c GroupDMChannel) String() string {
	return channelString(c)
}

func (_ GroupDMChannel) Type() ChannelType {
	return ChannelTypeGuildText
}
func (c GroupDMChannel) ID() Snowflake {
	return c.ChannelID
}

func (c GroupDMChannel) Name() string {
	return c.ChannelName
}

func (_ GroupDMChannel) channel()        {}
func (_ GroupDMChannel) messageChannel() {}

var (
	_ Channel      = (*GuildCategoryChannel)(nil)
	_ GuildChannel = (*GuildCategoryChannel)(nil)
)

type GuildCategoryChannel struct {
	ChannelID                   Snowflake             `json:"id"`
	ChannelGuildID              Snowflake             `json:"guild_id"`
	Position                    int                   `json:"position"`
	ChannelPermissionOverwrites []PermissionOverwrite `json:"permission_overwrites"`
	ChannelName                 string                `json:"name"`
	NSFW                        bool                  `json:"nsfw"`
	InteractionPermissions      Permissions           `json:"permissions,omitempty"`
}

func (c *GuildCategoryChannel) UnmarshalJSON(data []byte) error {
	type guildCategoryChannel GuildCategoryChannel
	var v struct {
		ChannelPermissionOverwrites []UnmarshalPermissionOverwrite `json:"permission_overwrites"`
		guildCategoryChannel
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*c = GuildCategoryChannel(v.guildCategoryChannel)
	c.ChannelPermissionOverwrites = parsePermissionOverwrites(v.ChannelPermissionOverwrites)
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

func (c GuildCategoryChannel) String() string {
	return channelString(c)
}

func (c GuildCategoryChannel) Mention() string {
	return channelMention(c.ID())
}

func (_ GuildCategoryChannel) Type() ChannelType {
	return ChannelTypeGuildCategory
}

func (c GuildCategoryChannel) ID() Snowflake {
	return c.ChannelID
}

func (c GuildCategoryChannel) Name() string {
	return c.ChannelName
}

func (c GuildCategoryChannel) GuildID() Snowflake {
	return c.ChannelGuildID
}

func (_ GuildCategoryChannel) channel()      {}
func (_ GuildCategoryChannel) guildChannel() {}

var (
	_ Channel                 = (*GuildNewsChannel)(nil)
	_ GuildChannel            = (*GuildNewsChannel)(nil)
	_ MessageChannel          = (*GuildNewsChannel)(nil)
	_ BaseGuildMessageChannel = (*GuildNewsChannel)(nil)
	_ GuildMessageChannel     = (*GuildNewsChannel)(nil)
)

type GuildNewsChannel struct {
	ChannelID                   Snowflake             `json:"id"`
	ChannelGuildID              Snowflake             `json:"guild_id,omitempty"`
	Position                    int                   `json:"position,omitempty"`
	ChannelPermissionOverwrites []PermissionOverwrite `json:"permission_overwrites"`
	ChannelName                 string                `json:"name,omitempty"`
	Topic                       *string               `json:"topic,omitempty"`
	NSFW                        bool                  `json:"nsfw,omitempty"`
	LastMessageID               *Snowflake            `json:"last_message_id,omitempty"`
	RateLimitPerUser            int                   `json:"rate_limit_per_user,omitempty"`
	ParentID                    *Snowflake            `json:"parent_id,omitempty"`
	LastPinTimestamp            *Time                 `json:"last_pin_timestamp,omitempty"`
	DefaultAutoArchiveDuration  AutoArchiveDuration   `json:"default_auto_archive_duration"`
	InteractionPermissions      Permissions           `json:"permissions,omitempty"`
}

func (c *GuildNewsChannel) UnmarshalJSON(data []byte) error {
	type guildNewsChannel GuildNewsChannel
	var v struct {
		ChannelPermissionOverwrites []UnmarshalPermissionOverwrite `json:"permission_overwrites"`
		guildNewsChannel
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*c = GuildNewsChannel(v.guildNewsChannel)
	c.ChannelPermissionOverwrites = parsePermissionOverwrites(v.ChannelPermissionOverwrites)
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

func (c GuildNewsChannel) String() string {
	return channelString(c)
}

func (c GuildNewsChannel) Mention() string {
	return channelMention(c.ID())
}

func (_ GuildNewsChannel) Type() ChannelType {
	return ChannelTypeGuildNews
}

func (c GuildNewsChannel) ID() Snowflake {
	return c.ChannelID
}

func (c GuildNewsChannel) Name() string {
	return c.ChannelName
}

func (c GuildNewsChannel) GuildID() Snowflake {
	return c.ChannelGuildID
}

func (_ GuildNewsChannel) channel()                 {}
func (_ GuildNewsChannel) guildChannel()            {}
func (_ GuildNewsChannel) messageChannel()          {}
func (_ GuildNewsChannel) baseGuildMessageChannel() {}
func (_ GuildNewsChannel) guildMessageChannel()     {}

var (
	_ Channel      = (*GuildStoreChannel)(nil)
	_ GuildChannel = (*GuildStoreChannel)(nil)
)

type GuildStoreChannel struct {
	ChannelID                   Snowflake             `json:"id"`
	ChannelGuildID              Snowflake             `json:"guild_id"`
	Position                    int                   `json:"position"`
	ChannelPermissionOverwrites []PermissionOverwrite `json:"permission_overwrites"`
	ChannelName                 string                `json:"name"`
	NSFW                        bool                  `json:"nsfw,omitempty"`
	ParentID                    *Snowflake            `json:"parent_id"`
	InteractionPermissions      Permissions           `json:"permissions,omitempty"`
}

func (c *GuildStoreChannel) UnmarshalJSON(data []byte) error {
	type guildStoreChannel GuildStoreChannel
	var v struct {
		ChannelPermissionOverwrites []UnmarshalPermissionOverwrite `json:"permission_overwrites"`
		guildStoreChannel
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*c = GuildStoreChannel(v.guildStoreChannel)
	c.ChannelPermissionOverwrites = parsePermissionOverwrites(v.ChannelPermissionOverwrites)
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

func (c GuildStoreChannel) String() string {
	return channelString(c)
}

func (c GuildStoreChannel) Mention() string {
	return channelMention(c.ID())
}

func (_ GuildStoreChannel) Type() ChannelType {
	return ChannelTypeGuildStore
}

func (c GuildStoreChannel) ID() Snowflake {
	return c.ChannelID
}

func (c GuildStoreChannel) Name() string {
	return c.ChannelName
}

func (c GuildStoreChannel) GuildID() Snowflake {
	return c.ChannelGuildID
}

func (_ GuildStoreChannel) channel()      {}
func (_ GuildStoreChannel) guildChannel() {}

var (
	_ Channel                 = (*GuildNewsThread)(nil)
	_ GuildChannel            = (*GuildNewsThread)(nil)
	_ MessageChannel          = (*GuildNewsThread)(nil)
	_ BaseGuildMessageChannel = (*GuildNewsThread)(nil)
	_ GuildThread             = (*GuildNewsThread)(nil)
)

type GuildNewsThread struct {
	ChannelID        Snowflake      `json:"id"`
	ChannelGuildID   Snowflake      `json:"guild_id"`
	ChannelName      string         `json:"name"`
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

func (c GuildNewsThread) String() string {
	return channelString(c)
}

func (c GuildNewsThread) Mention() string {
	return channelMention(c.ID())
}

func (_ GuildNewsThread) Type() ChannelType {
	return ChannelTypeGuildNewsThread
}

func (c GuildNewsThread) ID() Snowflake {
	return c.ChannelID
}

func (c GuildNewsThread) Name() string {
	return c.ChannelName
}

func (c GuildNewsThread) GuildID() Snowflake {
	return c.ChannelGuildID
}

func (_ GuildNewsThread) channel()                 {}
func (_ GuildNewsThread) guildChannel()            {}
func (_ GuildNewsThread) messageChannel()          {}
func (_ GuildNewsThread) baseGuildMessageChannel() {}
func (_ GuildNewsThread) guildThread()             {}

var (
	_ Channel                 = (*GuildPublicThread)(nil)
	_ GuildChannel            = (*GuildPublicThread)(nil)
	_ MessageChannel          = (*GuildPublicThread)(nil)
	_ BaseGuildMessageChannel = (*GuildPublicThread)(nil)
	_ GuildThread             = (*GuildPublicThread)(nil)
)

type GuildPublicThread struct {
	ChannelID        Snowflake      `json:"id"`
	ChannelGuildID   Snowflake      `json:"guild_id"`
	ChannelName      string         `json:"name"`
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

func (c GuildPublicThread) String() string {
	return channelString(c)
}

func (c GuildPublicThread) Mention() string {
	return channelMention(c.ID())
}

func (_ GuildPublicThread) Type() ChannelType {
	return ChannelTypeGuildPublicThread
}

func (c GuildPublicThread) ID() Snowflake {
	return c.ChannelID
}

func (c GuildPublicThread) Name() string {
	return c.ChannelName
}

func (c GuildPublicThread) GuildID() Snowflake {
	return c.ChannelGuildID
}

func (_ GuildPublicThread) channel()                 {}
func (_ GuildPublicThread) guildChannel()            {}
func (_ GuildPublicThread) messageChannel()          {}
func (_ GuildPublicThread) baseGuildMessageChannel() {}
func (_ GuildPublicThread) guildThread()             {}

var (
	_ Channel                 = (*GuildPrivateThread)(nil)
	_ GuildChannel            = (*GuildPrivateThread)(nil)
	_ MessageChannel          = (*GuildPrivateThread)(nil)
	_ BaseGuildMessageChannel = (*GuildPrivateThread)(nil)
	_ GuildThread             = (*GuildPrivateThread)(nil)
)

type GuildPrivateThread struct {
	ChannelID        Snowflake      `json:"id"`
	ChannelGuildID   Snowflake      `json:"guild_id"`
	ChannelName      string         `json:"name"`
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

func (c GuildPrivateThread) String() string {
	return channelString(c)
}

func (c GuildPrivateThread) Mention() string {
	return channelMention(c.ID())
}

func (_ GuildPrivateThread) Type() ChannelType {
	return ChannelTypeGuildPrivateThread
}

func (c GuildPrivateThread) ID() Snowflake {
	return c.ChannelID
}

func (c GuildPrivateThread) Name() string {
	return c.ChannelName
}

func (c GuildPrivateThread) GuildID() Snowflake {
	return c.ChannelGuildID
}

func (_ GuildPrivateThread) channel()                 {}
func (_ GuildPrivateThread) guildChannel()            {}
func (_ GuildPrivateThread) messageChannel()          {}
func (_ GuildPrivateThread) baseGuildMessageChannel() {}
func (_ GuildPrivateThread) guildThread()             {}

var (
	_ Channel           = (*GuildStageVoiceChannel)(nil)
	_ GuildChannel      = (*GuildStageVoiceChannel)(nil)
	_ GuildAudioChannel = (*GuildStageVoiceChannel)(nil)
)

type GuildStageVoiceChannel struct {
	ChannelID                   Snowflake             `json:"id"`
	ChannelGuildID              Snowflake             `json:"guild_id,omitempty"`
	Position                    int                   `json:"position,omitempty"`
	ChannelPermissionOverwrites []PermissionOverwrite `json:"permission_overwrites"`
	ChannelName                 string                `json:"name,omitempty"`
	Topic                       *string               `json:"topic,omitempty"`
	Bitrate                     int                   `json:"bitrate,omitempty"`
	UserLimit                   int                   `json:"user_limit,omitempty"`
	ParentID                    *Snowflake            `json:"parent_id,omitempty"`
	RTCRegion                   string                `json:"rtc_region"`
	VideoQualityMode            VideoQualityMode      `json:"video_quality_mode"`
	InteractionPermissions      Permissions           `json:"permissions,omitempty"`
}

func (c *GuildStageVoiceChannel) UnmarshalJSON(data []byte) error {
	type guildStageVoiceChannel GuildStageVoiceChannel
	var v struct {
		ChannelPermissionOverwrites []UnmarshalPermissionOverwrite `json:"permission_overwrites"`
		guildStageVoiceChannel
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*c = GuildStageVoiceChannel(v.guildStageVoiceChannel)
	c.ChannelPermissionOverwrites = parsePermissionOverwrites(v.ChannelPermissionOverwrites)
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

func (c GuildStageVoiceChannel) String() string {
	return channelString(c)
}

func (c GuildStageVoiceChannel) Mention() string {
	return channelMention(c.ID())
}

func (_ GuildStageVoiceChannel) Type() ChannelType {
	return ChannelTypeGuildStageVoice
}

func (c GuildStageVoiceChannel) ID() Snowflake {
	return c.ChannelID
}

func (c GuildStageVoiceChannel) Name() string {
	return c.ChannelName
}

func (c GuildStageVoiceChannel) GuildID() Snowflake {
	return c.ChannelGuildID
}

func (_ GuildStageVoiceChannel) channel()           {}
func (_ GuildStageVoiceChannel) guildChannel()      {}
func (_ GuildStageVoiceChannel) guildAudioChannel() {}

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

func channelString(channel Channel) string {
	return fmt.Sprintf("%s:%s(%s)", channel.Type().String(), channel.Name(), channel.ID())
}
