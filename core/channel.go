package core

import (
	"fmt"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/rest/route"
)

var _ Mentionable = (*Channel)(nil)

type Channel struct {
	discord.Channel
	Bot                *Bot
	StageInstanceID    *discord.Snowflake
	ConnectedMemberIDs map[discord.Snowflake]struct{}
}

func (c *Channel) String() string {
	return fmt.Sprintf("<#%s>", c.ID)
}

func (c *Channel) Mention() string {
	return c.String()
}

func (c *Channel) Guild() *Guild {
	if !c.IsGuildChannel() {
		unsupportedChannelType(c)
	}
	return c.Bot.Caches.GuildCache().Get(c.GuildID)
}

func (c *Channel) Channels() []*Channel {
	if !c.IsCategory() {
		unsupportedChannelType(c)
	}
	return c.Bot.Caches.ChannelCache().FindAll(func(channel *Channel) bool {
		return channel.ParentID != nil && *channel.ParentID == c.ID
	})
}

func (c *Channel) Members() []*Member {
	if c.IsStoreChannel() {
		unsupportedChannelType(c)
	}
	var members []*Member
	if c.IsCategory() {
		memberIds := make(map[discord.Snowflake]struct{})
		for _, channel := range c.Channels() {
			if channel.IsStoreChannel() {
				continue
			}
			for _, member := range channel.Members() {
				if _, ok := memberIds[member.ID]; ok {
					continue
				}
				members = append(members, member)
				memberIds[member.ID] = struct{}{}
			}
		}
		return members
	} else if c.IsTextChannel() || c.IsNewsChannel() {
		members = c.Bot.Caches.MemberCache().FindAll(func(member *Member) bool {
			return member.ChannelPermissions(c).Has(discord.PermissionViewChannel)
		})
	} else if c.IsVoiceChannel() || c.IsStageChannel() {
		members = c.Bot.Caches.MemberCache().FindAll(func(member *Member) bool {
			_, ok := c.ConnectedMemberIDs[member.ID]
			return ok
		})
	}
	return members
}

func (c *Channel) PermissionOverwrite(overwriteType discord.PermissionOverwriteType, id discord.Snowflake) *discord.PermissionOverwrite {
	for _, overwrite := range c.PermissionOverwrites {
		if overwrite.Type == overwriteType && overwrite.ID == id {
			return &overwrite
		}
	}
	return nil
}

func (c *Channel) IsMessageChannel() bool {
	return c.IsTextChannel() || c.IsNewsChannel() || c.IsDMChannel()
}

func (c *Channel) IsGuildChannel() bool {
	return c.IsCategory() || c.IsNewsChannel() || c.IsTextChannel() || c.IsVoiceChannel()
}

func (c *Channel) IsAudioChannel() bool {
	return c.IsVoiceChannel() || c.IsStageChannel()
}

func (c *Channel) IsDMChannel() bool {
	return c.Type == discord.ChannelTypeDM
}

func (c *Channel) IsTextChannel() bool {
	return c.Type == discord.ChannelTypeGuildText
}

func (c *Channel) IsVoiceChannel() bool {
	return c.Type == discord.ChannelTypeGuildVoice
}

func (c *Channel) IsCategory() bool {
	return c.Type == discord.ChannelTypeGuildCategory
}

func (c *Channel) IsNewsChannel() bool {
	return c.Type == discord.ChannelTypeGuildNews
}

func (c *Channel) IsStoreChannel() bool {
	return c.Type == discord.ChannelTypeGuildStore
}

func (c *Channel) IsStageChannel() bool {
	return c.Type == discord.ChannelTypeGuildStageVoice
}

// MessageFilter used to filter Message(s) in a MessageCollector
type MessageFilter func(message *Message) bool

func (c *Channel) CollectMessages(filter MessageFilter) (<-chan *Message, func()) {
	if !c.IsMessageChannel() {
		unsupportedChannelType(c)
	}
	return c.Bot.EventManager.Config().NewMessageCollector(c, filter)
}

// CreateMessage sends a Message to a Channel
func (c *Channel) CreateMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*Message, rest.Error) {
	if !c.IsMessageChannel() {
		unsupportedChannelType(c)
	}
	message, err := c.Bot.RestServices.ChannelService().CreateMessage(c.ID, messageCreate, opts...)
	if err != nil {
		return nil, err
	}
	return c.Bot.EntityBuilder.CreateMessage(*message, CacheStrategyNoWs), nil
}

// UpdateMessage edits a Message in this Channel
func (c *Channel) UpdateMessage(messageID discord.Snowflake, messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*Message, rest.Error) {
	if !c.IsMessageChannel() {
		unsupportedChannelType(c)
	}
	message, err := c.Bot.RestServices.ChannelService().UpdateMessage(c.ID, messageID, messageUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return c.Bot.EntityBuilder.CreateMessage(*message, CacheStrategyNoWs), nil
}

// DeleteMessage allows you to edit an existing Message sent by you
func (c *Channel) DeleteMessage(messageID discord.Snowflake, opts ...rest.RequestOpt) rest.Error {
	if !c.IsMessageChannel() {
		unsupportedChannelType(c)
	}
	return c.Bot.RestServices.ChannelService().DeleteMessage(c.ID, messageID, opts...)
}

// BulkDeleteMessages allows you bulk delete Message(s)
func (c *Channel) BulkDeleteMessages(messageIDs []discord.Snowflake, opts ...rest.RequestOpt) rest.Error {
	if !c.IsMessageChannel() {
		unsupportedChannelType(c)
	}
	return c.Bot.RestServices.ChannelService().BulkDeleteMessages(c.ID, messageIDs, opts...)
}

// GetMessage gets a Message from the Channel
func (c *Channel) GetMessage(messageID discord.Snowflake, opts ...rest.RequestOpt) (*Message, rest.Error) {
	if !c.IsMessageChannel() {
		unsupportedChannelType(c)
	}
	message, err := c.Bot.RestServices.ChannelService().GetMessage(c.ID, messageID, opts...)
	if err != nil {
		return nil, err
	}
	return c.Bot.EntityBuilder.CreateMessage(*message, CacheStrategyNoWs), nil
}

// GetMessages gets multiple Message(s) from the Channel
func (c *Channel) GetMessages(around discord.Snowflake, before discord.Snowflake, after discord.Snowflake, limit int, opts ...rest.RequestOpt) ([]*Message, rest.Error) {
	if !c.IsMessageChannel() {
		unsupportedChannelType(c)
	}
	messages, err := c.Bot.RestServices.ChannelService().GetMessages(c.ID, around, before, after, limit, opts...)
	if err != nil {
		return nil, err
	}
	coreMessages := make([]*Message, len(messages))
	for i, message := range messages {
		coreMessages[i] = c.Bot.EntityBuilder.CreateMessage(message, CacheStrategyNoWs)
	}
	return coreMessages, nil
}

func (c *Channel) Parent() *Channel {
	if c.ParentID == nil {
		return nil
	}
	return c.Bot.Caches.ChannelCache().Get(*c.Channel.ParentID)
}

func (c *Channel) Update(channelUpdate discord.ChannelUpdate, opts ...rest.RequestOpt) (*Channel, rest.Error) {
	if !c.IsGuildChannel() {
		unsupportedChannelType(c)
	}
	channel, err := c.Bot.RestServices.ChannelService().UpdateChannel(c.ID, channelUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return c.Bot.EntityBuilder.CreateChannel(*channel, CacheStrategyNoWs), nil
}

func (c *Channel) Connect() error {
	if !c.IsAudioChannel() {
		unsupportedChannelType(c)
	}
	return c.Bot.AudioController.Connect(c.GuildID, c.ID)
}

func (c *Channel) CrosspostMessage(messageID discord.Snowflake, opts ...rest.RequestOpt) (*Message, rest.Error) {
	message, err := c.Bot.RestServices.ChannelService().CrosspostMessage(c.ID, messageID, opts...)
	if err != nil {
		return nil, err
	}
	return c.Bot.EntityBuilder.CreateMessage(*message, CacheStrategyNoWs), nil
}

func (c *Channel) StageInstance() *StageInstance {
	if !c.IsStageChannel() {
		unsupportedChannelType(c)
	}
	if c.StageInstanceID == nil {
		return nil
	}
	return c.Bot.Caches.StageInstanceCache().Get(*c.StageInstanceID)
}

func (c *Channel) CreateStageInstance(stageInstanceCreate discord.StageInstanceCreate, opts ...rest.RequestOpt) (*StageInstance, rest.Error) {
	if !c.IsStageChannel() {
		unsupportedChannelType(c)
	}
	stageInstance, err := c.Bot.RestServices.StageInstanceService().CreateStageInstance(stageInstanceCreate, opts...)
	if err != nil {
		return nil, err
	}
	return c.Bot.EntityBuilder.CreateStageInstance(*stageInstance, CacheStrategyNoWs), nil
}

func (c *Channel) UpdateStageInstance(stageInstanceUpdate discord.StageInstanceUpdate, opts ...rest.RequestOpt) (*StageInstance, rest.Error) {
	if !c.IsStageChannel() {
		unsupportedChannelType(c)
	}
	stageInstance, err := c.Bot.RestServices.StageInstanceService().UpdateStageInstance(c.ID, stageInstanceUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return c.Bot.EntityBuilder.CreateStageInstance(*stageInstance, CacheStrategyNoWs), nil
}

func (c *Channel) DeleteStageInstance(opts ...rest.RequestOpt) rest.Error {
	if !c.IsStageChannel() {
		unsupportedChannelType(c)
	}
	return c.Bot.RestServices.StageInstanceService().DeleteStageInstance(c.ID, opts...)
}

// GetIconURL returns the Icon URL of this channel.
// This will be nil for every discord.ChannelType except discord.ChannelTypeGroupDM
func (c *Channel) GetIconURL(size int) *string {
	return discord.FormatAssetURL(route.ChannelIcon, c.ID, c.Icon, size)
}

func (c *Channel) IsModerator(member *Member) bool {
	if !c.IsStageChannel() {
		unsupportedChannelType(c)
	}
	return member.Permissions().Has(discord.PermissionsStageModerator)
}

func unsupportedChannelType(c *Channel) {
	panic(fmt.Sprintf("unsupported ChannelType operation for '%d'", c.Type))
}
