package discord

import (
	"fmt"
	"time"

	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/internal/flags"
)

// ChannelType for interacting with discord's channels
type ChannelType int

// Channel constants
const (
	ChannelTypeGuildText ChannelType = iota
	ChannelTypeDM
	ChannelTypeGuildVoice
	ChannelTypeGroupDM
	ChannelTypeGuildCategory
	ChannelTypeGuildNews
	_
	_
	_
	_
	ChannelTypeGuildNewsThread
	ChannelTypeGuildPublicThread
	ChannelTypeGuildPrivateThread
	ChannelTypeGuildStageVoice
	ChannelTypeGuildDirectory
	ChannelTypeGuildForum
	ChannelTypeGuildMedia
)

type ChannelFlags int

const (
	ChannelFlagPinned ChannelFlags = 1 << (iota + 1)
	_
	_
	ChannelFlagRequireTag
	ChannelFlagHideMediaDownloadOptions ChannelFlags = 1 << 15
	ChannelFlagsNone                    ChannelFlags = 0
)

// Add allows you to add multiple bits together, producing a new bit
func (f ChannelFlags) Add(bits ...ChannelFlags) ChannelFlags {
	return flags.Add(f, bits...)
}

// Remove allows you to subtract multiple bits from the first, producing a new bit
func (f ChannelFlags) Remove(bits ...ChannelFlags) ChannelFlags {
	return flags.Remove(f, bits...)
}

// Has will ensure that the bit includes all the bits entered
func (f ChannelFlags) Has(bits ...ChannelFlags) bool {
	return flags.Has(f, bits...)
}

// Missing will check whether the bit is missing any one of the bits
func (f ChannelFlags) Missing(bits ...ChannelFlags) bool {
	return flags.Missing(f, bits...)
}

type Channel interface {
	json.Marshaler
	fmt.Stringer

	// Type returns the ChannelType of the Channel.
	Type() ChannelType

	// ID returns the Snowflake ID of the Channel.
	ID() snowflake.ID

	// Name returns the name of the Channel.
	Name() string

	// CreatedAt returns the creation time of the Channel.
	CreatedAt() time.Time

	channel()
}

type MessageChannel interface {
	Channel

	// LastMessageID returns the ID of the last Message sent in this MessageChannel.
	// This is nil if no Message has been sent yet.
	LastMessageID() *snowflake.ID

	// LastPinTimestamp returns when the last Message in this MessageChannel was pinned.
	// This is nil if no Message has been pinned yet.
	LastPinTimestamp() *time.Time

	messageChannel()
}

type GuildChannel interface {
	Channel
	Mentionable

	// GuildID returns the Guild ID of the GuildChannel
	GuildID() snowflake.ID

	// Position returns the position of the GuildChannel in the channel list.
	// This is always 0 for GuildThread(s).
	Position() int

	// ParentID returns the parent Channel ID of the GuildChannel.
	// This is never nil for GuildThread(s).
	ParentID() *snowflake.ID

	// PermissionOverwrites returns the GuildChannel's PermissionOverwrites for Role(s) and Member(s).
	// This is always nil for GuildThread(s).
	PermissionOverwrites() PermissionOverwrites

	guildChannel()
}

type GuildMessageChannel interface {
	GuildChannel
	MessageChannel

	// Topic returns the topic of a GuildMessageChannel.
	// This is always nil for GuildThread(s).
	Topic() *string

	// NSFW returns whether the GuildMessageChannel is marked as not safe for work.
	NSFW() bool

	// DefaultAutoArchiveDuration returns the default AutoArchiveDuration for GuildThread(s) in this GuildMessageChannel.
	// This is always 0 for GuildThread(s).
	DefaultAutoArchiveDuration() AutoArchiveDuration
	RateLimitPerUser() int

	guildMessageChannel()
}

type GuildAudioChannel interface {
	GuildChannel

	// Bitrate returns the configured bitrate of the GuildAudioChannel.
	Bitrate() int

	// RTCRegion returns the configured voice server region of the GuildAudioChannel.
	RTCRegion() string

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

	case ChannelTypeGuildNewsThread, ChannelTypeGuildPublicThread, ChannelTypeGuildPrivateThread:
		var v GuildThread
		err = json.Unmarshal(data, &v)
		channel = v

	case ChannelTypeGuildStageVoice:
		var v GuildStageVoiceChannel
		err = json.Unmarshal(data, &v)
		channel = v

	case ChannelTypeGuildForum:
		var v GuildForumChannel
		err = json.Unmarshal(data, &v)
		channel = v

	case ChannelTypeGuildMedia:
		var v GuildMediaChannel
		err = json.Unmarshal(data, &v)
		channel = v

	default:
		err = fmt.Errorf("unknown channel with type %d received", cType.Type)
	}

	if err != nil {
		return err
	}

	u.Channel = channel
	return nil
}

var (
	_ Channel             = (*GuildTextChannel)(nil)
	_ GuildChannel        = (*GuildTextChannel)(nil)
	_ MessageChannel      = (*GuildTextChannel)(nil)
	_ GuildMessageChannel = (*GuildTextChannel)(nil)
)

type GuildTextChannel struct {
	id                         snowflake.ID
	guildID                    snowflake.ID
	position                   int
	permissionOverwrites       PermissionOverwrites
	name                       string
	topic                      *string
	nsfw                       bool
	lastMessageID              *snowflake.ID
	rateLimitPerUser           int
	parentID                   *snowflake.ID
	lastPinTimestamp           *time.Time
	defaultAutoArchiveDuration AutoArchiveDuration
}

func (c *GuildTextChannel) UnmarshalJSON(data []byte) error {
	var v guildTextChannel
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	c.id = v.ID
	c.guildID = v.GuildID
	c.position = v.Position
	c.permissionOverwrites = v.PermissionOverwrites
	c.name = v.Name
	c.topic = v.Topic
	c.nsfw = v.NSFW
	c.lastMessageID = v.LastMessageID
	c.rateLimitPerUser = v.RateLimitPerUser
	c.parentID = v.ParentID
	c.lastPinTimestamp = v.LastPinTimestamp
	c.defaultAutoArchiveDuration = v.DefaultAutoArchiveDuration
	return nil
}

func (c GuildTextChannel) MarshalJSON() ([]byte, error) {
	return json.Marshal(guildTextChannel{
		ID:                         c.id,
		Type:                       c.Type(),
		GuildID:                    c.guildID,
		Position:                   c.position,
		PermissionOverwrites:       c.permissionOverwrites,
		Name:                       c.name,
		Topic:                      c.topic,
		NSFW:                       c.nsfw,
		LastMessageID:              c.lastMessageID,
		RateLimitPerUser:           c.rateLimitPerUser,
		ParentID:                   c.parentID,
		LastPinTimestamp:           c.lastPinTimestamp,
		DefaultAutoArchiveDuration: c.defaultAutoArchiveDuration,
	})
}

func (c GuildTextChannel) String() string {
	return channelString(c)
}

func (c GuildTextChannel) Mention() string {
	return ChannelMention(c.ID())
}

func (c GuildTextChannel) ID() snowflake.ID {
	return c.id
}

func (GuildTextChannel) Type() ChannelType {
	return ChannelTypeGuildText
}

func (c GuildTextChannel) Name() string {
	return c.name
}

func (c GuildTextChannel) GuildID() snowflake.ID {
	return c.guildID
}

func (c GuildTextChannel) PermissionOverwrites() PermissionOverwrites {
	return c.permissionOverwrites
}

func (c GuildTextChannel) Position() int {
	return c.position
}

func (c GuildTextChannel) ParentID() *snowflake.ID {
	return c.parentID
}

func (c GuildTextChannel) LastMessageID() *snowflake.ID {
	return c.lastMessageID
}

func (c GuildTextChannel) RateLimitPerUser() int {
	return c.rateLimitPerUser
}

func (c GuildTextChannel) LastPinTimestamp() *time.Time {
	return c.lastPinTimestamp
}

func (c GuildTextChannel) Topic() *string {
	return c.topic
}

func (c GuildTextChannel) NSFW() bool {
	return c.nsfw
}

func (c GuildTextChannel) DefaultAutoArchiveDuration() AutoArchiveDuration {
	return c.defaultAutoArchiveDuration
}

func (c GuildTextChannel) CreatedAt() time.Time {
	return c.id.Time()
}

func (GuildTextChannel) channel()             {}
func (GuildTextChannel) guildChannel()        {}
func (GuildTextChannel) messageChannel()      {}
func (GuildTextChannel) guildMessageChannel() {}

var (
	_ Channel        = (*DMChannel)(nil)
	_ MessageChannel = (*DMChannel)(nil)
)

type DMChannel struct {
	id               snowflake.ID
	lastMessageID    *snowflake.ID
	recipients       []User
	lastPinTimestamp *time.Time
}

func (c *DMChannel) UnmarshalJSON(data []byte) error {
	var v dmChannel
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	c.id = v.ID
	c.lastMessageID = v.LastMessageID
	c.recipients = v.Recipients
	c.lastPinTimestamp = v.LastPinTimestamp
	return nil
}

func (c DMChannel) MarshalJSON() ([]byte, error) {
	return json.Marshal(dmChannel{
		ID:               c.id,
		Type:             c.Type(),
		LastMessageID:    c.lastMessageID,
		Recipients:       c.recipients,
		LastPinTimestamp: c.lastPinTimestamp,
	})
}

func (c DMChannel) String() string {
	return channelString(c)
}

func (c DMChannel) ID() snowflake.ID {
	return c.id
}

func (DMChannel) Type() ChannelType {
	return ChannelTypeDM
}

func (c DMChannel) Name() string {
	return c.recipients[0].Username
}

func (c DMChannel) LastMessageID() *snowflake.ID {
	return c.lastMessageID
}

func (c DMChannel) LastPinTimestamp() *time.Time {
	return c.lastPinTimestamp
}

func (c DMChannel) CreatedAt() time.Time {
	return c.id.Time()
}

func (DMChannel) channel()        {}
func (DMChannel) messageChannel() {}

var (
	_ Channel        = (*GroupDMChannel)(nil)
	_ MessageChannel = (*GroupDMChannel)(nil)
)

type GroupDMChannel struct {
	id               snowflake.ID
	ownerID          *snowflake.ID
	name             string
	lastPinTimestamp *time.Time
	lastMessageID    *snowflake.ID
	icon             *string
}

func (c *GroupDMChannel) UnmarshalJSON(data []byte) error {
	var v groupDMChannel
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	c.id = v.ID
	c.ownerID = v.OwnerID
	c.name = v.Name
	c.lastPinTimestamp = v.LastPinTimestamp
	c.lastMessageID = v.LastMessageID
	c.icon = v.Icon
	return nil
}

func (c GroupDMChannel) MarshalJSON() ([]byte, error) {
	return json.Marshal(groupDMChannel{
		ID:               c.id,
		Type:             c.Type(),
		OwnerID:          c.ownerID,
		Name:             c.name,
		LastPinTimestamp: c.lastPinTimestamp,
		LastMessageID:    c.lastMessageID,
		Icon:             c.icon,
	})
}

func (c GroupDMChannel) String() string {
	return channelString(c)
}

func (c GroupDMChannel) ID() snowflake.ID {
	return c.id
}

func (GroupDMChannel) Type() ChannelType {
	return ChannelTypeGroupDM
}

func (c GroupDMChannel) OwnerID() *snowflake.ID {
	return c.ownerID
}

func (c GroupDMChannel) Name() string {
	return c.name
}

func (c GroupDMChannel) LastPinTimestamp() *time.Time {
	return c.lastPinTimestamp
}

func (c GroupDMChannel) LastMessageID() *snowflake.ID {
	return c.lastMessageID
}

func (c GroupDMChannel) CreatedAt() time.Time {
	return c.id.Time()
}

// IconURL returns the icon URL of this group DM or nil if not set
func (c GroupDMChannel) IconURL(opts ...CDNOpt) *string {
	if c.icon == nil {
		return nil
	}
	url := formatAssetURL(ChannelIcon, opts, c.id, *c.icon)
	return &url
}

func (GroupDMChannel) channel()        {}
func (GroupDMChannel) messageChannel() {}

var (
	_ Channel             = (*GuildVoiceChannel)(nil)
	_ GuildChannel        = (*GuildVoiceChannel)(nil)
	_ GuildAudioChannel   = (*GuildVoiceChannel)(nil)
	_ GuildMessageChannel = (*GuildVoiceChannel)(nil)
)

type GuildVoiceChannel struct {
	id                   snowflake.ID
	guildID              snowflake.ID
	position             int
	permissionOverwrites []PermissionOverwrite
	name                 string
	bitrate              int
	UserLimit            int
	parentID             *snowflake.ID
	rtcRegion            string
	VideoQualityMode     VideoQualityMode
	lastMessageID        *snowflake.ID
	nsfw                 bool
	rateLimitPerUser     int
}

func (c *GuildVoiceChannel) UnmarshalJSON(data []byte) error {
	var v guildVoiceChannel
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	c.id = v.ID
	c.guildID = v.GuildID
	c.position = v.Position
	c.permissionOverwrites = v.PermissionOverwrites
	c.name = v.Name
	c.bitrate = v.Bitrate
	c.UserLimit = v.UserLimit
	c.parentID = v.ParentID
	c.rtcRegion = v.RTCRegion
	c.VideoQualityMode = v.VideoQualityMode
	c.lastMessageID = v.LastMessageID
	c.nsfw = v.NSFW
	c.rateLimitPerUser = v.RateLimitPerUser
	return nil
}

func (c GuildVoiceChannel) MarshalJSON() ([]byte, error) {
	return json.Marshal(guildVoiceChannel{
		ID:                   c.id,
		Type:                 c.Type(),
		GuildID:              c.guildID,
		Position:             c.position,
		PermissionOverwrites: c.permissionOverwrites,
		Name:                 c.name,
		Bitrate:              c.bitrate,
		UserLimit:            c.UserLimit,
		ParentID:             c.parentID,
		RTCRegion:            c.rtcRegion,
		VideoQualityMode:     c.VideoQualityMode,
		LastMessageID:        c.lastMessageID,
		NSFW:                 c.nsfw,
		RateLimitPerUser:     c.rateLimitPerUser,
	})
}

func (c GuildVoiceChannel) String() string {
	return channelString(c)
}

func (c GuildVoiceChannel) Mention() string {
	return ChannelMention(c.ID())
}

func (GuildVoiceChannel) Type() ChannelType {
	return ChannelTypeGuildVoice
}

func (c GuildVoiceChannel) ID() snowflake.ID {
	return c.id
}

func (c GuildVoiceChannel) Name() string {
	return c.name
}

func (c GuildVoiceChannel) GuildID() snowflake.ID {
	return c.guildID
}

func (c GuildVoiceChannel) PermissionOverwrites() PermissionOverwrites {
	return c.permissionOverwrites
}

func (c GuildVoiceChannel) Bitrate() int {
	return c.bitrate
}

func (c GuildVoiceChannel) RTCRegion() string {
	return c.rtcRegion
}

func (c GuildVoiceChannel) Position() int {
	return c.position
}

func (c GuildVoiceChannel) ParentID() *snowflake.ID {
	return c.parentID
}

func (c GuildVoiceChannel) LastMessageID() *snowflake.ID {
	return c.lastMessageID
}

// LastPinTimestamp always returns nil for GuildVoiceChannel(s) as they cannot have pinned messages.
func (c GuildVoiceChannel) LastPinTimestamp() *time.Time {
	return nil
}

// Topic always returns nil for GuildVoiceChannel(s) as they do not have their own topic.
func (c GuildVoiceChannel) Topic() *string {
	return nil
}

func (c GuildVoiceChannel) NSFW() bool {
	return c.nsfw
}

// DefaultAutoArchiveDuration is always 0 for GuildVoiceChannel(s) as they do not have their own AutoArchiveDuration.
func (c GuildVoiceChannel) DefaultAutoArchiveDuration() AutoArchiveDuration {
	return 0
}

func (c GuildVoiceChannel) RateLimitPerUser() int {
	return c.rateLimitPerUser
}

func (c GuildVoiceChannel) CreatedAt() time.Time {
	return c.id.Time()
}

func (GuildVoiceChannel) channel()             {}
func (GuildVoiceChannel) messageChannel()      {}
func (GuildVoiceChannel) guildChannel()        {}
func (GuildVoiceChannel) guildAudioChannel()   {}
func (GuildVoiceChannel) guildMessageChannel() {}

var (
	_ Channel      = (*GuildCategoryChannel)(nil)
	_ GuildChannel = (*GuildCategoryChannel)(nil)
)

type GuildCategoryChannel struct {
	id                   snowflake.ID
	guildID              snowflake.ID
	position             int
	permissionOverwrites PermissionOverwrites
	name                 string
}

func (c *GuildCategoryChannel) UnmarshalJSON(data []byte) error {
	var v guildCategoryChannel
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	c.id = v.ID
	c.guildID = v.GuildID
	c.position = v.Position
	c.permissionOverwrites = v.PermissionOverwrites
	c.name = v.Name
	return nil
}

func (c GuildCategoryChannel) MarshalJSON() ([]byte, error) {
	return json.Marshal(guildCategoryChannel{
		ID:                   c.id,
		Type:                 c.Type(),
		GuildID:              c.guildID,
		Position:             c.position,
		PermissionOverwrites: c.permissionOverwrites,
		Name:                 c.name,
	})
}

func (c GuildCategoryChannel) String() string {
	return channelString(c)
}

func (c GuildCategoryChannel) Mention() string {
	return ChannelMention(c.ID())
}

func (GuildCategoryChannel) Type() ChannelType {
	return ChannelTypeGuildCategory
}

func (c GuildCategoryChannel) ID() snowflake.ID {
	return c.id
}

func (c GuildCategoryChannel) Name() string {
	return c.name
}

func (c GuildCategoryChannel) GuildID() snowflake.ID {
	return c.guildID
}

func (c GuildCategoryChannel) PermissionOverwrites() PermissionOverwrites {
	return c.permissionOverwrites
}

func (c GuildCategoryChannel) Position() int {
	return c.position
}

// ParentID always returns nil for GuildCategoryChannel as they can't be nested.
func (c GuildCategoryChannel) ParentID() *snowflake.ID {
	return nil
}

func (c GuildCategoryChannel) CreatedAt() time.Time {
	return c.id.Time()
}

func (GuildCategoryChannel) channel()      {}
func (GuildCategoryChannel) guildChannel() {}

var (
	_ Channel             = (*GuildNewsChannel)(nil)
	_ GuildChannel        = (*GuildNewsChannel)(nil)
	_ MessageChannel      = (*GuildNewsChannel)(nil)
	_ GuildMessageChannel = (*GuildNewsChannel)(nil)
)

type GuildNewsChannel struct {
	id                         snowflake.ID
	guildID                    snowflake.ID
	position                   int
	permissionOverwrites       PermissionOverwrites
	name                       string
	topic                      *string
	nsfw                       bool
	lastMessageID              *snowflake.ID
	rateLimitPerUser           int
	parentID                   *snowflake.ID
	lastPinTimestamp           *time.Time
	defaultAutoArchiveDuration AutoArchiveDuration
}

func (c *GuildNewsChannel) UnmarshalJSON(data []byte) error {
	var v guildNewsChannel
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	c.id = v.ID
	c.guildID = v.GuildID
	c.position = v.Position
	c.permissionOverwrites = v.PermissionOverwrites
	c.name = v.Name
	c.topic = v.Topic
	c.nsfw = v.NSFW
	c.lastMessageID = v.LastMessageID
	c.rateLimitPerUser = v.RateLimitPerUser
	c.parentID = v.ParentID
	c.lastPinTimestamp = v.LastPinTimestamp
	c.defaultAutoArchiveDuration = v.DefaultAutoArchiveDuration
	return nil
}

func (c GuildNewsChannel) MarshalJSON() ([]byte, error) {
	return json.Marshal(guildNewsChannel{
		ID:                         c.id,
		Type:                       c.Type(),
		GuildID:                    c.guildID,
		Position:                   c.position,
		PermissionOverwrites:       c.permissionOverwrites,
		Name:                       c.name,
		Topic:                      c.topic,
		NSFW:                       c.nsfw,
		LastMessageID:              c.lastMessageID,
		RateLimitPerUser:           c.rateLimitPerUser,
		ParentID:                   c.parentID,
		LastPinTimestamp:           c.lastPinTimestamp,
		DefaultAutoArchiveDuration: c.defaultAutoArchiveDuration,
	})
}

func (c GuildNewsChannel) String() string {
	return channelString(c)
}

func (c GuildNewsChannel) Mention() string {
	return ChannelMention(c.ID())
}

func (GuildNewsChannel) Type() ChannelType {
	return ChannelTypeGuildNews
}

func (c GuildNewsChannel) ID() snowflake.ID {
	return c.id
}

func (c GuildNewsChannel) Name() string {
	return c.name
}

func (c GuildNewsChannel) GuildID() snowflake.ID {
	return c.guildID
}

func (c GuildNewsChannel) PermissionOverwrites() PermissionOverwrites {
	return c.permissionOverwrites
}

func (c GuildNewsChannel) Topic() *string {
	return c.topic
}

func (c GuildNewsChannel) NSFW() bool {
	return c.nsfw
}

func (c GuildNewsChannel) DefaultAutoArchiveDuration() AutoArchiveDuration {
	return c.defaultAutoArchiveDuration
}

func (c GuildNewsChannel) LastMessageID() *snowflake.ID {
	return c.lastMessageID
}

func (c GuildNewsChannel) RateLimitPerUser() int {
	return c.rateLimitPerUser
}

func (c GuildNewsChannel) LastPinTimestamp() *time.Time {
	return c.lastPinTimestamp
}

func (c GuildNewsChannel) Position() int {
	return c.position
}

func (c GuildNewsChannel) ParentID() *snowflake.ID {
	return c.parentID
}

func (c GuildNewsChannel) CreatedAt() time.Time {
	return c.id.Time()
}

func (GuildNewsChannel) channel()             {}
func (GuildNewsChannel) guildChannel()        {}
func (GuildNewsChannel) messageChannel()      {}
func (GuildNewsChannel) guildMessageChannel() {}

var (
	_ Channel             = (*GuildThread)(nil)
	_ GuildChannel        = (*GuildThread)(nil)
	_ MessageChannel      = (*GuildThread)(nil)
	_ GuildMessageChannel = (*GuildThread)(nil)
)

type GuildThread struct {
	id               snowflake.ID
	channelType      ChannelType
	guildID          snowflake.ID
	name             string
	nsfw             bool
	lastMessageID    *snowflake.ID
	lastPinTimestamp *time.Time
	rateLimitPerUser int
	OwnerID          snowflake.ID
	parentID         snowflake.ID
	MessageCount     int
	TotalMessageSent int
	AppliedTags      []snowflake.ID
	MemberCount      int
	ThreadMetadata   ThreadMetadata
}

func (c *GuildThread) UnmarshalJSON(data []byte) error {
	var v guildThread
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	c.id = v.ID
	c.channelType = v.Type
	c.guildID = v.GuildID
	c.name = v.Name
	c.nsfw = v.NSFW
	c.lastMessageID = v.LastMessageID
	c.lastPinTimestamp = v.LastPinTimestamp
	c.rateLimitPerUser = v.RateLimitPerUser
	c.OwnerID = v.OwnerID
	c.parentID = v.ParentID
	c.MessageCount = v.MessageCount
	c.TotalMessageSent = v.TotalMessageSent
	c.AppliedTags = v.AppliedTags
	c.MemberCount = v.MemberCount
	c.ThreadMetadata = v.ThreadMetadata
	return nil
}

func (c GuildThread) MarshalJSON() ([]byte, error) {
	return json.Marshal(guildThread{
		ID:               c.id,
		Type:             c.channelType,
		GuildID:          c.guildID,
		Name:             c.name,
		NSFW:             c.nsfw,
		LastMessageID:    c.lastMessageID,
		LastPinTimestamp: c.lastPinTimestamp,
		RateLimitPerUser: c.rateLimitPerUser,
		OwnerID:          c.OwnerID,
		ParentID:         c.parentID,
		MessageCount:     c.MessageCount,
		TotalMessageSent: c.TotalMessageSent,
		AppliedTags:      c.AppliedTags,
		MemberCount:      c.MemberCount,
		ThreadMetadata:   c.ThreadMetadata,
	})
}

func (c GuildThread) String() string {
	return channelString(c)
}

func (c GuildThread) Mention() string {
	return ChannelMention(c.ID())
}

func (c GuildThread) Type() ChannelType {
	return c.channelType
}

func (c GuildThread) ID() snowflake.ID {
	return c.id
}

// PermissionOverwrites always returns nil for GuildThread(s) as they do not have their own PermissionOverwrites.
func (c GuildThread) PermissionOverwrites() PermissionOverwrites {
	return nil
}

// Topic always returns nil for GuildThread(s) as they do not have their own topic.
func (c GuildThread) Topic() *string {
	return nil
}

func (c GuildThread) NSFW() bool {
	return c.nsfw
}

func (c GuildThread) Name() string {
	return c.name
}

func (c GuildThread) GuildID() snowflake.ID {
	return c.guildID
}

func (c GuildThread) LastMessageID() *snowflake.ID {
	return c.lastMessageID
}

func (c GuildThread) RateLimitPerUser() int {
	return 0
}

func (c GuildThread) LastPinTimestamp() *time.Time {
	return c.lastPinTimestamp
}

// Position always returns 0 for GuildThread(s) as they do not have their own position.
func (c GuildThread) Position() int {
	return 0
}

// ParentID is never nil for GuildThread(s).
func (c GuildThread) ParentID() *snowflake.ID {
	return &c.parentID
}

// DefaultAutoArchiveDuration is always 0 for GuildThread(s) as they do not have their own AutoArchiveDuration.
func (c GuildThread) DefaultAutoArchiveDuration() AutoArchiveDuration {
	return 0
}

func (c GuildThread) CreatedAt() time.Time {
	return c.id.Time()
}

func (GuildThread) channel()             {}
func (GuildThread) guildChannel()        {}
func (GuildThread) messageChannel()      {}
func (GuildThread) guildMessageChannel() {}

var (
	_ Channel             = (*GuildStageVoiceChannel)(nil)
	_ GuildChannel        = (*GuildStageVoiceChannel)(nil)
	_ GuildAudioChannel   = (*GuildStageVoiceChannel)(nil)
	_ GuildMessageChannel = (*GuildStageVoiceChannel)(nil)
)

type GuildStageVoiceChannel struct {
	id                   snowflake.ID
	guildID              snowflake.ID
	position             int
	permissionOverwrites PermissionOverwrites
	name                 string
	bitrate              int
	parentID             *snowflake.ID
	rtcRegion            string
	VideoQualityMode     VideoQualityMode
	lastMessageID        *snowflake.ID
	nsfw                 bool
	rateLimitPerUser     int
}

func (c *GuildStageVoiceChannel) UnmarshalJSON(data []byte) error {
	var v guildStageVoiceChannel
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	c.id = v.ID
	c.guildID = v.GuildID
	c.position = v.Position
	c.permissionOverwrites = v.PermissionOverwrites
	c.name = v.Name
	c.bitrate = v.Bitrate
	c.parentID = v.ParentID
	c.rtcRegion = v.RTCRegion
	c.VideoQualityMode = v.VideoQualityMode
	c.lastMessageID = v.LastMessageID
	c.nsfw = v.NSFW
	c.rateLimitPerUser = v.RateLimitPerUser
	return nil
}

func (c GuildStageVoiceChannel) MarshalJSON() ([]byte, error) {
	return json.Marshal(guildStageVoiceChannel{
		ID:                   c.id,
		Type:                 c.Type(),
		GuildID:              c.guildID,
		Position:             c.position,
		PermissionOverwrites: c.permissionOverwrites,
		Name:                 c.name,
		Bitrate:              c.bitrate,
		ParentID:             c.parentID,
		RTCRegion:            c.rtcRegion,
		VideoQualityMode:     c.VideoQualityMode,
		LastMessageID:        c.lastMessageID,
		NSFW:                 c.nsfw,
		RateLimitPerUser:     c.rateLimitPerUser,
	})
}

func (c GuildStageVoiceChannel) String() string {
	return channelString(c)
}

func (c GuildStageVoiceChannel) Mention() string {
	return ChannelMention(c.ID())
}

func (GuildStageVoiceChannel) Type() ChannelType {
	return ChannelTypeGuildStageVoice
}

func (c GuildStageVoiceChannel) ID() snowflake.ID {
	return c.id
}

func (c GuildStageVoiceChannel) Name() string {
	return c.name
}

func (c GuildStageVoiceChannel) GuildID() snowflake.ID {
	return c.guildID
}

func (c GuildStageVoiceChannel) PermissionOverwrites() PermissionOverwrites {
	return c.permissionOverwrites
}

func (c GuildStageVoiceChannel) Bitrate() int {
	return c.bitrate
}

func (c GuildStageVoiceChannel) RTCRegion() string {
	return c.rtcRegion
}

func (c GuildStageVoiceChannel) Position() int {
	return c.position
}

func (c GuildStageVoiceChannel) ParentID() *snowflake.ID {
	return c.parentID
}

func (c GuildStageVoiceChannel) LastMessageID() *snowflake.ID {
	return c.lastMessageID
}

// LastPinTimestamp always returns nil for GuildStageVoiceChannel(s) as they cannot have pinned messages.
func (c GuildStageVoiceChannel) LastPinTimestamp() *time.Time {
	return nil
}

// Topic always returns nil for GuildStageVoiceChannel(s) as they do not have their own topic.
func (c GuildStageVoiceChannel) Topic() *string {
	return nil
}

func (c GuildStageVoiceChannel) NSFW() bool {
	return c.nsfw
}

// DefaultAutoArchiveDuration is always 0 for GuildStageVoiceChannel(s) as they do not have their own AutoArchiveDuration.
func (c GuildStageVoiceChannel) DefaultAutoArchiveDuration() AutoArchiveDuration {
	return 0
}

func (c GuildStageVoiceChannel) RateLimitPerUser() int {
	return c.rateLimitPerUser
}

func (c GuildStageVoiceChannel) CreatedAt() time.Time {
	return c.id.Time()
}

func (GuildStageVoiceChannel) channel()             {}
func (GuildStageVoiceChannel) messageChannel()      {}
func (GuildStageVoiceChannel) guildChannel()        {}
func (GuildStageVoiceChannel) guildAudioChannel()   {}
func (GuildStageVoiceChannel) guildMessageChannel() {}

var (
	_ Channel      = (*GuildForumChannel)(nil)
	_ GuildChannel = (*GuildForumChannel)(nil)
)

type GuildForumChannel struct {
	id                            snowflake.ID
	guildID                       snowflake.ID
	position                      int
	permissionOverwrites          PermissionOverwrites
	name                          string
	parentID                      *snowflake.ID
	LastPostID                    *snowflake.ID
	Topic                         *string
	NSFW                          bool
	RateLimitPerUser              int
	Flags                         ChannelFlags
	AvailableTags                 []ChannelTag
	DefaultReactionEmoji          *DefaultReactionEmoji
	DefaultThreadRateLimitPerUser int
	DefaultSortOrder              *DefaultSortOrder
	DefaultForumLayout            DefaultForumLayout
}

func (c *GuildForumChannel) UnmarshalJSON(data []byte) error {
	var v guildForumChannel
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	c.id = v.ID
	c.guildID = v.GuildID
	c.position = v.Position
	c.permissionOverwrites = v.PermissionOverwrites
	c.name = v.Name
	c.parentID = v.ParentID
	c.LastPostID = v.LastPostID
	c.Topic = v.Topic
	c.NSFW = v.NSFW
	c.RateLimitPerUser = v.RateLimitPerUser
	c.Flags = v.Flags
	c.AvailableTags = v.AvailableTags
	c.DefaultReactionEmoji = v.DefaultReactionEmoji
	c.DefaultThreadRateLimitPerUser = v.DefaultThreadRateLimitPerUser
	c.DefaultSortOrder = v.DefaultSortOrder
	c.DefaultForumLayout = v.DefaultForumLayout
	return nil
}

func (c GuildForumChannel) MarshalJSON() ([]byte, error) {
	return json.Marshal(guildForumChannel{
		ID:                            c.id,
		Type:                          c.Type(),
		GuildID:                       c.guildID,
		Position:                      c.position,
		PermissionOverwrites:          c.permissionOverwrites,
		Name:                          c.name,
		ParentID:                      c.parentID,
		LastPostID:                    c.LastPostID,
		Topic:                         c.Topic,
		NSFW:                          c.NSFW,
		RateLimitPerUser:              c.RateLimitPerUser,
		Flags:                         c.Flags,
		AvailableTags:                 c.AvailableTags,
		DefaultReactionEmoji:          c.DefaultReactionEmoji,
		DefaultThreadRateLimitPerUser: c.DefaultThreadRateLimitPerUser,
		DefaultSortOrder:              c.DefaultSortOrder,
		DefaultForumLayout:            c.DefaultForumLayout,
	})
}

func (c GuildForumChannel) String() string {
	return channelString(c)
}

func (c GuildForumChannel) Mention() string {
	return ChannelMention(c.ID())
}

func (GuildForumChannel) Type() ChannelType {
	return ChannelTypeGuildForum
}

func (c GuildForumChannel) ID() snowflake.ID {
	return c.id
}

func (c GuildForumChannel) Name() string {
	return c.name
}

func (c GuildForumChannel) GuildID() snowflake.ID {
	return c.guildID
}

func (c GuildForumChannel) PermissionOverwrites() PermissionOverwrites {
	return c.permissionOverwrites
}

func (c GuildForumChannel) Position() int {
	return c.position
}

func (c GuildForumChannel) ParentID() *snowflake.ID {
	return c.parentID
}

func (c GuildForumChannel) CreatedAt() time.Time {
	return c.id.Time()
}

func (GuildForumChannel) channel()      {}
func (GuildForumChannel) guildChannel() {}

var (
	_ Channel      = (*GuildMediaChannel)(nil)
	_ GuildChannel = (*GuildMediaChannel)(nil)
)

type GuildMediaChannel struct {
	id                            snowflake.ID
	guildID                       snowflake.ID
	position                      int
	permissionOverwrites          PermissionOverwrites
	name                          string
	parentID                      *snowflake.ID
	LastPostID                    *snowflake.ID
	Topic                         *string
	NSFW                          bool
	RateLimitPerUser              int
	Flags                         ChannelFlags
	AvailableTags                 []ChannelTag
	DefaultReactionEmoji          *DefaultReactionEmoji
	DefaultThreadRateLimitPerUser int
	DefaultSortOrder              *DefaultSortOrder
}

func (c *GuildMediaChannel) UnmarshalJSON(data []byte) error {
	var v guildMediaChannel
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	c.id = v.ID
	c.guildID = v.GuildID
	c.position = v.Position
	c.permissionOverwrites = v.PermissionOverwrites
	c.name = v.Name
	c.parentID = v.ParentID
	c.LastPostID = v.LastPostID
	c.Topic = v.Topic
	c.NSFW = v.NSFW
	c.RateLimitPerUser = v.RateLimitPerUser
	c.Flags = v.Flags
	c.AvailableTags = v.AvailableTags
	c.DefaultReactionEmoji = v.DefaultReactionEmoji
	c.DefaultThreadRateLimitPerUser = v.DefaultThreadRateLimitPerUser
	c.DefaultSortOrder = v.DefaultSortOrder
	return nil
}

func (c GuildMediaChannel) MarshalJSON() ([]byte, error) {
	return json.Marshal(guildMediaChannel{
		ID:                            c.id,
		Type:                          c.Type(),
		GuildID:                       c.guildID,
		Position:                      c.position,
		PermissionOverwrites:          c.permissionOverwrites,
		Name:                          c.name,
		ParentID:                      c.parentID,
		LastPostID:                    c.LastPostID,
		Topic:                         c.Topic,
		NSFW:                          c.NSFW,
		RateLimitPerUser:              c.RateLimitPerUser,
		Flags:                         c.Flags,
		AvailableTags:                 c.AvailableTags,
		DefaultReactionEmoji:          c.DefaultReactionEmoji,
		DefaultThreadRateLimitPerUser: c.DefaultThreadRateLimitPerUser,
		DefaultSortOrder:              c.DefaultSortOrder,
	})
}

func (c GuildMediaChannel) String() string {
	return channelString(c)
}

func (c GuildMediaChannel) Mention() string {
	return ChannelMention(c.ID())
}

func (GuildMediaChannel) Type() ChannelType {
	return ChannelTypeGuildMedia
}

func (c GuildMediaChannel) ID() snowflake.ID {
	return c.id
}

func (c GuildMediaChannel) Name() string {
	return c.name
}

func (c GuildMediaChannel) GuildID() snowflake.ID {
	return c.guildID
}

func (c GuildMediaChannel) PermissionOverwrites() PermissionOverwrites {
	return c.permissionOverwrites
}

func (c GuildMediaChannel) Position() int {
	return c.position
}

func (c GuildMediaChannel) ParentID() *snowflake.ID {
	return c.parentID
}

func (c GuildMediaChannel) CreatedAt() time.Time {
	return c.id.Time()
}

func (GuildMediaChannel) channel()      {}
func (GuildMediaChannel) guildChannel() {}

type FollowedChannel struct {
	ChannelID snowflake.ID `json:"channel_id"`
	WebhookID snowflake.ID `json:"webhook_id"`
}

type FollowChannel struct {
	ChannelID snowflake.ID `json:"webhook_channel_id"`
}

type PartialChannel struct {
	ID   snowflake.ID `json:"id"`
	Type ChannelType  `json:"type"`
}

// VideoQualityMode https://com/developers/docs/resources/channel#channel-object-video-quality-modes
type VideoQualityMode int

const (
	VideoQualityModeAuto VideoQualityMode = iota + 1
	VideoQualityModeFull
)

type ThreadMetadata struct {
	Archived            bool                `json:"archived"`
	AutoArchiveDuration AutoArchiveDuration `json:"auto_archive_duration"`
	ArchiveTimestamp    time.Time           `json:"archive_timestamp"`
	Locked              bool                `json:"locked"`
	Invitable           bool                `json:"invitable"`
	CreateTimestamp     time.Time           `json:"create_timestamp"`
}

type ChannelTag struct {
	ID        snowflake.ID  `json:"id"`
	Name      string        `json:"name"`
	Moderated bool          `json:"moderated"`
	EmojiID   *snowflake.ID `json:"emoji_id"`
	EmojiName *string       `json:"emoji_name"`
}

type DefaultReactionEmoji struct {
	EmojiID   *snowflake.ID `json:"emoji_id"`
	EmojiName *string       `json:"emoji_name"`
}

type DefaultSortOrder int

const (
	DefaultSortOrderLatestActivity DefaultSortOrder = iota
	DefaultSortOrderCreationDate
)

type DefaultForumLayout int

const (
	DefaultForumLayoutNotSet DefaultForumLayout = iota
	DefaultForumLayoutListView
	DefaultForumLayoutGalleryView
)

type AutoArchiveDuration int

const (
	AutoArchiveDuration1h  AutoArchiveDuration = 60
	AutoArchiveDuration24h AutoArchiveDuration = 1440
	AutoArchiveDuration3d  AutoArchiveDuration = 4320
	AutoArchiveDuration1w  AutoArchiveDuration = 10080
)

func channelString(channel Channel) string {
	return fmt.Sprintf("%d:%s(%s)", channel.Type(), channel.Name(), channel.ID())
}

func ApplyGuildIDToThread(guildThread GuildThread, guildID snowflake.ID) GuildThread {
	guildThread.guildID = guildID
	return guildThread
}

func ApplyGuildIDToChannel(channel GuildChannel, guildID snowflake.ID) GuildChannel {
	switch c := channel.(type) {
	case GuildTextChannel:
		c.guildID = guildID
		return c
	case GuildVoiceChannel:
		c.guildID = guildID
		return c
	case GuildCategoryChannel:
		c.guildID = guildID
		return c
	case GuildNewsChannel:
		c.guildID = guildID
		return c
	case GuildStageVoiceChannel:
		c.guildID = guildID
		return c
	case GuildThread:
		c.guildID = guildID
		return c
	case GuildForumChannel:
		c.guildID = guildID
		return c
	case GuildMediaChannel:
		c.guildID = guildID
		return c
	default:
		return channel
	}
}

func ApplyLastMessageIDToChannel(channel GuildMessageChannel, lastMessageID snowflake.ID) GuildMessageChannel {
	switch c := channel.(type) {
	case GuildTextChannel:
		c.lastMessageID = &lastMessageID
		return c
	case GuildVoiceChannel:
		c.lastMessageID = &lastMessageID
		return c
	case GuildNewsChannel:
		c.lastMessageID = &lastMessageID
		return c
	case GuildThread:
		c.lastMessageID = &lastMessageID
		return c
	default:
		return channel
	}
}

func ApplyLastPinTimestampToChannel(channel GuildMessageChannel, lastPinTimestamp *time.Time) GuildMessageChannel {
	switch c := channel.(type) {
	case GuildTextChannel:
		c.lastPinTimestamp = lastPinTimestamp
		return c
	case GuildNewsChannel:
		c.lastPinTimestamp = lastPinTimestamp
		return c
	case GuildThread:
		c.lastPinTimestamp = lastPinTimestamp
		return c
	default:
		return channel
	}
}
