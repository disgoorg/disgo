package core

import (
	"fmt"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type Channel interface {
	Disgo() Disgo
	ID() discord.Snowflake
	Name() string
	Type() discord.ChannelType

	IsMessageChannel() bool
	IsGuildChannel() bool
	IsTextChannel() bool
	IsVoiceChannel() bool
	IsDMChannel() bool
	IsCategory() bool
	IsNewsChannel() bool
	IsStoreChannel() bool
	IsStageChannel() bool
}

// MessageChannel is used for sending Message(s) to User(s)
type MessageChannel interface {
	Channel
	LastMessageID() *discord.Snowflake
	LastPinTimestamp() *discord.Time
	CreateMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*Message, rest.Error)
	UpdateMessage(messageID discord.Snowflake, messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*Message, rest.Error)
	DeleteMessage(messageID discord.Snowflake, opts ...rest.RequestOpt) rest.Error
	BulkDeleteMessages(messageIDs []discord.Snowflake, opts ...rest.RequestOpt) rest.Error

	//CollectMessages(filter collectors.MessageFilter) (chan *Message, func())
}

// DMChannel is used for interacting in private Message(s) with users
type DMChannel interface {
	MessageChannel
}

// GuildChannel is a generic type for all server channels
type GuildChannel interface {
	Channel
	Guild() *Guild
	GuildID() discord.Snowflake
	Permissions() discord.Permissions
	ParentID() *discord.Snowflake
	Parent() Category
	Position() *int
	Update(channelUpdate discord.ChannelUpdate, opts ...rest.RequestOpt) (GuildChannel, rest.Error)
}

// Category groups text & voice channels in servers together
type Category interface {
	GuildChannel
}

// VoiceChannel adds methods specifically for interacting with discord's voice
type VoiceChannel interface {
	GuildChannel
	Connect() error
	Bitrate() int
}

// TextChannel allows you to interact with discord's text channels
type TextChannel interface {
	GuildChannel
	MessageChannel
	NSFW() bool
	Topic() *string
}

// NewsChannel allows you to interact with discord's text channels
type NewsChannel interface {
	TextChannel
	CrosspostMessage(messageID discord.Snowflake, opts ...rest.RequestOpt) (*Message, rest.Error)
}

// StoreChannel allows you to interact with discord's store channels
type StoreChannel interface {
	GuildChannel
}

type StageChannel interface {
	VoiceChannel
	StageInstance() *StageInstance
	CreateStageInstance(stageInstanceCreate discord.StageInstanceCreate, opts ...rest.RequestOpt) (*StageInstance, rest.Error)
	UpdateStageInstance(stageInstanceUpdate discord.StageInstanceUpdate, opts ...rest.RequestOpt) (*StageInstance, rest.Error)
	DeleteStageInstance(opts ...rest.RequestOpt) rest.Error
	IsModerator(member *Member) bool
}

var _ Channel = (*ChannelImpl)(nil)

type ChannelImpl struct {
	discord.Channel
	disgo           Disgo
	stageInstanceID *discord.Snowflake
}

func (c *ChannelImpl) Guild() *Guild {
	return c.Disgo().Caches().GuildCache().Get(c.GuildID())
}

func (c *ChannelImpl) Disgo() Disgo {
	return c.disgo
}

func (c *ChannelImpl) ID() discord.Snowflake {
	return c.Channel.ID
}

func (c *ChannelImpl) Name() string {
	return *c.Channel.Name
}

func (c *ChannelImpl) Type() discord.ChannelType {
	return c.Channel.Type
}

func (c *ChannelImpl) IsMessageChannel() bool {
	return c.IsTextChannel() || c.IsNewsChannel() || c.IsDMChannel()
}

func (c *ChannelImpl) IsGuildChannel() bool {
	return c.IsCategory() || c.IsNewsChannel() || c.IsTextChannel() || c.IsVoiceChannel()
}

func (c *ChannelImpl) IsDMChannel() bool {
	return c.Type() != discord.ChannelTypeDM
}

func (c *ChannelImpl) IsTextChannel() bool {
	return c.Type() != discord.ChannelTypeText
}

func (c *ChannelImpl) IsVoiceChannel() bool {
	return c.Type() != discord.ChannelTypeVoice
}

func (c *ChannelImpl) IsCategory() bool {
	return c.Type() != discord.ChannelTypeCategory
}

func (c *ChannelImpl) IsNewsChannel() bool {
	return c.Type() != discord.ChannelTypeNews
}

func (c *ChannelImpl) IsStoreChannel() bool {
	return c.Type() != discord.ChannelTypeStore
}

func (c *ChannelImpl) IsStageChannel() bool {
	return c.Type() != discord.ChannelTypeStage
}

var _ MessageChannel = (*ChannelImpl)(nil)

func (c *ChannelImpl) LastMessageID() *discord.Snowflake {
	return c.Channel.LastMessageID
}

func (c *ChannelImpl) LastPinTimestamp() *discord.Time {
	return c.Channel.LastPinTimestamp
}

// CreateMessage sends a Message to a TextChannel
func (c *ChannelImpl) CreateMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*Message, rest.Error) {
	message, err := c.Disgo().RestServices().ChannelService().CreateMessage(c.ID(), messageCreate, opts...)
	if err != nil {
		return nil, err
	}
	return c.Disgo().EntityBuilder().CreateMessage(*message, CacheStrategyNoWs), nil
}

// UpdateMessage edits a Message in this TextChannel
func (c *ChannelImpl) UpdateMessage(messageID discord.Snowflake, messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*Message, rest.Error) {
	message, err := c.Disgo().RestServices().ChannelService().UpdateMessage(c.ID(), messageID, messageUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return c.Disgo().EntityBuilder().CreateMessage(*message, CacheStrategyNoWs), nil
}

// DeleteMessage allows you to edit an existing Message sent by you
func (c *ChannelImpl) DeleteMessage(messageID discord.Snowflake, opts ...rest.RequestOpt) rest.Error {
	return c.Disgo().RestServices().ChannelService().DeleteMessage(c.ID(), messageID, opts...)
}

// BulkDeleteMessages allows you bulk delete Message(s)
func (c *ChannelImpl) BulkDeleteMessages(messageIDs []discord.Snowflake, opts ...rest.RequestOpt) rest.Error {
	return c.Disgo().RestServices().ChannelService().BulkDeleteMessages(c.ID(), messageIDs, opts...)
}

/* func (c *channelImpl) CollectMessages(filter collectors.MessageFilter) (chan *Message, func()) {
	var guildID *discord.Snowflake = nil
	if c.IsTextChannel() {
		id := c.GuildID()
		guildID = &id
	}
	return collectors.NewMessageCollector(c.Disgo(), c.ID(), guildID, filter)
}*/

var _ DMChannel = (*ChannelImpl)(nil)

var _ GuildChannel = (*ChannelImpl)(nil)

// GuildID returns the channel's Guild ID
func (c *ChannelImpl) GuildID() discord.Snowflake {
	if !c.IsGuildChannel() || c.Channel.GuildID == nil {
		unsupported(c)
	}
	return *c.Channel.GuildID
}

func (c *ChannelImpl) Permissions() discord.Permissions {
	if !c.IsGuildChannel() {
		unsupported(c)
	}
	return *c.Channel.InteractionPermissions
}

func (c *ChannelImpl) ParentID() *discord.Snowflake {
	if !c.IsGuildChannel() {
		unsupported(c)
	}
	return c.Channel.ParentID
}

func (c *ChannelImpl) Parent() Category {
	if c.ParentID() == nil {
		return nil
	}
	return c.Disgo().Caches().CategoryCache().Get(*c.Channel.ParentID)
}

func (c *ChannelImpl) Position() *int {
	if !c.IsGuildChannel() {
		unsupported(c)
	}

	return c.Channel.Position
}

func (c *ChannelImpl) Update(channelUpdate discord.ChannelUpdate, opts ...rest.RequestOpt) (GuildChannel, rest.Error) {
	if !c.IsGuildChannel() {
		unsupported(c)
	}
	channel, err := c.Disgo().RestServices().ChannelService().UpdateChannel(c.ID(), channelUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return c.Disgo().EntityBuilder().CreateChannel(*channel, CacheStrategyNoWs).(GuildChannel), nil
}

var _ Category = (*ChannelImpl)(nil)
var _ VoiceChannel = (*ChannelImpl)(nil)

func (c *ChannelImpl) Connect() error {
	if !c.IsVoiceChannel() {
		unsupported(c)
	}
	return c.Disgo().AudioController().Connect(c.GuildID(), c.ID())
}

func (c *ChannelImpl) Bitrate() int {
	if !c.IsVoiceChannel() {
		unsupported(c)
	}
	return *c.Channel.Bitrate
}

var _ TextChannel = (*ChannelImpl)(nil)

func (c *ChannelImpl) NSFW() bool {
	if !c.IsTextChannel() {
		unsupported(c)
	}
	return *c.Channel.NSFW
}

func (c *ChannelImpl) Topic() *string {
	return c.Channel.Topic
}

var _ NewsChannel = (*ChannelImpl)(nil)

func (c *ChannelImpl) CrosspostMessage(messageID discord.Snowflake, opts ...rest.RequestOpt) (*Message, rest.Error) {
	message, err := c.Disgo().RestServices().ChannelService().CrosspostMessage(c.ID(), messageID, opts...)
	if err != nil {
		return nil, err
	}
	return c.Disgo().EntityBuilder().CreateMessage(*message, CacheStrategyNoWs), nil
}

var _ StoreChannel = (*ChannelImpl)(nil)

var _ StageChannel = (*ChannelImpl)(nil)

func (c *ChannelImpl) StageInstance() *StageInstance {
	if !c.IsStageChannel() {
		unsupported(c)
	}
	if c.stageInstanceID == nil {
		return nil
	}
	return c.Disgo().Caches().StageInstanceCache().Get(*c.stageInstanceID)
}

func (c *ChannelImpl) CreateStageInstance(stageInstanceCreate discord.StageInstanceCreate, opts ...rest.RequestOpt) (*StageInstance, rest.Error) {
	if !c.IsStageChannel() {
		unsupported(c)
	}
	stageInstance, err := c.Disgo().RestServices().StageInstanceService().CreateStageInstance(stageInstanceCreate, opts...)
	if err != nil {
		return nil, err
	}
	return c.Disgo().EntityBuilder().CreateStageInstance(*stageInstance, CacheStrategyNoWs), nil
}

func (c *ChannelImpl) UpdateStageInstance(stageInstanceUpdate discord.StageInstanceUpdate, opts ...rest.RequestOpt) (*StageInstance, rest.Error) {
	if !c.IsStageChannel() {
		unsupported(c)
	}
	stageInstance, err := c.Disgo().RestServices().StageInstanceService().UpdateStageInstance(c.ID(), stageInstanceUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return c.Disgo().EntityBuilder().CreateStageInstance(*stageInstance, CacheStrategyNoWs), nil
}

func (c *ChannelImpl) DeleteStageInstance(opts ...rest.RequestOpt) rest.Error {
	if !c.IsStageChannel() {
		unsupported(c)
	}
	return c.Disgo().RestServices().StageInstanceService().DeleteStageInstance(c.ID(), opts...)
}

func (c *ChannelImpl) IsModerator(member *Member) bool {
	if !c.IsStageChannel() {
		unsupported(c)
	}
	return member.Permissions().Has(discord.PermissionsStageModerator)
}

func unsupported(c *ChannelImpl) {
	panic(fmt.Sprintf("unsupported operation for '%d'", c.Type()))
}
