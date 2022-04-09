package discord

import (
	"fmt"

	"github.com/disgoorg/disgo/json"
	"github.com/disgoorg/snowflake"
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
)

type Channel interface {
	json.Marshaler
	fmt.Stringer

	Type() ChannelType
	ID() snowflake.Snowflake
	Name() string

	channel()
}

type MessageChannel interface {
	Channel

	LastMessageID() *snowflake.Snowflake
	LastPinTimestamp() *Time

	messageChannel()
}

type GuildChannel interface {
	Channel
	Mentionable

	GuildID() snowflake.Snowflake
	Position() int
	ParentID() *snowflake.Snowflake
	PermissionOverwrites() []PermissionOverwrite

	guildChannel()
}

type GuildMessageChannel interface {
	GuildChannel
	MessageChannel

	Topic() *string
	NSFW() bool
	DefaultAutoArchiveDuration() AutoArchiveDuration

	guildMessageChannel()
}

type GuildAudioChannel interface {
	GuildChannel

	Bitrate() int
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
	_ Channel             = (*GuildTextChannel)(nil)
	_ GuildChannel        = (*GuildTextChannel)(nil)
	_ MessageChannel      = (*GuildTextChannel)(nil)
	_ GuildMessageChannel = (*GuildTextChannel)(nil)
)

type GuildTextChannel struct {
	id                         snowflake.Snowflake
	guildID                    snowflake.Snowflake
	position                   int
	permissionOverwrites       []PermissionOverwrite
	name                       string
	topic                      *string
	nsfw                       bool
	lastMessageID              *snowflake.Snowflake
	rateLimitPerUser           int
	parentID                   *snowflake.Snowflake
	lastPinTimestamp           *Time
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

func (c GuildTextChannel) ID() snowflake.Snowflake {
	return c.id
}

func (GuildTextChannel) Type() ChannelType {
	return ChannelTypeGuildText
}

func (c GuildTextChannel) Name() string {
	return c.name
}

func (c GuildTextChannel) GuildID() snowflake.Snowflake {
	return c.guildID
}

func (c GuildTextChannel) PermissionOverwrites() []PermissionOverwrite {
	return c.permissionOverwrites
}

func (c GuildTextChannel) Position() int {
	return c.position
}

func (c GuildTextChannel) ParentID() *snowflake.Snowflake {
	return c.parentID
}

func (c GuildTextChannel) LastMessageID() *snowflake.Snowflake {
	return c.lastMessageID
}

func (c GuildTextChannel) LastPinTimestamp() *Time {
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

func (GuildTextChannel) channel()             {}
func (GuildTextChannel) guildChannel()        {}
func (GuildTextChannel) messageChannel()      {}
func (GuildTextChannel) guildMessageChannel() {}

var (
	_ Channel        = (*DMChannel)(nil)
	_ MessageChannel = (*DMChannel)(nil)
)

type DMChannel struct {
	id               snowflake.Snowflake
	lastMessageID    *snowflake.Snowflake
	recipients       []User
	lastPinTimestamp *Time
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

func (c DMChannel) ID() snowflake.Snowflake {
	return c.id
}

func (DMChannel) Type() ChannelType {
	return ChannelTypeDM
}

func (c DMChannel) Name() string {
	return c.recipients[0].Username
}

func (c DMChannel) LastMessageID() *snowflake.Snowflake {
	return c.lastMessageID
}

func (c DMChannel) LastPinTimestamp() *Time {
	return c.lastPinTimestamp
}

func (DMChannel) channel()        {}
func (DMChannel) messageChannel() {}

var (
	_ Channel             = (*GuildVoiceChannel)(nil)
	_ GuildChannel        = (*GuildVoiceChannel)(nil)
	_ GuildAudioChannel   = (*GuildVoiceChannel)(nil)
	_ GuildMessageChannel = (*GuildVoiceChannel)(nil)
)

type GuildVoiceChannel struct {
	id                         snowflake.Snowflake
	guildID                    snowflake.Snowflake
	position                   int
	permissionOverwrites       []PermissionOverwrite
	name                       string
	bitrate                    int
	UserLimit                  int
	parentID                   *snowflake.Snowflake
	rtcRegion                  string
	VideoQualityMode           VideoQualityMode
	lastMessageID              *snowflake.Snowflake
	lastPinTimestamp           *Time
	topic                      *string
	nsfw                       bool
	defaultAutoArchiveDuration AutoArchiveDuration
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
	c.lastPinTimestamp = v.LastPinTimestamp
	c.topic = v.Topic
	c.nsfw = v.NSFW
	c.defaultAutoArchiveDuration = v.DefaultAutoArchiveDuration
	return nil
}

func (c GuildVoiceChannel) MarshalJSON() ([]byte, error) {
	return json.Marshal(guildVoiceChannel{
		ID:                         c.id,
		Type:                       c.Type(),
		GuildID:                    c.guildID,
		Position:                   c.position,
		PermissionOverwrites:       c.permissionOverwrites,
		Name:                       c.name,
		Bitrate:                    c.bitrate,
		UserLimit:                  c.UserLimit,
		ParentID:                   c.parentID,
		RTCRegion:                  c.rtcRegion,
		VideoQualityMode:           c.VideoQualityMode,
		LastMessageID:              c.lastMessageID,
		LastPinTimestamp:           c.lastPinTimestamp,
		Topic:                      c.topic,
		NSFW:                       c.nsfw,
		DefaultAutoArchiveDuration: c.defaultAutoArchiveDuration,
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

func (c GuildVoiceChannel) ID() snowflake.Snowflake {
	return c.id
}

func (c GuildVoiceChannel) Name() string {
	return c.name
}

func (c GuildVoiceChannel) GuildID() snowflake.Snowflake {
	return c.guildID
}

func (c GuildVoiceChannel) PermissionOverwrites() []PermissionOverwrite {
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

func (c GuildVoiceChannel) ParentID() *snowflake.Snowflake {
	return c.parentID
}

func (c GuildVoiceChannel) LastMessageID() *snowflake.Snowflake {
	return c.lastMessageID
}

func (c GuildVoiceChannel) LastPinTimestamp() *Time {
	return c.lastPinTimestamp
}

func (c GuildVoiceChannel) Topic() *string {
	return c.topic
}

func (c GuildVoiceChannel) NSFW() bool {
	return c.nsfw
}

func (c GuildVoiceChannel) DefaultAutoArchiveDuration() AutoArchiveDuration {
	return c.defaultAutoArchiveDuration
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
	id                   snowflake.Snowflake
	guildID              snowflake.Snowflake
	position             int
	permissionOverwrites []PermissionOverwrite
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

func (c GuildCategoryChannel) ID() snowflake.Snowflake {
	return c.id
}

func (c GuildCategoryChannel) Name() string {
	return c.name
}

func (c GuildCategoryChannel) GuildID() snowflake.Snowflake {
	return c.guildID
}

func (c GuildCategoryChannel) PermissionOverwrites() []PermissionOverwrite {
	return c.permissionOverwrites
}

func (c GuildCategoryChannel) Position() int {
	return c.position
}

func (c GuildCategoryChannel) ParentID() *snowflake.Snowflake {
	return nil
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
	id                         snowflake.Snowflake
	guildID                    snowflake.Snowflake
	position                   int
	permissionOverwrites       []PermissionOverwrite
	name                       string
	topic                      *string
	nsfw                       bool
	lastMessageID              *snowflake.Snowflake
	rateLimitPerUser           int
	parentID                   *snowflake.Snowflake
	lastPinTimestamp           *Time
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

func (c GuildNewsChannel) ID() snowflake.Snowflake {
	return c.id
}

func (c GuildNewsChannel) Name() string {
	return c.name
}

func (c GuildNewsChannel) GuildID() snowflake.Snowflake {
	return c.guildID
}

func (c GuildNewsChannel) PermissionOverwrites() []PermissionOverwrite {
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

func (c GuildNewsChannel) LastMessageID() *snowflake.Snowflake {
	return c.lastMessageID
}

func (c GuildNewsChannel) LastPinTimestamp() *Time {
	return c.lastPinTimestamp
}

func (c GuildNewsChannel) Position() int {
	return c.position
}

func (c GuildNewsChannel) ParentID() *snowflake.Snowflake {
	return c.parentID
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
	id               snowflake.Snowflake
	channelType      ChannelType
	guildID          snowflake.Snowflake
	name             string
	nsfw             bool
	lastMessageID    *snowflake.Snowflake
	lastPinTimestamp *Time
	RateLimitPerUser int
	OwnerID          snowflake.Snowflake
	parentID         snowflake.Snowflake
	MessageCount     int
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
	c.RateLimitPerUser = v.RateLimitPerUser
	c.OwnerID = v.OwnerID
	c.parentID = v.ParentID
	c.MessageCount = v.MessageCount
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
		RateLimitPerUser: c.RateLimitPerUser,
		OwnerID:          c.OwnerID,
		ParentID:         c.parentID,
		MessageCount:     c.MessageCount,
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

func (c GuildThread) ID() snowflake.Snowflake {
	return c.id
}

func (c GuildThread) PermissionOverwrites() []PermissionOverwrite {
	return nil
}

func (c GuildThread) Topic() *string {
	return nil
}

func (c GuildThread) NSFW() bool {
	return c.nsfw
}

func (c GuildThread) Name() string {
	return c.name
}

func (c GuildThread) GuildID() snowflake.Snowflake {
	return c.guildID
}

func (c GuildThread) LastMessageID() *snowflake.Snowflake {
	return c.lastMessageID
}

func (c GuildThread) LastPinTimestamp() *Time {
	return c.lastPinTimestamp
}

func (c GuildThread) Position() int {
	return 0
}

func (c GuildThread) ParentID() *snowflake.Snowflake {
	return &c.parentID
}

func (c GuildThread) DefaultAutoArchiveDuration() AutoArchiveDuration {
	return 0
}

func (GuildThread) channel()             {}
func (GuildThread) guildChannel()        {}
func (GuildThread) messageChannel()      {}
func (GuildThread) guildMessageChannel() {}

var (
	_ Channel           = (*GuildStageVoiceChannel)(nil)
	_ GuildChannel      = (*GuildStageVoiceChannel)(nil)
	_ GuildAudioChannel = (*GuildStageVoiceChannel)(nil)
)

type GuildStageVoiceChannel struct {
	id                   snowflake.Snowflake
	guildID              snowflake.Snowflake
	position             int
	permissionOverwrites []PermissionOverwrite
	name                 string
	bitrate              int
	parentID             *snowflake.Snowflake
	rtcRegion            string
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

func (c GuildStageVoiceChannel) ID() snowflake.Snowflake {
	return c.id
}

func (c GuildStageVoiceChannel) Name() string {
	return c.name
}

func (c GuildStageVoiceChannel) GuildID() snowflake.Snowflake {
	return c.guildID
}

func (c GuildStageVoiceChannel) PermissionOverwrites() []PermissionOverwrite {
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

func (c GuildStageVoiceChannel) ParentID() *snowflake.Snowflake {
	return c.parentID
}

func (GuildStageVoiceChannel) channel()           {}
func (GuildStageVoiceChannel) guildChannel()      {}
func (GuildStageVoiceChannel) guildAudioChannel() {}

// VideoQualityMode https://com/developers/docs/resources/channel#channel-object-video-quality-modes
type VideoQualityMode int

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
	CreateTimestamp     Time                `json:"create_timestamp"`
}

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

func ApplyGuildIDToThread(guildThread GuildThread, guildID snowflake.Snowflake) GuildThread {
	guildThread.guildID = guildID
	return guildThread
}

func ApplyGuildIDToChannel(channel GuildChannel, guildID snowflake.Snowflake) GuildChannel {
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
	default:
		panic("unknown channel type")
	}
}
