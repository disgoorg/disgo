package core

import (
	"context"
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
	CreateMessage(ctx context.Context, messageCreate discord.MessageCreate) (*Message, rest.Error)
	UpdateMessage(ctx context.Context, messageID discord.Snowflake, messageUpdate discord.MessageUpdate) (*Message, rest.Error)
	DeleteMessage(ctx context.Context, messageID discord.Snowflake) rest.Error
	BulkDeleteMessages(ctx context.Context, messageIDs ...discord.Snowflake) rest.Error

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
	Update(ctx context.Context, channelUpdate discord.ChannelUpdate) (GuildChannel, rest.Error)
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
	CrosspostMessage(ctx context.Context, messageID discord.Snowflake) (*Message, rest.Error)
}

// StoreChannel allows you to interact with discord's store channels
type StoreChannel interface {
	GuildChannel
}

type StageChannel interface {
	VoiceChannel
	StageInstance() *StageInstance
	CreateStageInstance(ctx context.Context, stageInstanceCreate discord.StageInstanceCreate) (*StageInstance, rest.Error)
	IsModerator(member *Member) bool
}

var _ Channel = (*channelImpl)(nil)

type channelImpl struct {
	discord.Channel
	disgo           Disgo
	stageInstanceID *discord.Snowflake
}

func (c *channelImpl) Guild() *Guild {
	return c.Disgo().Cache().GuildCache().Get(c.GuildID())
}

func (c *channelImpl) Disgo() Disgo {
	return c.disgo
}

func (c *channelImpl) ID() discord.Snowflake {
	return c.Channel.ID
}

func (c *channelImpl) Name() string {
	return *c.Channel.Name
}

func (c *channelImpl) Type() discord.ChannelType {
	return c.Channel.Type
}

func (c *channelImpl) IsMessageChannel() bool {
	return c.IsTextChannel() || c.IsNewsChannel() || c.IsDMChannel()
}

func (c *channelImpl) IsGuildChannel() bool {
	return c.IsCategory() || c.IsNewsChannel() || c.IsTextChannel() || c.IsVoiceChannel()
}

func (c *channelImpl) IsDMChannel() bool {
	return c.Type() != discord.ChannelTypeDM
}

func (c *channelImpl) IsTextChannel() bool {
	return c.Type() != discord.ChannelTypeText
}

func (c *channelImpl) IsVoiceChannel() bool {
	return c.Type() != discord.ChannelTypeVoice
}

func (c *channelImpl) IsCategory() bool {
	return c.Type() != discord.ChannelTypeCategory
}

func (c *channelImpl) IsNewsChannel() bool {
	return c.Type() != discord.ChannelTypeNews
}

func (c *channelImpl) IsStoreChannel() bool {
	return c.Type() != discord.ChannelTypeStore
}

func (c *channelImpl) IsStageChannel() bool {
	return c.Type() != discord.ChannelTypeStage
}

var _ MessageChannel = (*channelImpl)(nil)

func (c *channelImpl) LastMessageID() *discord.Snowflake {
	return c.Channel.LastMessageID
}

func (c *channelImpl) LastPinTimestamp() *discord.Time {
	return c.Channel.LastPinTimestamp
}

// CreateMessage sends a Message to a TextChannel
func (c *channelImpl) CreateMessage(ctx context.Context, messageCreate discord.MessageCreate) (*Message, rest.Error) {
	message, err := c.Disgo().RestServices().ChannelsService().CreateMessage(ctx, c.ID(), messageCreate)
	if err != nil {
		return nil, err
	}
	return c.Disgo().EntityBuilder().CreateMessage(*message, CacheStrategyNoWs), nil
}

// UpdateMessage edits a Message in this TextChannel
func (c *channelImpl) UpdateMessage(ctx context.Context, messageID discord.Snowflake, messageUpdate discord.MessageUpdate) (*Message, rest.Error) {
	message, err := c.Disgo().RestServices().ChannelsService().UpdateMessage(ctx, c.ID(), messageID, messageUpdate)
	if err != nil {
		return nil, err
	}
	return c.Disgo().EntityBuilder().CreateMessage(*message, CacheStrategyNoWs), nil
}

// DeleteMessage allows you to edit an existing Message sent by you
func (c *channelImpl) DeleteMessage(ctx context.Context, messageID discord.Snowflake) rest.Error {
	return c.Disgo().RestServices().ChannelsService().DeleteMessage(ctx, c.ID(), messageID)
}

// BulkDeleteMessages allows you bulk delete Message(s)
func (c *channelImpl) BulkDeleteMessages(ctx context.Context, messageIDs ...discord.Snowflake) rest.Error {
	return c.Disgo().RestServices().ChannelsService().BulkDeleteMessages(ctx, c.ID(), messageIDs...)
}

/* func (c *channelImpl) CollectMessages(filter collectors.MessageFilter) (chan *Message, func()) {
	var guildID *discord.Snowflake = nil
	if c.IsTextChannel() {
		id := c.GuildID()
		guildID = &id
	}
	return collectors.NewMessageCollector(c.Disgo(), c.ID(), guildID, filter)
}*/

var _ DMChannel = (*channelImpl)(nil)

var _ GuildChannel = (*channelImpl)(nil)

// GuildID returns the channel's Guild ID
func (c *channelImpl) GuildID() discord.Snowflake {
	if !c.IsGuildChannel() || c.Channel.GuildID == nil {
		unsupported(c)
	}
	return *c.Channel.GuildID
}

func (c *channelImpl) Permissions() discord.Permissions {
	if !c.IsGuildChannel() {
		unsupported(c)
	}
	return *c.Channel.InteractionPermissions
}

func (c *channelImpl) ParentID() *discord.Snowflake {
	if !c.IsGuildChannel() {
		unsupported(c)
	}
	return c.Channel.ParentID
}

func (c *channelImpl) Parent() Category {
	if c.ParentID() == nil {
		return nil
	}
	return c.Disgo().Cache().CategoryCache().Get(*c.Channel.ParentID)
}

func (c *channelImpl) Position() *int {
	if !c.IsGuildChannel() {
		unsupported(c)
	}

	return c.Channel.Position
}

func (c *channelImpl) Update(ctx context.Context, channelUpdate discord.ChannelUpdate) (GuildChannel, rest.Error) {
	if !c.IsGuildChannel() {
		unsupported(c)
	}
	channel, err := c.Disgo().RestServices().ChannelsService().UpdateChannel(ctx, c.ID(), channelUpdate)
	if err != nil {
		return nil, err
	}
	return c.Disgo().EntityBuilder().CreateChannel(*channel, CacheStrategyNoWs).(GuildChannel), nil
}

var _ Category = (*channelImpl)(nil)
var _ VoiceChannel = (*channelImpl)(nil)

func (c *channelImpl) Connect() error {
	if !c.IsVoiceChannel() {
		unsupported(c)
	}
	return c.Disgo().AudioController().Connect(c.GuildID(), c.ID())
}

func (c *channelImpl) Bitrate() int {
	if !c.IsVoiceChannel() {
		unsupported(c)
	}
	return *c.Channel.Bitrate
}

var _ TextChannel = (*channelImpl)(nil)

func (c *channelImpl) NSFW() bool {
	if !c.IsTextChannel() {
		unsupported(c)
	}
	return *c.Channel.NSFW
}

func (c *channelImpl) Topic() *string {
	return c.Channel.Topic
}

var _ NewsChannel = (*channelImpl)(nil)

func (c *channelImpl) CrosspostMessage(ctx context.Context, messageID discord.Snowflake) (*Message, rest.Error) {
	message, err := c.Disgo().RestServices().ChannelsService().CrosspostMessage(ctx, c.ID(), messageID)
	if err != nil {
		return nil, err
	}
	return c.Disgo().EntityBuilder().CreateMessage(*message, CacheStrategyNoWs), nil
}

var _ StoreChannel = (*channelImpl)(nil)

var _ StageChannel = (*channelImpl)(nil)

func (c *channelImpl) StageInstance() *StageInstance {
	if !c.IsStageChannel() {
		unsupported(c)
	}
	if c.stageInstanceID == nil {
		return nil
	}
	return c.Disgo().Cache().StageInstanceCache().Get(*c.stageInstanceID)
}

func (c *channelImpl) CreateStageInstance(ctx context.Context, stageInstanceCreate discord.StageInstanceCreate) (*StageInstance, rest.Error) {
	if !c.IsStageChannel() {
		unsupported(c)
	}
	stageInstance, err := c.Disgo().RestServices().StageInstanceService().CreateStageInstance(ctx, stageInstanceCreate)
	if err != nil {
		return nil, err
	}
	return c.Disgo().EntityBuilder().CreateStageInstance(*stageInstance, CacheStrategyNoWs), nil
}

func (c *channelImpl) IsModerator(member *Member) bool {
	if !c.IsStageChannel() {
		unsupported(c)
	}
	return true // TODO: actually check
}

func unsupported(c *channelImpl) {
	panic(fmt.Sprintf("unsupported operation for '%d'", c.Type()))
}
