package core

import (
	"fmt"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/rest/route"
)

type Channel interface {
	discord.Channel
}

type GuildChannel interface {
	discord.GuildChannel
}

type MessageChannel interface {
	discord.MessageChannel
}

type GuildMessageChannel interface {
	discord.GuildMessageChannel
}

type GuildThread interface {
	discord.GuildThread
}

type AudioChannel interface {
	discord.AudioChannel
}

type GuildTextChannel struct {
	discord.GuildTextChannel
	Bot *Bot
}

func (c *GuildTextChannel) String() string {
	return channelMention(c.ID)
}

func (c *GuildTextChannel) Guild() *Guild {
	return channelGuild(c.Bot, c.GuildID)
}

func (c *GuildTextChannel) Parent() *GuildCategoryChannel {
	if c.ParentID == nil {
		return nil
	}
	return c.Bot.Caches.ChannelCache().Get(*c.ParentID).(*GuildCategoryChannel)
}

func (c *GuildTextChannel) Members() []*Member {
	return c.Bot.Caches.MemberCache().FindAll(func(member *Member) bool {
		return member.ChannelPermissions(c).Has(discord.PermissionViewChannel)
	})
}

type DMChannel struct {
	discord.DMChannel
	Bot          *Bot
	RecipientIDs []discord.Snowflake
}

func (c *DMChannel) String() string {
	return channelMention(c.ID)
}

type GuildVoiceChannel struct {
	discord.GuildVoiceChannel
	Bot                *Bot
	ConnectedMemberIDs map[discord.Snowflake]struct{}
}

func (c *GuildVoiceChannel) String() string {
	return channelMention(c.ID)
}

func (c *GuildVoiceChannel) Guild() *Guild {
	return channelGuild(c.Bot, c.GuildID)
}

func (c *GuildVoiceChannel) Parent() *GuildCategoryChannel {
	if c.ParentID == nil {
		return nil
	}
	return c.Bot.Caches.ChannelCache().Get(*c.ParentID).(*GuildCategoryChannel)
}

func (c *GuildVoiceChannel) Connect() error {
	return c.Bot.AudioController.Connect(c.GuildID, c.ID)
}

func (c *GuildVoiceChannel) Members() []*Member {
	return c.Bot.Caches.MemberCache().FindAll(func(member *Member) bool {
		_, ok := c.ConnectedMemberIDs[member.ID]
		return ok
	})
}

type GroupDMChannel struct {
	discord.GroupDMChannel
	Bot          *Bot
	RecipientIDs []discord.Snowflake
}

func (c *GroupDMChannel) String() string {
	return channelMention(c.ID)
}

// GetIconURL returns the Icon URL of this channel.
func (c *GroupDMChannel) GetIconURL(size int) *string {
	return discord.FormatAssetURL(route.ChannelIcon, c.ID, c.Icon, size)
}

type GuildCategoryChannel struct {
	discord.GuildCategoryChannel
	Bot *Bot
}

func (c *GuildCategoryChannel) String() string {
	return channelMention(c.ID)
}

func (c *GuildCategoryChannel) Guild() *Guild {
	return channelGuild(c.Bot, c.GuildID)
}

func (c *GuildCategoryChannel) Channels() []GuildChannel {
	channels := c.Bot.Caches.ChannelCache().FindAll(func(channel Channel) bool {
		switch ch := channel.(type) {
		case *GuildTextChannel:
			return ch.ParentID != nil && *ch.ParentID == c.ID

		default:
			return false
		}
	})
	guildChannels := make([]GuildChannel, len(channels))
	for i := range channels {
		guildChannels[i] = channels[i].(GuildChannel)
	}
	return guildChannels
}

func (c *GuildCategoryChannel) Members() []*Member {
	var members []*Member
	memberIds := make(map[discord.Snowflake]struct{})
	for _, channel := range c.Channels() {
		var chMembers []*Member
		switch ch := channel.(type) {
		case *GuildTextChannel:
			chMembers = ch.Members()

		default:
			continue
		}
		for _, member := range chMembers {
			if _, ok := memberIds[member.ID]; ok {
				continue
			}
			members = append(members, member)
			memberIds[member.ID] = struct{}{}
		}
	}
	return members
}

type GuildNewsChannel struct {
	discord.GuildNewsChannel
	Bot *Bot
}

func (c *GuildNewsChannel) String() string {
	return channelMention(c.ID)
}

func (c *GuildNewsChannel) Guild() *Guild {
	return channelGuild(c.Bot, c.GuildID)
}

func (c *GuildNewsChannel) Parent() *GuildCategoryChannel {
	if c.ParentID == nil {
		return nil
	}
	return c.Bot.Caches.ChannelCache().Get(*c.ParentID).(*GuildCategoryChannel)
}

func (c *GuildNewsChannel) CrosspostMessage(messageID discord.Snowflake, opts ...rest.RequestOpt) (*Message, error) {
	message, err := c.Bot.RestServices.ChannelService().CrosspostMessage(c.ID, messageID, opts...)
	if err != nil {
		return nil, err
	}
	return c.Bot.EntityBuilder.CreateMessage(*message, CacheStrategyNoWs), nil
}

type GuildStoreChannel struct {
	discord.GuildStoreChannel
	Bot *Bot
}

func (c *GuildStoreChannel) String() string {
	return channelMention(c.ID)
}

func (c *GuildStoreChannel) Guild() *Guild {
	return channelGuild(c.Bot, c.GuildID)
}

func (c *GuildStoreChannel) Parent() *GuildCategoryChannel {
	if c.ParentID == nil {
		return nil
	}
	return c.Bot.Caches.ChannelCache().Get(*c.ParentID).(*GuildCategoryChannel)
}

type GuildNewsThread struct {
	discord.GuildNewsThread
	Bot *Bot
}

func (c *GuildNewsThread) String() string {
	return channelMention(c.ID)
}

func (c *GuildNewsThread) Guild() *Guild {
	return channelGuild(c.Bot, c.GuildID)
}

func (c *GuildNewsThread) Parent() *GuildNewsChannel {
	return c.Bot.Caches.ChannelCache().Get(c.ParentID).(*GuildNewsChannel)
}

func (c *GuildNewsThread) Members() []*Member {
	return c.Bot.Caches.MemberCache().FindAll(func(member *Member) bool {
		return member.ChannelPermissions(c).Has(discord.PermissionViewChannel)
	})
}

func (c *GuildNewsThread) ThreadMembers() []*ThreadMember {
	return c.Bot.Caches.ThreadMemberCache().ThreadAll(c.ID)
}

func (c *GuildNewsThread) ThreadMembersCache() map[discord.Snowflake]*ThreadMember {
	return c.Bot.Caches.ThreadMemberCache().ThreadCache(c.ID)
}

type GuildPublicThread struct {
	discord.GuildNewsThread
	Bot *Bot
}

func (c *GuildPublicThread) String() string {
	return channelMention(c.ID)
}

func (c *GuildPublicThread) Guild() *Guild {
	return channelGuild(c.Bot, c.GuildID)
}

func (c *GuildPublicThread) Parent() *GuildTextChannel {
	return c.Bot.Caches.ChannelCache().Get(c.ParentID).(*GuildTextChannel)
}

func (c *GuildPublicThread) Members() []*Member {
	return c.Bot.Caches.MemberCache().FindAll(func(member *Member) bool {
		return member.ChannelPermissions(c).Has(discord.PermissionViewChannel)
	})
}

func (c *GuildPublicThread) ThreadMembers() []*ThreadMember {
	return c.Bot.Caches.ThreadMemberCache().ThreadAll(c.ID)
}

func (c *GuildPublicThread) ThreadMembersCache() map[discord.Snowflake]*ThreadMember {
	return c.Bot.Caches.ThreadMemberCache().ThreadCache(c.ID)
}

type GuildPrivateThread struct {
	discord.GuildNewsThread
	Bot *Bot
}

func (c *GuildPrivateThread) String() string {
	return channelMention(c.ID)
}

func (c *GuildPrivateThread) Guild() *Guild {
	return channelGuild(c.Bot, c.GuildID)
}

func (c *GuildPrivateThread) Parent() *GuildTextChannel {
	return c.Bot.Caches.ChannelCache().Get(c.ParentID).(*GuildTextChannel)
}

func (c *GuildPrivateThread) Members() []*Member {
	return c.Bot.Caches.MemberCache().FindAll(func(member *Member) bool {
		return member.ChannelPermissions(c).Has(discord.PermissionViewChannel)
	})
}

func (c *GuildPrivateThread) ThreadMembers() []*ThreadMember {
	return c.Bot.Caches.ThreadMemberCache().ThreadAll(c.ID)
}

func (c *GuildPrivateThread) ThreadMembersCache() map[discord.Snowflake]*ThreadMember {
	return c.Bot.Caches.ThreadMemberCache().ThreadCache(c.ID)
}

type GuildStageVoiceChannel struct {
	discord.GuildStageVoiceChannel
	Bot                *Bot
	StageInstanceID    *discord.Snowflake
	ConnectedMemberIDs map[discord.Snowflake]struct{}
}

func (c *GuildStageVoiceChannel) String() string {
	return channelMention(c.ID)
}

func (c *GuildStageVoiceChannel) Update(channelUpdate discord.GuildStageVoiceChannelUpdate, opts ...rest.RequestOpt) (*GuildStageVoiceChannel, error) {
	channel, err := updateChannel(c.Bot, c.ID, channelUpdate, opts...)
}

func (c *GuildStageVoiceChannel) Guild() *Guild {
	return channelGuild(c.Bot, c.GuildID)
}

func (c *GuildStageVoiceChannel) Parent() *GuildCategoryChannel {
	if c.ParentID == nil {
		return nil
	}
	return c.Bot.Caches.ChannelCache().Get(*c.ParentID).(*GuildCategoryChannel)
}

func (c *GuildStageVoiceChannel) Connect() error {
	return c.Bot.AudioController.Connect(c.GuildID, c.ID)
}

func (c *GuildStageVoiceChannel) Members() []*Member {
	return c.Bot.Caches.MemberCache().FindAll(func(member *Member) bool {
		_, ok := c.ConnectedMemberIDs[member.ID]
		return ok
	})
}

func (c *GuildStageVoiceChannel) IsModerator(member *Member) bool {
	return member.Permissions().Has(discord.PermissionsStageModerator)
}

func (c *GuildStageVoiceChannel) StageInstance() *StageInstance {
	if c.StageInstanceID == nil {
		return nil
	}
	return c.Bot.Caches.StageInstanceCache().Get(*c.StageInstanceID)
}

func (c *GuildStageVoiceChannel) CreateStageInstance(stageInstanceCreate discord.StageInstanceCreate, opts ...rest.RequestOpt) (*StageInstance, error) {
	stageInstance, err := c.Bot.RestServices.StageInstanceService().CreateStageInstance(stageInstanceCreate, opts...)
	if err != nil {
		return nil, err
	}
	return c.Bot.EntityBuilder.CreateStageInstance(*stageInstance, CacheStrategyNoWs), nil
}

func (c *GuildStageVoiceChannel) UpdateStageInstance(stageInstanceUpdate discord.StageInstanceUpdate, opts ...rest.RequestOpt) (*StageInstance, error) {
	stageInstance, err := c.Bot.RestServices.StageInstanceService().UpdateStageInstance(c.ID, stageInstanceUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return c.Bot.EntityBuilder.CreateStageInstance(*stageInstance, CacheStrategyNoWs), nil
}

func (c *GuildStageVoiceChannel) DeleteStageInstance(opts ...rest.RequestOpt) error {
	return c.Bot.RestServices.StageInstanceService().DeleteStageInstance(c.ID, opts...)
}

//--------------------------------------------

func (c *Channel) PermissionOverwrite(overwriteType discord.PermissionOverwriteType, id discord.Snowflake) *discord.PermissionOverwrite {
	for _, overwrite := range c.PermissionOverwrites {
		if overwrite.Type == overwriteType && overwrite.ID == id {
			return &overwrite
		}
	}
	return nil
}



func channelMention(id discord.Snowflake) string {
	return fmt.Sprintf("<#%s>", id)
}

func channelGuild(bot *Bot, guildID discord.Snowflake) *Guild {
	return bot.Caches.GuildCache().Get(guildID)
}

func createThread(bot *Bot, channelID discord.Snowflake, threadCreate discord.ThreadCreate, opts ...rest.RequestOpt) (discord.GuildThread, error) {
	channel, err := bot.RestServices.ThreadService().CreateThread(channelID, threadCreate, opts...)
	if err != nil {
		return nil, err
	}
	return bot.EntityBuilder.CreateChannel(channel, CacheStrategyNo).(discord.GuildThread), nil
}

func createThreadWithMessage(bot *Bot, channelID discord.Snowflake, messageID discord.Snowflake, threadCreateWithMessage discord.ThreadCreateWithMessage, opts ...rest.RequestOpt) (discord.GuildThread, error) {
	channel, err := bot.RestServices.ThreadService().CreateThreadWithMessage(channelID, messageID, threadCreateWithMessage, opts...)
	if err != nil {
		return nil, err
	}
	return bot.EntityBuilder.CreateChannel(channel, CacheStrategyNo).(discord.GuildThread), nil
}

func createMessage(bot *Bot, channelID discord.Snowflake, messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*Message, error) {
	message, err := bot.RestServices.ChannelService().CreateMessage(channelID, messageCreate, opts...)
	if err != nil {
		return nil, err
	}
	return bot.EntityBuilder.CreateMessage(*message, CacheStrategyNoWs), nil
}

func updateMessage(bot *Bot, channelID discord.Snowflake, messageID discord.Snowflake, messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*Message, error) {
	message, err := bot.RestServices.ChannelService().UpdateMessage(channelID, messageID, messageUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return bot.EntityBuilder.CreateMessage(*message, CacheStrategyNoWs), nil
}

func deleteMessage(bot *Bot, channelID discord.Snowflake, messageID discord.Snowflake, opts ...rest.RequestOpt) error {
	return bot.RestServices.ChannelService().DeleteMessage(channelID, messageID, opts...)
}

func bulkDeleteMessages(bot *Bot, channelID discord.Snowflake, messageIDs []discord.Snowflake, opts ...rest.RequestOpt) error {
	return bot.RestServices.ChannelService().BulkDeleteMessages(channelID, messageIDs, opts...)
}

func getMessage(bot *Bot, channelID discord.Snowflake, messageID discord.Snowflake, opts ...rest.RequestOpt) (*Message, error) {
	message, err := bot.RestServices.ChannelService().GetMessage(channelID, messageID, opts...)
	if err != nil {
		return nil, err
	}
	return bot.EntityBuilder.CreateMessage(*message, CacheStrategyNoWs), nil
}

func getMessages(bot *Bot, channelID discord.Snowflake, around discord.Snowflake, before discord.Snowflake, after discord.Snowflake, limit int, opts ...rest.RequestOpt) ([]*Message, error) {
	messages, err := bot.RestServices.ChannelService().GetMessages(channelID, around, before, after, limit, opts...)
	if err != nil {
		return nil, err
	}
	coreMessages := make([]*Message, len(messages))
	for i, message := range messages {
		coreMessages[i] = bot.EntityBuilder.CreateMessage(message, CacheStrategyNoWs)
	}
	return coreMessages, nil
}

func updateChannel(bot *Bot, channelID discord.Snowflake, channelUpdate discord.ChannelUpdate, opts ...rest.RequestOpt) (Channel, error) {
	channel, err := bot.RestServices.ChannelService().UpdateChannel(channelID, channelUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return bot.EntityBuilder.CreateChannel(channel, CacheStrategyNoWs), nil
}
