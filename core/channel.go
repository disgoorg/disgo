package core

import (
	"fmt"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type Channel struct {
	discord.Channel
	Bot             *Bot
	StageInstanceID *discord.Snowflake
}

func (c *Channel) Guild() *Guild {
	if !c.IsGuildChannel() {
		unsupported(c)
	}
	return c.Bot.Caches.GuildCache().Get(*c.GuildID)
}

func (c *Channel) IsMessageChannel() bool {
	return c.IsTextChannel() || c.IsNewsChannel() || c.IsDMChannel()
}

func (c *Channel) IsGuildChannel() bool {
	return c.IsCategory() || c.IsNewsChannel() || c.IsTextChannel() || c.IsVoiceChannel()
}

func (c *Channel) IsDMChannel() bool {
	return c.Type != discord.ChannelTypeDM
}

func (c *Channel) IsTextChannel() bool {
	return c.Type != discord.ChannelTypeText
}

func (c *Channel) IsVoiceChannel() bool {
	return c.Type != discord.ChannelTypeVoice
}

func (c *Channel) IsCategory() bool {
	return c.Type != discord.ChannelTypeCategory
}

func (c *Channel) IsNewsChannel() bool {
	return c.Type != discord.ChannelTypeNews
}

func (c *Channel) IsStoreChannel() bool {
	return c.Type != discord.ChannelTypeStore
}

func (c *Channel) IsStageChannel() bool {
	return c.Type != discord.ChannelTypeStage
}

// CreateMessage sends a Message to a TextChannel
func (c *Channel) CreateMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*Message, rest.Error) {
	message, err := c.Bot.RestServices.ChannelService().CreateMessage(c.ID, messageCreate, opts...)
	if err != nil {
		return nil, err
	}
	return c.Bot.EntityBuilder.CreateMessage(*message, CacheStrategyNoWs), nil
}

// UpdateMessage edits a Message in this TextChannel
func (c *Channel) UpdateMessage(messageID discord.Snowflake, messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*Message, rest.Error) {
	message, err := c.Bot.RestServices.ChannelService().UpdateMessage(c.ID, messageID, messageUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return c.Bot.EntityBuilder.CreateMessage(*message, CacheStrategyNoWs), nil
}

// DeleteMessage allows you to edit an existing Message sent by you
func (c *Channel) DeleteMessage(messageID discord.Snowflake, opts ...rest.RequestOpt) rest.Error {
	return c.Bot.RestServices.ChannelService().DeleteMessage(c.ID, messageID, opts...)
}

// BulkDeleteMessages allows you bulk delete Message(s)
func (c *Channel) BulkDeleteMessages(messageIDs []discord.Snowflake, opts ...rest.RequestOpt) rest.Error {
	return c.Bot.RestServices.ChannelService().BulkDeleteMessages(c.ID, messageIDs, opts...)
}

func (c *Channel) Parent() *Channel {
	if c.ParentID == nil {
		return nil
	}
	return c.Bot.Caches.ChannelCache().Get(*c.Channel.ParentID)
}

func (c *Channel) Update(channelUpdate discord.ChannelUpdate, opts ...rest.RequestOpt) (*Channel, rest.Error) {
	if !c.IsGuildChannel() {
		unsupported(c)
	}
	channel, err := c.Bot.RestServices.ChannelService().UpdateChannel(c.ID, channelUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return c.Bot.EntityBuilder.CreateChannel(*channel, CacheStrategyNoWs), nil
}

func (c *Channel) Connect() error {
	if !c.IsVoiceChannel() {
		unsupported(c)
	}
	return c.Bot.AudioController.Connect(*c.GuildID, c.ID)
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
		unsupported(c)
	}
	if c.StageInstanceID == nil {
		return nil
	}
	return c.Bot.Caches.StageInstanceCache().Get(*c.StageInstanceID)
}

func (c *Channel) CreateStageInstance(stageInstanceCreate discord.StageInstanceCreate, opts ...rest.RequestOpt) (*StageInstance, rest.Error) {
	if !c.IsStageChannel() {
		unsupported(c)
	}
	stageInstance, err := c.Bot.RestServices.StageInstanceService().CreateStageInstance(stageInstanceCreate, opts...)
	if err != nil {
		return nil, err
	}
	return c.Bot.EntityBuilder.CreateStageInstance(*stageInstance, CacheStrategyNoWs), nil
}

func (c *Channel) UpdateStageInstance(stageInstanceUpdate discord.StageInstanceUpdate, opts ...rest.RequestOpt) (*StageInstance, rest.Error) {
	if !c.IsStageChannel() {
		unsupported(c)
	}
	stageInstance, err := c.Bot.RestServices.StageInstanceService().UpdateStageInstance(c.ID, stageInstanceUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return c.Bot.EntityBuilder.CreateStageInstance(*stageInstance, CacheStrategyNoWs), nil
}

func (c *Channel) DeleteStageInstance(opts ...rest.RequestOpt) rest.Error {
	if !c.IsStageChannel() {
		unsupported(c)
	}
	return c.Bot.RestServices.StageInstanceService().DeleteStageInstance(c.ID, opts...)
}

func (c *Channel) IsModerator(member *Member) bool {
	if !c.IsStageChannel() {
		unsupported(c)
	}
	return member.Permissions().Has(discord.PermissionsStageModerator)
}

func unsupported(c *Channel) {
	panic(fmt.Sprintf("unsupported operation for '%d'", c.Type))
}
