package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/rest/route"
	"github.com/pkg/errors"
)

type Channel interface {
	discord.Channel
	set(channel Channel) Channel
}

type GuildChannel interface {
	discord.GuildChannel
	Channel
	Guild() *Guild

	UpdateGuildChannel(guildChannelUpdate discord.GuildChannelUpdate, opts ...rest.RequestOpt) (GuildChannel, error)
	Delete(opts ...rest.RequestOpt) error

	PermissionOverwrites() []discord.PermissionOverwrite
	PermissionOverwrite(overwriteType discord.PermissionOverwriteType, id discord.Snowflake) discord.PermissionOverwrite
	RolePermissionOverwrite(id discord.Snowflake) *discord.RolePermissionOverwrite
	MemberPermissionOverwrite(id discord.Snowflake) *discord.MemberPermissionOverwrite
	SetPermissionOverwrite(overwriteType discord.PermissionOverwriteType, id discord.Snowflake, allow discord.Permissions, deny discord.Permissions, opts ...rest.RequestOpt) error
	UpdatePermissionOverwrite(overwriteType discord.PermissionOverwriteType, id discord.Snowflake, allow discord.Permissions, deny discord.Permissions, opts ...rest.RequestOpt) error
	DeletePermissionOverwrite(id discord.Snowflake, opts ...rest.RequestOpt) error

	Members() []*Member
}

type MessageChannel interface {
	discord.MessageChannel
	Channel

	SendTyping(opts ...rest.RequestOpt) error

	GetMessage(messageID discord.Snowflake, opts ...rest.RequestOpt) (*Message, error)
	GetMessages(around discord.Snowflake, before discord.Snowflake, after discord.Snowflake, limit int, opts ...rest.RequestOpt) ([]*Message, error)
	CreateMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*Message, error)
	UpdateMessage(messageID discord.Snowflake, messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*Message, error)
	DeleteMessage(messageID discord.Snowflake, opts ...rest.RequestOpt) error
	BulkDeleteMessages(messageIDs []discord.Snowflake, opts ...rest.RequestOpt) error

	AddReaction(messageID discord.Snowflake, emoji string, opts ...rest.RequestOpt) error
	RemoveOwnReaction(messageID discord.Snowflake, emoji string, opts ...rest.RequestOpt) error
	RemoveUserReaction(messageID discord.Snowflake, emoji string, userID discord.Snowflake, opts ...rest.RequestOpt) error
	RemoveAllReactions(messageID discord.Snowflake, opts ...rest.RequestOpt) error
	RemoveAllReactionsForEmoji(messageID discord.Snowflake, emoji string, opts ...rest.RequestOpt) error
}

type BaseGuildMessageChannel interface {
	discord.BaseGuildMessageChannel
	GuildChannel
	MessageChannel
}

type GuildMessageChannel interface {
	discord.GuildMessageChannel
	BaseGuildMessageChannel

	GetWebhooks(opts ...rest.RequestOpt) ([]Webhook, error)
	CreateWebhook(webhookCreate discord.WebhookCreate, opts ...rest.RequestOpt) (Webhook, error)
	DeleteWebhook(webhookID discord.Snowflake, opts ...rest.RequestOpt) error

	Threads() []GuildThread
	Thread(threadID discord.Snowflake) GuildThread

	CreateThread(theadCreate discord.ThreadCreate, opts ...rest.RequestOpt) (GuildThread, error)
	CreateThreadWithMessage(messageID discord.Snowflake, threadCreateWithMessage discord.ThreadCreateWithMessage, opts ...rest.RequestOpt) (GuildThread, error)

	GetPublicArchivedThreads(before discord.Time, limit int, opts ...rest.RequestOpt) ([]GuildThread, map[discord.Snowflake]*ThreadMember, bool, error)
}

type GuildThread interface {
	discord.GuildThread
	BaseGuildMessageChannel

	ParentMessageChannel() GuildMessageChannel
	SelfThreadMember() *ThreadMember
	ThreadMember(userID discord.Snowflake) *ThreadMember
	ThreadMembers() []*ThreadMember
	Join(opts ...rest.RequestOpt) error
	Leave(opts ...rest.RequestOpt) error
	AddThreadMember(userID discord.Snowflake, opts ...rest.RequestOpt) error
	RemoveThreadMember(userID discord.Snowflake, opts ...rest.RequestOpt) error
	GetThreadMember(userID discord.Snowflake, opts ...rest.RequestOpt) (*ThreadMember, error)
	GetThreadMembers(opts ...rest.RequestOpt) ([]*ThreadMember, error)
}

type GuildAudioChannel interface {
	discord.GuildAudioChannel
	GuildChannel

	Connect() error
	connectedMembers() map[discord.Snowflake]struct{}
}

var (
	_ Channel             = (*GuildTextChannel)(nil)
	_ GuildChannel        = (*GuildTextChannel)(nil)
	_ MessageChannel      = (*GuildTextChannel)(nil)
	_ GuildMessageChannel = (*GuildTextChannel)(nil)
)

type GuildTextChannel struct {
	discord.GuildTextChannel
	Bot *Bot
}

func (c *GuildTextChannel) set(channel Channel) Channel {
	switch ch := channel.(type) {
	case *GuildTextChannel:
		*c = *ch
		return c

	default:
		return c
	}
}

func (c *GuildTextChannel) UpdateGuildChannel(guildChannelUpdate discord.GuildChannelUpdate, opts ...rest.RequestOpt) (GuildChannel, error) {
	channel, err := updateChannel(c.Bot, c.ID(), guildChannelUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return channel.(GuildChannel), nil
}

// Update updates the GuildNewsChannel which can return either a GuildNewsChannel or a GuildTextChannel
func (c *GuildTextChannel) Update(channelUpdate discord.GuildTextChannelUpdate, opts ...rest.RequestOpt) (GuildMessageChannel, error) {
	channel, err := c.UpdateGuildChannel(channelUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return channel.(GuildMessageChannel), nil
}

func (c *GuildTextChannel) Delete(opts ...rest.RequestOpt) error {
	return deleteChannel(c.Bot, c.ID(), opts...)
}

func (c *GuildTextChannel) PermissionOverwrites() []discord.PermissionOverwrite {
	return c.ChannelPermissionOverwrites
}

func (c *GuildTextChannel) PermissionOverwrite(overwriteType discord.PermissionOverwriteType, id discord.Snowflake) discord.PermissionOverwrite {
	return getPermissionOverwrite(c, overwriteType, id)
}

func (c *GuildTextChannel) RolePermissionOverwrite(id discord.Snowflake) *discord.RolePermissionOverwrite {
	return getPermissionOverwrite(c, discord.PermissionOverwriteTypeRole, id).(*discord.RolePermissionOverwrite)
}

func (c *GuildTextChannel) MemberPermissionOverwrite(id discord.Snowflake) *discord.MemberPermissionOverwrite {
	return getPermissionOverwrite(c, discord.PermissionOverwriteTypeMember, id).(*discord.MemberPermissionOverwrite)
}

func (c *GuildTextChannel) SetPermissionOverwrite(overwriteType discord.PermissionOverwriteType, id discord.Snowflake, allow discord.Permissions, deny discord.Permissions, opts ...rest.RequestOpt) error {
	return setPermissionOverwrite(c.Bot, c.ID(), overwriteType, id, allow, deny, opts...)
}

func (c *GuildTextChannel) UpdatePermissionOverwrite(overwriteType discord.PermissionOverwriteType, id discord.Snowflake, allow discord.Permissions, deny discord.Permissions, opts ...rest.RequestOpt) error {
	return updatePermissionOverwrite(c.Bot, c, overwriteType, id, allow, deny, opts...)
}

func (c *GuildTextChannel) DeletePermissionOverwrite(id discord.Snowflake, opts ...rest.RequestOpt) error {
	return deletePermissionOverwrite(c.Bot, c.ID(), id, opts...)
}

func (c *GuildTextChannel) GetMessage(messageID discord.Snowflake, opts ...rest.RequestOpt) (*Message, error) {
	return getMessage(c.Bot, c.ID(), messageID, opts...)
}

func (c *GuildTextChannel) GetMessages(around discord.Snowflake, before discord.Snowflake, after discord.Snowflake, limit int, opts ...rest.RequestOpt) ([]*Message, error) {
	return getMessages(c.Bot, c.ID(), around, before, after, limit, opts...)
}

func (c *GuildTextChannel) CreateMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*Message, error) {
	return createMessage(c.Bot, c.ID(), messageCreate, opts...)
}

func (c *GuildTextChannel) UpdateMessage(messageID discord.Snowflake, messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*Message, error) {
	return updateMessage(c.Bot, c.ID(), messageID, messageUpdate, opts...)
}

func (c *GuildTextChannel) DeleteMessage(messageID discord.Snowflake, opts ...rest.RequestOpt) error {
	return deleteMessage(c.Bot, c.ID(), messageID, opts...)
}

func (c *GuildTextChannel) BulkDeleteMessages(messageIDs []discord.Snowflake, opts ...rest.RequestOpt) error {
	return bulkDeleteMessages(c.Bot, c.ID(), messageIDs, opts...)
}

func (c *GuildTextChannel) SendTyping(opts ...rest.RequestOpt) error {
	return sendTying(c.Bot, c.ID(), opts...)
}

func (c *GuildTextChannel) AddReaction(messageID discord.Snowflake, emoji string, opts ...rest.RequestOpt) error {
	return addReaction(c.Bot, c.ID(), messageID, emoji, opts...)
}

func (c *GuildTextChannel) RemoveOwnReaction(messageID discord.Snowflake, emoji string, opts ...rest.RequestOpt) error {
	return removeOwnReaction(c.Bot, c.ID(), messageID, emoji, opts...)
}

func (c *GuildTextChannel) RemoveUserReaction(messageID discord.Snowflake, emoji string, userID discord.Snowflake, opts ...rest.RequestOpt) error {
	return removeUserReaction(c.Bot, c.ID(), messageID, emoji, userID, opts...)
}

func (c *GuildTextChannel) RemoveAllReactions(messageID discord.Snowflake, opts ...rest.RequestOpt) error {
	return removeAllReactions(c.Bot, c.ID(), messageID, opts...)
}

func (c *GuildTextChannel) RemoveAllReactionsForEmoji(messageID discord.Snowflake, emoji string, opts ...rest.RequestOpt) error {
	return removeAllReactionsForEmoji(c.Bot, c.ID(), messageID, emoji, opts...)
}

func (c *GuildTextChannel) GetWebhooks(opts ...rest.RequestOpt) ([]Webhook, error) {
	return getWebhooks(c.Bot, c.ID(), opts...)
}

func (c *GuildTextChannel) CreateWebhook(webhookCreate discord.WebhookCreate, opts ...rest.RequestOpt) (Webhook, error) {
	return createWebhook(c.Bot, c.ID(), webhookCreate, opts...)
}

func (c *GuildTextChannel) DeleteWebhook(webhookID discord.Snowflake, opts ...rest.RequestOpt) error {
	return deleteWebhook(c.Bot, webhookID, opts...)
}

func (c *GuildTextChannel) Threads() []GuildThread {
	var threads []GuildThread
	c.Bot.Caches.Channels().ForAll(func(channel Channel) {
		if thread, ok := channel.(GuildThread); ok && thread.ParentID() == c.ID() {
			threads = append(threads, thread)
		}
	})
	return threads
}

func (c *GuildTextChannel) Thread(threadID discord.Snowflake) GuildThread {
	if thread, ok := c.Bot.Caches.Channels().Get(threadID).(GuildThread); ok {
		return thread
	}
	return nil
}

func (c *GuildTextChannel) PrivateThreads() []*GuildPrivateThread {
	var threads []*GuildPrivateThread
	c.Bot.Caches.Channels().ForAll(func(channel Channel) {
		if thread, ok := channel.(*GuildPrivateThread); ok && thread.ParentID() == c.ID() {
			threads = append(threads, thread)
		}
	})
	return threads
}

func (c *GuildTextChannel) PublicThreads() []*GuildPublicThread {
	var threads []*GuildPublicThread
	c.Bot.Caches.Channels().ForAll(func(channel Channel) {
		if thread, ok := channel.(*GuildPublicThread); ok && thread.ParentID() == c.ID() {
			threads = append(threads, thread)
		}
	})
	return threads
}

func (c *GuildTextChannel) CreateThread(theadCreate discord.ThreadCreate, opts ...rest.RequestOpt) (GuildThread, error) {
	return createThread(c.Bot, c.ID(), theadCreate, opts...)
}

func (c *GuildTextChannel) CreateThreadWithMessage(messageID discord.Snowflake, threadCreateWithMessage discord.ThreadCreateWithMessage, opts ...rest.RequestOpt) (GuildThread, error) {
	return createThreadWithMessage(c.Bot, c.ID(), messageID, threadCreateWithMessage, opts...)
}

func createThreadMembers(bot *Bot, members []discord.ThreadMember) map[discord.Snowflake]*ThreadMember {
	threadMembers := make(map[discord.Snowflake]*ThreadMember, len(members))
	for i := range members {
		threadMembers[members[i].ThreadID] = bot.EntityBuilder.CreateThreadMember(members[i], CacheStrategyNo)
	}
	return threadMembers
}

func (c *GuildTextChannel) GetPublicArchivedThreads(before discord.Time, limit int, opts ...rest.RequestOpt) ([]GuildThread, map[discord.Snowflake]*ThreadMember, bool, error) {
	return getPublicArchivedThreads(c.Bot, c.ID(), before, limit, opts...)
}

func createGuildPrivateThreads(bot *Bot, threads []discord.GuildThread) []*GuildPrivateThread {
	privateThreads := make([]*GuildPrivateThread, len(threads))
	for i := range threads {
		privateThreads[i] = bot.EntityBuilder.CreateChannel(threads[i], CacheStrategyNo).(*GuildPrivateThread)
	}
	return privateThreads
}

func (c *GuildTextChannel) GetPrivateArchivedThreads(before discord.Time, limit int, opts ...rest.RequestOpt) ([]*GuildPrivateThread, map[discord.Snowflake]*ThreadMember, bool, error) {
	getThreads, err := c.Bot.RestServices.ThreadService().GetPrivateArchivedThreads(c.ID(), before, limit, opts...)
	if err != nil {
		return nil, nil, false, err
	}

	return createGuildPrivateThreads(c.Bot, getThreads.Threads), createThreadMembers(c.Bot, getThreads.Members), getThreads.HasMore, nil
}

func (c *GuildTextChannel) GetJoinedPrivateAchievedThreads(before discord.Time, limit int, opts ...rest.RequestOpt) ([]*GuildPrivateThread, map[discord.Snowflake]*ThreadMember, bool, error) {
	getThreads, err := c.Bot.RestServices.ThreadService().GetJoinedPrivateAchievedThreads(c.ID(), before, limit, opts...)
	if err != nil {
		return nil, nil, false, err
	}

	return createGuildPrivateThreads(c.Bot, getThreads.Threads), createThreadMembers(c.Bot, getThreads.Members), getThreads.HasMore, nil
}

func (c *GuildTextChannel) Guild() *Guild {
	return channelGuild(c.Bot, c.GuildID())
}

func (c *GuildTextChannel) Parent() *GuildCategoryChannel {
	if c.ParentID == nil {
		return nil
	}
	return c.Bot.Caches.Channels().Get(*c.ParentID).(*GuildCategoryChannel)
}

func (c *GuildTextChannel) Members() []*Member {
	return viewMembers(c.Bot, c)
}

var (
	_ Channel        = (*DMChannel)(nil)
	_ MessageChannel = (*DMChannel)(nil)
)

type DMChannel struct {
	discord.DMChannel
	Bot          *Bot
	RecipientIDs []discord.Snowflake
}

func (c *DMChannel) set(channel Channel) Channel {
	switch ch := channel.(type) {
	case *DMChannel:
		*c = *ch
		return c

	default:
		return c
	}
}

func (c *DMChannel) Delete(opts ...rest.RequestOpt) error {
	return deleteChannel(c.Bot, c.ID(), opts...)
}

func (c *DMChannel) GetMessage(messageID discord.Snowflake, opts ...rest.RequestOpt) (*Message, error) {
	return getMessage(c.Bot, c.ID(), messageID, opts...)
}

func (c *DMChannel) GetMessages(around discord.Snowflake, before discord.Snowflake, after discord.Snowflake, limit int, opts ...rest.RequestOpt) ([]*Message, error) {
	return getMessages(c.Bot, c.ID(), around, before, after, limit, opts...)
}

func (c *DMChannel) CreateMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*Message, error) {
	return createMessage(c.Bot, c.ID(), messageCreate, opts...)
}

func (c *DMChannel) UpdateMessage(messageID discord.Snowflake, messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*Message, error) {
	return updateMessage(c.Bot, c.ID(), messageID, messageUpdate, opts...)
}

func (c *DMChannel) DeleteMessage(messageID discord.Snowflake, opts ...rest.RequestOpt) error {
	return deleteMessage(c.Bot, c.ID(), messageID, opts...)
}

func (c *DMChannel) BulkDeleteMessages(messageIDs []discord.Snowflake, opts ...rest.RequestOpt) error {
	return bulkDeleteMessages(c.Bot, c.ID(), messageIDs, opts...)
}

func (c *DMChannel) SendTyping(opts ...rest.RequestOpt) error {
	return sendTying(c.Bot, c.ID(), opts...)
}

func (c *DMChannel) AddReaction(messageID discord.Snowflake, emoji string, opts ...rest.RequestOpt) error {
	return addReaction(c.Bot, c.ID(), messageID, emoji, opts...)
}

func (c *DMChannel) RemoveOwnReaction(messageID discord.Snowflake, emoji string, opts ...rest.RequestOpt) error {
	return removeOwnReaction(c.Bot, c.ID(), messageID, emoji, opts...)
}

func (c *DMChannel) RemoveUserReaction(messageID discord.Snowflake, emoji string, userID discord.Snowflake, opts ...rest.RequestOpt) error {
	return removeUserReaction(c.Bot, c.ID(), messageID, emoji, userID, opts...)
}

func (c *DMChannel) RemoveAllReactions(messageID discord.Snowflake, opts ...rest.RequestOpt) error {
	return removeAllReactions(c.Bot, c.ID(), messageID, opts...)
}

func (c *DMChannel) RemoveAllReactionsForEmoji(messageID discord.Snowflake, emoji string, opts ...rest.RequestOpt) error {
	return removeAllReactionsForEmoji(c.Bot, c.ID(), messageID, emoji, opts...)
}

var (
	_ Channel           = (*GuildVoiceChannel)(nil)
	_ GuildChannel      = (*GuildVoiceChannel)(nil)
	_ GuildAudioChannel = (*GuildVoiceChannel)(nil)
)

type GuildVoiceChannel struct {
	discord.GuildVoiceChannel
	Bot                *Bot
	ConnectedMemberIDs map[discord.Snowflake]struct{}
}

func (c *GuildVoiceChannel) set(channel Channel) Channel {
	switch ch := channel.(type) {
	case *GuildVoiceChannel:
		*c = *ch
		return c

	default:
		return c
	}
}

func (c *GuildVoiceChannel) UpdateGuildChannel(guildChannelUpdate discord.GuildChannelUpdate, opts ...rest.RequestOpt) (GuildChannel, error) {
	channel, err := updateChannel(c.Bot, c.ID(), guildChannelUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return channel.(GuildChannel), nil
}

// Update updates the GuildNewsChannel which can return either a GuildNewsChannel or a GuildTextChannel
func (c *GuildVoiceChannel) Update(channelUpdate discord.GuildVoiceChannelUpdate, opts ...rest.RequestOpt) (*GuildVoiceChannel, error) {
	channel, err := c.UpdateGuildChannel(channelUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return channel.(*GuildVoiceChannel), nil
}

func (c *GuildVoiceChannel) Delete(opts ...rest.RequestOpt) error {
	return deleteChannel(c.Bot, c.ID(), opts...)
}

func (c *GuildVoiceChannel) PermissionOverwrites() []discord.PermissionOverwrite {
	return c.ChannelPermissionOverwrites
}

func (c *GuildVoiceChannel) PermissionOverwrite(overwriteType discord.PermissionOverwriteType, id discord.Snowflake) discord.PermissionOverwrite {
	return getPermissionOverwrite(c, overwriteType, id)
}

func (c *GuildVoiceChannel) RolePermissionOverwrite(id discord.Snowflake) *discord.RolePermissionOverwrite {
	return getPermissionOverwrite(c, discord.PermissionOverwriteTypeRole, id).(*discord.RolePermissionOverwrite)
}

func (c *GuildVoiceChannel) MemberPermissionOverwrite(id discord.Snowflake) *discord.MemberPermissionOverwrite {
	return getPermissionOverwrite(c, discord.PermissionOverwriteTypeMember, id).(*discord.MemberPermissionOverwrite)
}

func (c *GuildVoiceChannel) SetPermissionOverwrite(overwriteType discord.PermissionOverwriteType, id discord.Snowflake, allow discord.Permissions, deny discord.Permissions, opts ...rest.RequestOpt) error {
	return setPermissionOverwrite(c.Bot, c.ID(), overwriteType, id, allow, deny, opts...)
}

func (c *GuildVoiceChannel) UpdatePermissionOverwrite(overwriteType discord.PermissionOverwriteType, id discord.Snowflake, allow discord.Permissions, deny discord.Permissions, opts ...rest.RequestOpt) error {
	return updatePermissionOverwrite(c.Bot, c, overwriteType, id, allow, deny, opts...)
}

func (c *GuildVoiceChannel) DeletePermissionOverwrite(id discord.Snowflake, opts ...rest.RequestOpt) error {
	return deletePermissionOverwrite(c.Bot, c.ID(), id, opts...)
}

func (c *GuildVoiceChannel) Guild() *Guild {
	return channelGuild(c.Bot, c.GuildID())
}

func (c *GuildVoiceChannel) Parent() *GuildCategoryChannel {
	if c.ParentID == nil {
		return nil
	}
	return c.Bot.Caches.Channels().Get(*c.ParentID).(*GuildCategoryChannel)
}

func (c *GuildVoiceChannel) Connect() error {
	return c.Bot.AudioController.Connect(c.GuildID(), c.ID())
}

func (c *GuildVoiceChannel) Members() []*Member {
	return connectedMembers(c.Bot, c)
}

func (c *GuildVoiceChannel) connectedMembers() map[discord.Snowflake]struct{} {
	return c.ConnectedMemberIDs
}

var (
	_ Channel = (*GroupDMChannel)(nil)
	//_ MessageChannel = (*GroupDMChannel)(nil)
)

type GroupDMChannel struct {
	discord.GroupDMChannel
	Bot          *Bot
	RecipientIDs []discord.Snowflake
}

func (c *GroupDMChannel) set(channel Channel) Channel {
	switch ch := channel.(type) {
	case *GroupDMChannel:
		*c = *ch
		return c

	default:
		return c
	}
}

func (c *GroupDMChannel) Update(channelUpdate discord.GroupDMChannelUpdate, opts ...rest.RequestOpt) (*GroupDMChannel, error) {
	channel, err := updateChannel(c.Bot, c.ID(), channelUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return channel.(*GroupDMChannel), nil
}

func (c *GroupDMChannel) Delete(opts ...rest.RequestOpt) error {
	return deleteChannel(c.Bot, c.ID(), opts...)
}

// GetIconURL returns the Icon URL of this channel.
func (c *GroupDMChannel) GetIconURL(size int) *string {
	return discord.FormatAssetURL(route.ChannelIcon, c.ID(), c.Icon, size)
}

var (
	_ Channel      = (*GuildCategoryChannel)(nil)
	_ GuildChannel = (*GuildCategoryChannel)(nil)
)

type GuildCategoryChannel struct {
	discord.GuildCategoryChannel
	Bot *Bot
}

func (c *GuildCategoryChannel) set(channel Channel) Channel {
	switch ch := channel.(type) {
	case *GuildCategoryChannel:
		*c = *ch
		return c

	default:
		return c
	}
}

func (c *GuildCategoryChannel) UpdateGuildChannel(guildChannelUpdate discord.GuildChannelUpdate, opts ...rest.RequestOpt) (GuildChannel, error) {
	channel, err := updateChannel(c.Bot, c.ID(), guildChannelUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return channel.(GuildChannel), nil
}

// Update updates the GuildNewsChannel which can return either a GuildNewsChannel or a GuildTextChannel
func (c *GuildCategoryChannel) Update(channelUpdate discord.GuildCategoryChannelUpdate, opts ...rest.RequestOpt) (*GuildCategoryChannel, error) {
	channel, err := c.UpdateGuildChannel(channelUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return channel.(*GuildCategoryChannel), nil
}

func (c *GuildCategoryChannel) Delete(opts ...rest.RequestOpt) error {
	return deleteChannel(c.Bot, c.ChannelID, opts...)
}

func (c *GuildCategoryChannel) PermissionOverwrites() []discord.PermissionOverwrite {
	return c.ChannelPermissionOverwrites
}

func (c *GuildCategoryChannel) PermissionOverwrite(overwriteType discord.PermissionOverwriteType, id discord.Snowflake) discord.PermissionOverwrite {
	return getPermissionOverwrite(c, overwriteType, id)
}

func (c *GuildCategoryChannel) RolePermissionOverwrite(id discord.Snowflake) *discord.RolePermissionOverwrite {
	return getPermissionOverwrite(c, discord.PermissionOverwriteTypeRole, id).(*discord.RolePermissionOverwrite)
}

func (c *GuildCategoryChannel) MemberPermissionOverwrite(id discord.Snowflake) *discord.MemberPermissionOverwrite {
	return getPermissionOverwrite(c, discord.PermissionOverwriteTypeMember, id).(*discord.MemberPermissionOverwrite)
}

func (c *GuildCategoryChannel) SetPermissionOverwrite(overwriteType discord.PermissionOverwriteType, id discord.Snowflake, allow discord.Permissions, deny discord.Permissions, opts ...rest.RequestOpt) error {
	return setPermissionOverwrite(c.Bot, c.ChannelID, overwriteType, id, allow, deny, opts...)
}

func (c *GuildCategoryChannel) UpdatePermissionOverwrite(overwriteType discord.PermissionOverwriteType, id discord.Snowflake, allow discord.Permissions, deny discord.Permissions, opts ...rest.RequestOpt) error {
	return updatePermissionOverwrite(c.Bot, c, overwriteType, id, allow, deny, opts...)
}

func (c *GuildCategoryChannel) DeletePermissionOverwrite(id discord.Snowflake, opts ...rest.RequestOpt) error {
	return deletePermissionOverwrite(c.Bot, c.ChannelID, id, opts...)
}

func (c *GuildCategoryChannel) Guild() *Guild {
	return channelGuild(c.Bot, c.GuildID())
}

func (c *GuildCategoryChannel) Channels() []GuildChannel {
	channels := c.Bot.Caches.Channels().FindAll(func(channel Channel) bool {
		switch ch := channel.(type) {
		case *GuildTextChannel:
			return ch.ParentID != nil && *ch.ParentID == c.ChannelID

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
		for _, member := range channel.Members() {
			if _, ok := memberIds[member.User.ID]; ok {
				continue
			}
			members = append(members, member)
			memberIds[member.User.ID] = struct{}{}
		}
	}
	return members
}

var (
	_ Channel             = (*GuildNewsChannel)(nil)
	_ GuildChannel        = (*GuildNewsChannel)(nil)
	_ MessageChannel      = (*GuildNewsChannel)(nil)
	_ GuildMessageChannel = (*GuildNewsChannel)(nil)
)

type GuildNewsChannel struct {
	discord.GuildNewsChannel
	Bot *Bot
}

func (c *GuildNewsChannel) set(channel Channel) Channel {
	switch ch := channel.(type) {
	case *GuildNewsChannel:
		*c = *ch
		return c

	default:
		return c
	}
}

func (c *GuildNewsChannel) UpdateGuildChannel(guildChannelUpdate discord.GuildChannelUpdate, opts ...rest.RequestOpt) (GuildChannel, error) {
	channel, err := updateChannel(c.Bot, c.ID(), guildChannelUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return channel.(GuildChannel), nil
}

// Update updates the GuildNewsChannel which can return either a GuildNewsChannel or a GuildTextChannel
func (c *GuildNewsChannel) Update(channelUpdate discord.GuildNewsChannelUpdate, opts ...rest.RequestOpt) (GuildMessageChannel, error) {
	channel, err := c.UpdateGuildChannel(channelUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return channel.(GuildMessageChannel), nil
}

func (c *GuildNewsChannel) Delete(opts ...rest.RequestOpt) error {
	return deleteChannel(c.Bot, c.ID(), opts...)
}

func (c *GuildNewsChannel) PermissionOverwrites() []discord.PermissionOverwrite {
	return c.ChannelPermissionOverwrites
}

func (c *GuildNewsChannel) PermissionOverwrite(overwriteType discord.PermissionOverwriteType, id discord.Snowflake) discord.PermissionOverwrite {
	return getPermissionOverwrite(c, overwriteType, id)
}

func (c *GuildNewsChannel) RolePermissionOverwrite(id discord.Snowflake) *discord.RolePermissionOverwrite {
	return getPermissionOverwrite(c, discord.PermissionOverwriteTypeRole, id).(*discord.RolePermissionOverwrite)
}

func (c *GuildNewsChannel) MemberPermissionOverwrite(id discord.Snowflake) *discord.MemberPermissionOverwrite {
	return getPermissionOverwrite(c, discord.PermissionOverwriteTypeMember, id).(*discord.MemberPermissionOverwrite)
}

func (c *GuildNewsChannel) SetPermissionOverwrite(overwriteType discord.PermissionOverwriteType, id discord.Snowflake, allow discord.Permissions, deny discord.Permissions, opts ...rest.RequestOpt) error {
	return setPermissionOverwrite(c.Bot, c.ID(), overwriteType, id, allow, deny, opts...)
}

func (c *GuildNewsChannel) UpdatePermissionOverwrite(overwriteType discord.PermissionOverwriteType, id discord.Snowflake, allow discord.Permissions, deny discord.Permissions, opts ...rest.RequestOpt) error {
	return updatePermissionOverwrite(c.Bot, c, overwriteType, id, allow, deny, opts...)
}

func (c *GuildNewsChannel) DeletePermissionOverwrite(id discord.Snowflake, opts ...rest.RequestOpt) error {
	return deletePermissionOverwrite(c.Bot, c.ID(), id, opts...)
}

func (c *GuildNewsChannel) GetMessage(messageID discord.Snowflake, opts ...rest.RequestOpt) (*Message, error) {
	return getMessage(c.Bot, c.ID(), messageID, opts...)
}

func (c *GuildNewsChannel) GetMessages(around discord.Snowflake, before discord.Snowflake, after discord.Snowflake, limit int, opts ...rest.RequestOpt) ([]*Message, error) {
	return getMessages(c.Bot, c.ID(), around, before, after, limit, opts...)
}

func (c *GuildNewsChannel) CreateMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*Message, error) {
	return createMessage(c.Bot, c.ID(), messageCreate, opts...)
}

func (c *GuildNewsChannel) UpdateMessage(messageID discord.Snowflake, messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*Message, error) {
	return updateMessage(c.Bot, c.ID(), messageID, messageUpdate, opts...)
}

func (c *GuildNewsChannel) DeleteMessage(messageID discord.Snowflake, opts ...rest.RequestOpt) error {
	return deleteMessage(c.Bot, c.ID(), messageID, opts...)
}

func (c *GuildNewsChannel) BulkDeleteMessages(messageIDs []discord.Snowflake, opts ...rest.RequestOpt) error {
	return bulkDeleteMessages(c.Bot, c.ID(), messageIDs, opts...)
}

func (c *GuildNewsChannel) CrosspostMessage(messageID discord.Snowflake, opts ...rest.RequestOpt) (*Message, error) {
	message, err := c.Bot.RestServices.ChannelService().CrosspostMessage(c.ID(), messageID, opts...)
	if err != nil {
		return nil, err
	}
	return c.Bot.EntityBuilder.CreateMessage(*message, CacheStrategyNoWs), nil
}

func (c *GuildNewsChannel) SendTyping(opts ...rest.RequestOpt) error {
	return sendTying(c.Bot, c.ID(), opts...)
}

func (c *GuildNewsChannel) AddReaction(messageID discord.Snowflake, emoji string, opts ...rest.RequestOpt) error {
	return addReaction(c.Bot, c.ID(), messageID, emoji, opts...)
}

func (c *GuildNewsChannel) RemoveOwnReaction(messageID discord.Snowflake, emoji string, opts ...rest.RequestOpt) error {
	return removeOwnReaction(c.Bot, c.ID(), messageID, emoji, opts...)
}

func (c *GuildNewsChannel) RemoveUserReaction(messageID discord.Snowflake, emoji string, userID discord.Snowflake, opts ...rest.RequestOpt) error {
	return removeUserReaction(c.Bot, c.ID(), messageID, emoji, userID, opts...)
}

func (c *GuildNewsChannel) RemoveAllReactions(messageID discord.Snowflake, opts ...rest.RequestOpt) error {
	return removeAllReactions(c.Bot, c.ID(), messageID, opts...)
}

func (c *GuildNewsChannel) RemoveAllReactionsForEmoji(messageID discord.Snowflake, emoji string, opts ...rest.RequestOpt) error {
	return removeAllReactionsForEmoji(c.Bot, c.ID(), messageID, emoji, opts...)
}

func (c *GuildNewsChannel) GetWebhooks(opts ...rest.RequestOpt) ([]Webhook, error) {
	return getWebhooks(c.Bot, c.ID(), opts...)
}

func (c *GuildNewsChannel) CreateWebhook(webhookCreate discord.WebhookCreate, opts ...rest.RequestOpt) (Webhook, error) {
	return createWebhook(c.Bot, c.ID(), webhookCreate, opts...)
}

func (c *GuildNewsChannel) DeleteWebhook(webhookID discord.Snowflake, opts ...rest.RequestOpt) error {
	return deleteWebhook(c.Bot, webhookID, opts...)
}

func (c *GuildNewsChannel) Threads() []GuildThread {
	var threads []GuildThread
	c.Bot.Caches.Channels().ForAll(func(channel Channel) {
		if thread, ok := channel.(*GuildNewsThread); ok && thread.ParentID() == c.ID() {
			threads = append(threads, thread)
		}
	})
	return threads
}

func (c *GuildNewsChannel) NewsThreads() []*GuildNewsThread {
	var threads []*GuildNewsThread
	c.Bot.Caches.Channels().ForAll(func(channel Channel) {
		if thread, ok := channel.(*GuildNewsThread); ok && thread.ParentID() == c.ID() {
			threads = append(threads, thread)
		}
	})
	return threads
}

func (c *GuildNewsChannel) Thread(threadID discord.Snowflake) GuildThread {
	if thread, ok := c.Bot.Caches.Channels().Get(threadID).(GuildThread); ok {
		return thread
	}
	return nil
}

func (c *GuildNewsChannel) PrivateThreads() []*GuildPrivateThread {
	var threads []*GuildPrivateThread
	c.Bot.Caches.Channels().ForAll(func(channel Channel) {
		if thread, ok := channel.(*GuildPrivateThread); ok && thread.ParentID() == c.ID() {
			threads = append(threads, thread)
		}
	})
	return threads
}

func (c *GuildNewsChannel) PublicThreads() []*GuildPublicThread {
	var threads []*GuildPublicThread
	c.Bot.Caches.Channels().ForAll(func(channel Channel) {
		if thread, ok := channel.(*GuildPublicThread); ok && thread.ParentID() == c.ID() {
			threads = append(threads, thread)
		}
	})
	return threads
}

func (c *GuildNewsChannel) CreateThread(theadCreate discord.ThreadCreate, opts ...rest.RequestOpt) (GuildThread, error) {
	return createThread(c.Bot, c.ID(), theadCreate, opts...)
}

func (c *GuildNewsChannel) CreateThreadWithMessage(messageID discord.Snowflake, threadCreateWithMessage discord.ThreadCreateWithMessage, opts ...rest.RequestOpt) (GuildThread, error) {
	return createThreadWithMessage(c.Bot, c.ID(), messageID, threadCreateWithMessage, opts...)
}

func (c *GuildNewsChannel) GetPublicArchivedThreads(before discord.Time, limit int, opts ...rest.RequestOpt) ([]GuildThread, map[discord.Snowflake]*ThreadMember, bool, error) {
	return getPublicArchivedThreads(c.Bot, c.ID(), before, limit, opts...)
}

func (c *GuildNewsChannel) Guild() *Guild {
	return channelGuild(c.Bot, c.GuildID())
}

func (c *GuildNewsChannel) Parent() *GuildCategoryChannel {
	if c.ParentID == nil {
		return nil
	}
	return c.Bot.Caches.Channels().Get(*c.ParentID).(*GuildCategoryChannel)
}

func (c *GuildNewsChannel) Members() []*Member {
	return viewMembers(c.Bot, c)
}

var (
	_ Channel      = (*GuildStoreChannel)(nil)
	_ GuildChannel = (*GuildStoreChannel)(nil)
)

type GuildStoreChannel struct {
	discord.GuildStoreChannel
	Bot *Bot
}

func (c *GuildStoreChannel) set(channel Channel) Channel {
	switch ch := channel.(type) {
	case *GuildStoreChannel:
		*c = *ch
		return c

	default:
		return c
	}
}

func (c *GuildStoreChannel) UpdateGuildChannel(guildChannelUpdate discord.GuildChannelUpdate, opts ...rest.RequestOpt) (GuildChannel, error) {
	channel, err := updateChannel(c.Bot, c.ID(), guildChannelUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return channel.(GuildChannel), nil
}

// Update updates the GuildNewsChannel which can return either a GuildNewsChannel or a GuildTextChannel
func (c *GuildStoreChannel) Update(channelUpdate discord.GuildStoreChannelUpdate, opts ...rest.RequestOpt) (*GuildStoreChannel, error) {
	channel, err := c.UpdateGuildChannel(channelUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return channel.(*GuildStoreChannel), nil
}

func (c *GuildStoreChannel) Delete(opts ...rest.RequestOpt) error {
	return deleteChannel(c.Bot, c.ID(), opts...)
}

func (c *GuildStoreChannel) PermissionOverwrites() []discord.PermissionOverwrite {
	return c.ChannelPermissionOverwrites
}

func (c *GuildStoreChannel) PermissionOverwrite(overwriteType discord.PermissionOverwriteType, id discord.Snowflake) discord.PermissionOverwrite {
	return getPermissionOverwrite(c, overwriteType, id)
}

func (c *GuildStoreChannel) RolePermissionOverwrite(id discord.Snowflake) *discord.RolePermissionOverwrite {
	return getPermissionOverwrite(c, discord.PermissionOverwriteTypeRole, id).(*discord.RolePermissionOverwrite)
}

func (c *GuildStoreChannel) MemberPermissionOverwrite(id discord.Snowflake) *discord.MemberPermissionOverwrite {
	return getPermissionOverwrite(c, discord.PermissionOverwriteTypeMember, id).(*discord.MemberPermissionOverwrite)
}

func (c *GuildStoreChannel) SetPermissionOverwrite(overwriteType discord.PermissionOverwriteType, id discord.Snowflake, allow discord.Permissions, deny discord.Permissions, opts ...rest.RequestOpt) error {
	return setPermissionOverwrite(c.Bot, c.ID(), overwriteType, id, allow, deny, opts...)
}

func (c *GuildStoreChannel) UpdatePermissionOverwrite(overwriteType discord.PermissionOverwriteType, id discord.Snowflake, allow discord.Permissions, deny discord.Permissions, opts ...rest.RequestOpt) error {
	return updatePermissionOverwrite(c.Bot, c, overwriteType, id, allow, deny, opts...)
}

func (c *GuildStoreChannel) DeletePermissionOverwrite(id discord.Snowflake, opts ...rest.RequestOpt) error {
	return deletePermissionOverwrite(c.Bot, c.ID(), id, opts...)
}

func (c *GuildStoreChannel) Guild() *Guild {
	return channelGuild(c.Bot, c.GuildID())
}

func (c *GuildStoreChannel) Parent() *GuildCategoryChannel {
	if c.ParentID == nil {
		return nil
	}
	return c.Bot.Caches.Channels().Get(*c.ParentID).(*GuildCategoryChannel)
}

func (c *GuildStoreChannel) Members() []*Member {
	return viewMembers(c.Bot, c)
}

var (
	_ Channel                 = (*GuildNewsThread)(nil)
	_ GuildChannel            = (*GuildNewsThread)(nil)
	_ MessageChannel          = (*GuildNewsThread)(nil)
	_ BaseGuildMessageChannel = (*GuildNewsThread)(nil)
	_ GuildThread             = (*GuildNewsThread)(nil)
)

type GuildNewsThread struct {
	discord.GuildNewsThread
	Bot *Bot
}

func (c *GuildNewsThread) set(channel Channel) Channel {
	switch ch := channel.(type) {
	case *GuildNewsThread:
		*c = *ch
		return c

	default:
		return c
	}
}

func (c *GuildNewsThread) UpdateGuildChannel(guildChannelUpdate discord.GuildChannelUpdate, opts ...rest.RequestOpt) (GuildChannel, error) {
	channel, err := updateChannel(c.Bot, c.ID(), guildChannelUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return channel.(GuildChannel), nil
}

func (c *GuildNewsThread) Update(channelUpdate discord.GuildNewsThreadUpdate, opts ...rest.RequestOpt) (*GuildNewsThread, error) {
	channel, err := c.UpdateGuildChannel(channelUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return channel.(*GuildNewsThread), nil
}

func (c *GuildNewsThread) Delete(opts ...rest.RequestOpt) error {
	return deleteChannel(c.Bot, c.ID(), opts...)
}

func (c *GuildNewsThread) PermissionOverwrites() []discord.PermissionOverwrite {
	if parent := c.Parent(); parent != nil {
		return parent.PermissionOverwrites()
	}
	return nil
}

func (c *GuildNewsThread) PermissionOverwrite(overwriteType discord.PermissionOverwriteType, id discord.Snowflake) discord.PermissionOverwrite {
	if parent := c.Parent(); parent != nil {
		return getPermissionOverwrite(parent, overwriteType, id)
	}
	return nil
}

func (c *GuildNewsThread) RolePermissionOverwrite(id discord.Snowflake) *discord.RolePermissionOverwrite {
	if parent := c.Parent(); parent != nil {
		return getPermissionOverwrite(parent, discord.PermissionOverwriteTypeRole, id).(*discord.RolePermissionOverwrite)
	}
	return nil
}

func (c *GuildNewsThread) MemberPermissionOverwrite(id discord.Snowflake) *discord.MemberPermissionOverwrite {
	if parent := c.Parent(); parent != nil {
		return getPermissionOverwrite(parent, discord.PermissionOverwriteTypeMember, id).(*discord.MemberPermissionOverwrite)
	}
	return nil
}

func (c *GuildNewsThread) SetPermissionOverwrite(overwriteType discord.PermissionOverwriteType, id discord.Snowflake, allow discord.Permissions, deny discord.Permissions, opts ...rest.RequestOpt) error {
	return setPermissionOverwrite(c.Bot, c.ID(), overwriteType, id, allow, deny, opts...)
}

func (c *GuildNewsThread) UpdatePermissionOverwrite(overwriteType discord.PermissionOverwriteType, id discord.Snowflake, allow discord.Permissions, deny discord.Permissions, opts ...rest.RequestOpt) error {
	if parent := c.Parent(); parent != nil {
		return updatePermissionOverwrite(c.Bot, c.Parent(), overwriteType, id, allow, deny, opts...)
	}
	// TODO return error here
	return nil
}

func (c *GuildNewsThread) DeletePermissionOverwrite(id discord.Snowflake, opts ...rest.RequestOpt) error {
	return deletePermissionOverwrite(c.Bot, c.ID(), id, opts...)
}

func (c *GuildNewsThread) GetMessage(messageID discord.Snowflake, opts ...rest.RequestOpt) (*Message, error) {
	return getMessage(c.Bot, c.ID(), messageID, opts...)
}

func (c *GuildNewsThread) GetMessages(around discord.Snowflake, before discord.Snowflake, after discord.Snowflake, limit int, opts ...rest.RequestOpt) ([]*Message, error) {
	return getMessages(c.Bot, c.ID(), around, before, after, limit, opts...)
}

func (c *GuildNewsThread) CreateMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*Message, error) {
	return createMessage(c.Bot, c.ID(), messageCreate, opts...)
}

func (c *GuildNewsThread) UpdateMessage(messageID discord.Snowflake, messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*Message, error) {
	return updateMessage(c.Bot, c.ID(), messageID, messageUpdate, opts...)
}

func (c *GuildNewsThread) DeleteMessage(messageID discord.Snowflake, opts ...rest.RequestOpt) error {
	return deleteMessage(c.Bot, c.ID(), messageID, opts...)
}

func (c *GuildNewsThread) BulkDeleteMessages(messageIDs []discord.Snowflake, opts ...rest.RequestOpt) error {
	return bulkDeleteMessages(c.Bot, c.ID(), messageIDs, opts...)
}

func (c *GuildNewsThread) SendTyping(opts ...rest.RequestOpt) error {
	return sendTying(c.Bot, c.ID(), opts...)
}

func (c *GuildNewsThread) AddReaction(messageID discord.Snowflake, emoji string, opts ...rest.RequestOpt) error {
	return addReaction(c.Bot, c.ID(), messageID, emoji, opts...)
}

func (c *GuildNewsThread) RemoveOwnReaction(messageID discord.Snowflake, emoji string, opts ...rest.RequestOpt) error {
	return removeOwnReaction(c.Bot, c.ID(), messageID, emoji, opts...)
}

func (c *GuildNewsThread) RemoveUserReaction(messageID discord.Snowflake, emoji string, userID discord.Snowflake, opts ...rest.RequestOpt) error {
	return removeUserReaction(c.Bot, c.ID(), messageID, emoji, userID, opts...)
}

func (c *GuildNewsThread) RemoveAllReactions(messageID discord.Snowflake, opts ...rest.RequestOpt) error {
	return removeAllReactions(c.Bot, c.ID(), messageID, opts...)
}

func (c *GuildNewsThread) RemoveAllReactionsForEmoji(messageID discord.Snowflake, emoji string, opts ...rest.RequestOpt) error {
	return removeAllReactionsForEmoji(c.Bot, c.ID(), messageID, emoji, opts...)
}

func (c *GuildNewsThread) Guild() *Guild {
	return channelGuild(c.Bot, c.GuildID())
}

func (c *GuildNewsThread) ParentMessageChannel() GuildMessageChannel {
	return c.Bot.Caches.Channels().Get(c.ParentID()).(GuildMessageChannel)
}

func (c *GuildNewsThread) Parent() *GuildNewsChannel {
	return c.Bot.Caches.Channels().Get(c.ParentID()).(*GuildNewsChannel)
}

func (c *GuildNewsThread) Members() []*Member {
	return c.Bot.Caches.Members().FindAll(func(member *Member) bool {
		return member.ChannelPermissions(c).Has(discord.PermissionViewChannel)
	})
}

func (c *GuildNewsThread) SelfThreadMember() *ThreadMember {
	return c.ThreadMember(c.Bot.ApplicationID)
}

func (c *GuildNewsThread) ThreadMember(userID discord.Snowflake) *ThreadMember {
	return c.Bot.Caches.ThreadMembers().Get(c.ID(), userID)
}

func (c *GuildNewsThread) ThreadMembers() []*ThreadMember {
	return c.Bot.Caches.ThreadMembers().ThreadAll(c.ID())
}

func (c *GuildNewsThread) Join(opts ...rest.RequestOpt) error {
	return join(c.Bot, c.ID(), opts...)
}

func (c *GuildNewsThread) Leave(opts ...rest.RequestOpt) error {
	return leave(c.Bot, c.ID(), opts...)
}

func (c *GuildNewsThread) AddThreadMember(userID discord.Snowflake, opts ...rest.RequestOpt) error {
	return addThreadMember(c.Bot, c.ID(), userID, opts...)
}

func (c *GuildNewsThread) RemoveThreadMember(userID discord.Snowflake, opts ...rest.RequestOpt) error {
	return removeThreadMember(c.Bot, c.ID(), userID, opts...)
}

func (c *GuildNewsThread) GetThreadMember(userID discord.Snowflake, opts ...rest.RequestOpt) (*ThreadMember, error) {
	return getThreadMember(c.Bot, c.ID(), userID, opts...)
}

func (c *GuildNewsThread) GetThreadMembers(opts ...rest.RequestOpt) ([]*ThreadMember, error) {
	return getThreadMembers(c.Bot, c.ID(), opts...)
}

var (
	_ Channel                 = (*GuildPublicThread)(nil)
	_ GuildChannel            = (*GuildPublicThread)(nil)
	_ MessageChannel          = (*GuildPublicThread)(nil)
	_ BaseGuildMessageChannel = (*GuildPublicThread)(nil)
	_ GuildThread             = (*GuildPublicThread)(nil)
)

type GuildPublicThread struct {
	discord.GuildPublicThread
	Bot *Bot
}

func (c *GuildPublicThread) set(channel Channel) Channel {
	switch ch := channel.(type) {
	case *GuildPublicThread:
		*c = *ch
		return c

	default:
		return c
	}
}

func (c *GuildPublicThread) UpdateGuildChannel(guildChannelUpdate discord.GuildChannelUpdate, opts ...rest.RequestOpt) (GuildChannel, error) {
	channel, err := updateChannel(c.Bot, c.ID(), guildChannelUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return channel.(GuildChannel), nil
}

func (c *GuildPublicThread) Update(channelUpdate discord.GuildNewsThreadUpdate, opts ...rest.RequestOpt) (*GuildPublicThread, error) {
	channel, err := c.UpdateGuildChannel(channelUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return channel.(*GuildPublicThread), nil
}

func (c *GuildPublicThread) Delete(opts ...rest.RequestOpt) error {
	return deleteChannel(c.Bot, c.ID(), opts...)
}

func (c *GuildPublicThread) PermissionOverwrites() []discord.PermissionOverwrite {
	if parent := c.Parent(); parent != nil {
		return parent.PermissionOverwrites()
	}
	return nil
}

func (c *GuildPublicThread) PermissionOverwrite(overwriteType discord.PermissionOverwriteType, id discord.Snowflake) discord.PermissionOverwrite {
	if parent := c.Parent(); parent != nil {
		return getPermissionOverwrite(parent, overwriteType, id)
	}
	return nil
}

func (c *GuildPublicThread) RolePermissionOverwrite(id discord.Snowflake) *discord.RolePermissionOverwrite {
	if parent := c.Parent(); parent != nil {
		return getPermissionOverwrite(parent, discord.PermissionOverwriteTypeRole, id).(*discord.RolePermissionOverwrite)
	}
	return nil
}

func (c *GuildPublicThread) MemberPermissionOverwrite(id discord.Snowflake) *discord.MemberPermissionOverwrite {
	if parent := c.Parent(); parent != nil {
		return getPermissionOverwrite(parent, discord.PermissionOverwriteTypeMember, id).(*discord.MemberPermissionOverwrite)
	}
	return nil
}

func (c *GuildPublicThread) SetPermissionOverwrite(overwriteType discord.PermissionOverwriteType, id discord.Snowflake, allow discord.Permissions, deny discord.Permissions, opts ...rest.RequestOpt) error {
	return setPermissionOverwrite(c.Bot, c.ID(), overwriteType, id, allow, deny, opts...)
}

func (c *GuildPublicThread) UpdatePermissionOverwrite(overwriteType discord.PermissionOverwriteType, id discord.Snowflake, allow discord.Permissions, deny discord.Permissions, opts ...rest.RequestOpt) error {
	if parent := c.Parent(); parent != nil {
		return updatePermissionOverwrite(c.Bot, c.Parent(), overwriteType, id, allow, deny, opts...)
	}
	// TODO return error here
	return nil
}

func (c *GuildPublicThread) DeletePermissionOverwrite(id discord.Snowflake, opts ...rest.RequestOpt) error {
	return deletePermissionOverwrite(c.Bot, c.ID(), id, opts...)
}

func (c *GuildPublicThread) GetMessage(messageID discord.Snowflake, opts ...rest.RequestOpt) (*Message, error) {
	return getMessage(c.Bot, c.ID(), messageID, opts...)
}

func (c *GuildPublicThread) GetMessages(around discord.Snowflake, before discord.Snowflake, after discord.Snowflake, limit int, opts ...rest.RequestOpt) ([]*Message, error) {
	return getMessages(c.Bot, c.ID(), around, before, after, limit, opts...)
}

func (c *GuildPublicThread) CreateMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*Message, error) {
	return createMessage(c.Bot, c.ID(), messageCreate, opts...)
}

func (c *GuildPublicThread) UpdateMessage(messageID discord.Snowflake, messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*Message, error) {
	return updateMessage(c.Bot, c.ID(), messageID, messageUpdate, opts...)
}

func (c *GuildPublicThread) DeleteMessage(messageID discord.Snowflake, opts ...rest.RequestOpt) error {
	return deleteMessage(c.Bot, c.ID(), messageID, opts...)
}

func (c *GuildPublicThread) BulkDeleteMessages(messageIDs []discord.Snowflake, opts ...rest.RequestOpt) error {
	return bulkDeleteMessages(c.Bot, c.ID(), messageIDs, opts...)
}

func (c *GuildPublicThread) SendTyping(opts ...rest.RequestOpt) error {
	return sendTying(c.Bot, c.ID(), opts...)
}

func (c *GuildPublicThread) AddReaction(messageID discord.Snowflake, emoji string, opts ...rest.RequestOpt) error {
	return addReaction(c.Bot, c.ID(), messageID, emoji, opts...)
}

func (c *GuildPublicThread) RemoveOwnReaction(messageID discord.Snowflake, emoji string, opts ...rest.RequestOpt) error {
	return removeOwnReaction(c.Bot, c.ID(), messageID, emoji, opts...)
}

func (c *GuildPublicThread) RemoveUserReaction(messageID discord.Snowflake, emoji string, userID discord.Snowflake, opts ...rest.RequestOpt) error {
	return removeUserReaction(c.Bot, c.ID(), messageID, emoji, userID, opts...)
}

func (c *GuildPublicThread) RemoveAllReactions(messageID discord.Snowflake, opts ...rest.RequestOpt) error {
	return removeAllReactions(c.Bot, c.ID(), messageID, opts...)
}

func (c *GuildPublicThread) RemoveAllReactionsForEmoji(messageID discord.Snowflake, emoji string, opts ...rest.RequestOpt) error {
	return removeAllReactionsForEmoji(c.Bot, c.ID(), messageID, emoji, opts...)
}

func (c *GuildPublicThread) Guild() *Guild {
	return channelGuild(c.Bot, c.GuildID())
}

func (c *GuildPublicThread) Parent() *GuildTextChannel {
	return c.Bot.Caches.Channels().Get(c.ParentID()).(*GuildTextChannel)
}

func (c *GuildPublicThread) ParentMessageChannel() GuildMessageChannel {
	return c.Bot.Caches.Channels().Get(c.ParentID()).(GuildMessageChannel)
}

func (c *GuildPublicThread) Members() []*Member {
	return c.Bot.Caches.Members().FindAll(func(member *Member) bool {
		return member.ChannelPermissions(c).Has(discord.PermissionViewChannel)
	})
}

func (c *GuildPublicThread) SelfThreadMember() *ThreadMember {
	return c.ThreadMember(c.Bot.ApplicationID)
}

func (c *GuildPublicThread) ThreadMember(userID discord.Snowflake) *ThreadMember {
	return c.Bot.Caches.ThreadMembers().Get(c.ID(), userID)
}

func (c *GuildPublicThread) ThreadMembers() []*ThreadMember {
	return c.Bot.Caches.ThreadMembers().ThreadAll(c.ID())
}

func (c *GuildPublicThread) Join(opts ...rest.RequestOpt) error {
	return join(c.Bot, c.ID(), opts...)
}

func (c *GuildPublicThread) Leave(opts ...rest.RequestOpt) error {
	return leave(c.Bot, c.ID(), opts...)
}

func (c *GuildPublicThread) AddThreadMember(userID discord.Snowflake, opts ...rest.RequestOpt) error {
	return addThreadMember(c.Bot, c.ID(), userID, opts...)
}

func (c *GuildPublicThread) RemoveThreadMember(userID discord.Snowflake, opts ...rest.RequestOpt) error {
	return removeThreadMember(c.Bot, c.ID(), userID, opts...)
}

func (c *GuildPublicThread) GetThreadMember(userID discord.Snowflake, opts ...rest.RequestOpt) (*ThreadMember, error) {
	return getThreadMember(c.Bot, c.ID(), userID, opts...)
}

func (c *GuildPublicThread) GetThreadMembers(opts ...rest.RequestOpt) ([]*ThreadMember, error) {
	return getThreadMembers(c.Bot, c.ID(), opts...)
}

var (
	_ Channel                 = (*GuildPrivateThread)(nil)
	_ GuildChannel            = (*GuildPrivateThread)(nil)
	_ MessageChannel          = (*GuildPrivateThread)(nil)
	_ BaseGuildMessageChannel = (*GuildPrivateThread)(nil)
	_ GuildThread             = (*GuildPrivateThread)(nil)
)

type GuildPrivateThread struct {
	discord.GuildPrivateThread
	Bot *Bot
}

func (c *GuildPrivateThread) set(channel Channel) Channel {
	switch ch := channel.(type) {
	case *GuildPrivateThread:
		*c = *ch
		return c

	default:
		return c
	}
}

func (c *GuildPrivateThread) UpdateGuildChannel(guildChannelUpdate discord.GuildChannelUpdate, opts ...rest.RequestOpt) (GuildChannel, error) {
	channel, err := updateChannel(c.Bot, c.ID(), guildChannelUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return channel.(GuildChannel), nil
}

func (c *GuildPrivateThread) Update(channelUpdate discord.GuildNewsThreadUpdate, opts ...rest.RequestOpt) (*GuildPrivateThread, error) {
	channel, err := c.UpdateGuildChannel(channelUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return channel.(*GuildPrivateThread), nil
}

func (c *GuildPrivateThread) Delete(opts ...rest.RequestOpt) error {
	return deleteChannel(c.Bot, c.ID(), opts...)
}

func (c *GuildPrivateThread) PermissionOverwrites() []discord.PermissionOverwrite {
	if parent := c.Parent(); parent != nil {
		return parent.PermissionOverwrites()
	}
	return nil
}

func (c *GuildPrivateThread) PermissionOverwrite(overwriteType discord.PermissionOverwriteType, id discord.Snowflake) discord.PermissionOverwrite {
	if parent := c.Parent(); parent != nil {
		return getPermissionOverwrite(parent, overwriteType, id)
	}
	return nil
}

func (c *GuildPrivateThread) RolePermissionOverwrite(id discord.Snowflake) *discord.RolePermissionOverwrite {
	if parent := c.Parent(); parent != nil {
		return getPermissionOverwrite(parent, discord.PermissionOverwriteTypeRole, id).(*discord.RolePermissionOverwrite)
	}
	return nil
}

func (c *GuildPrivateThread) MemberPermissionOverwrite(id discord.Snowflake) *discord.MemberPermissionOverwrite {
	if parent := c.Parent(); parent != nil {
		return getPermissionOverwrite(parent, discord.PermissionOverwriteTypeMember, id).(*discord.MemberPermissionOverwrite)
	}
	return nil
}

func (c *GuildPrivateThread) SetPermissionOverwrite(overwriteType discord.PermissionOverwriteType, id discord.Snowflake, allow discord.Permissions, deny discord.Permissions, opts ...rest.RequestOpt) error {
	return setPermissionOverwrite(c.Bot, c.ID(), overwriteType, id, allow, deny, opts...)
}

func (c *GuildPrivateThread) UpdatePermissionOverwrite(overwriteType discord.PermissionOverwriteType, id discord.Snowflake, allow discord.Permissions, deny discord.Permissions, opts ...rest.RequestOpt) error {
	if parent := c.Parent(); parent != nil {
		return updatePermissionOverwrite(c.Bot, c.Parent(), overwriteType, id, allow, deny, opts...)
	}
	// TODO return error here
	return nil
}

func (c *GuildPrivateThread) DeletePermissionOverwrite(id discord.Snowflake, opts ...rest.RequestOpt) error {
	return deletePermissionOverwrite(c.Bot, c.ID(), id, opts...)
}

func (c *GuildPrivateThread) GetMessage(messageID discord.Snowflake, opts ...rest.RequestOpt) (*Message, error) {
	return getMessage(c.Bot, c.ID(), messageID, opts...)
}

func (c *GuildPrivateThread) GetMessages(around discord.Snowflake, before discord.Snowflake, after discord.Snowflake, limit int, opts ...rest.RequestOpt) ([]*Message, error) {
	return getMessages(c.Bot, c.ID(), around, before, after, limit, opts...)
}

func (c *GuildPrivateThread) CreateMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*Message, error) {
	return createMessage(c.Bot, c.ID(), messageCreate, opts...)
}

func (c *GuildPrivateThread) UpdateMessage(messageID discord.Snowflake, messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*Message, error) {
	return updateMessage(c.Bot, c.ID(), messageID, messageUpdate, opts...)
}

func (c *GuildPrivateThread) DeleteMessage(messageID discord.Snowflake, opts ...rest.RequestOpt) error {
	return deleteMessage(c.Bot, c.ID(), messageID, opts...)
}

func (c *GuildPrivateThread) BulkDeleteMessages(messageIDs []discord.Snowflake, opts ...rest.RequestOpt) error {
	return bulkDeleteMessages(c.Bot, c.ID(), messageIDs, opts...)
}

func (c *GuildPrivateThread) SendTyping(opts ...rest.RequestOpt) error {
	return sendTying(c.Bot, c.ID(), opts...)
}

func (c *GuildPrivateThread) AddReaction(messageID discord.Snowflake, emoji string, opts ...rest.RequestOpt) error {
	return addReaction(c.Bot, c.ID(), messageID, emoji, opts...)
}

func (c *GuildPrivateThread) RemoveOwnReaction(messageID discord.Snowflake, emoji string, opts ...rest.RequestOpt) error {
	return removeOwnReaction(c.Bot, c.ID(), messageID, emoji, opts...)
}

func (c *GuildPrivateThread) RemoveUserReaction(messageID discord.Snowflake, emoji string, userID discord.Snowflake, opts ...rest.RequestOpt) error {
	return removeUserReaction(c.Bot, c.ID(), messageID, emoji, userID, opts...)
}

func (c *GuildPrivateThread) RemoveAllReactions(messageID discord.Snowflake, opts ...rest.RequestOpt) error {
	return removeAllReactions(c.Bot, c.ID(), messageID, opts...)
}

func (c *GuildPrivateThread) RemoveAllReactionsForEmoji(messageID discord.Snowflake, emoji string, opts ...rest.RequestOpt) error {
	return removeAllReactionsForEmoji(c.Bot, c.ID(), messageID, emoji, opts...)
}

func (c *GuildPrivateThread) Guild() *Guild {
	return channelGuild(c.Bot, c.GuildID())
}

func (c *GuildPrivateThread) Parent() *GuildTextChannel {
	return c.Bot.Caches.Channels().Get(c.ParentID()).(*GuildTextChannel)
}

func (c *GuildPrivateThread) ParentMessageChannel() GuildMessageChannel {
	return c.Bot.Caches.Channels().Get(c.ParentID()).(GuildMessageChannel)
}

func (c *GuildPrivateThread) Members() []*Member {
	return c.Bot.Caches.Members().FindAll(func(member *Member) bool {
		return member.ChannelPermissions(c).Has(discord.PermissionViewChannel)
	})
}

func (c *GuildPrivateThread) SelfThreadMember() *ThreadMember {
	return c.ThreadMember(c.Bot.ApplicationID)
}

func (c *GuildPrivateThread) ThreadMember(userID discord.Snowflake) *ThreadMember {
	return c.Bot.Caches.ThreadMembers().Get(c.ID(), userID)
}

func (c *GuildPrivateThread) ThreadMembers() []*ThreadMember {
	return c.Bot.Caches.ThreadMembers().ThreadAll(c.ID())
}

func (c *GuildPrivateThread) Join(opts ...rest.RequestOpt) error {
	return join(c.Bot, c.ID(), opts...)
}

func (c *GuildPrivateThread) Leave(opts ...rest.RequestOpt) error {
	return leave(c.Bot, c.ID(), opts...)
}

func (c *GuildPrivateThread) AddThreadMember(userID discord.Snowflake, opts ...rest.RequestOpt) error {
	return addThreadMember(c.Bot, c.ID(), userID, opts...)
}

func (c *GuildPrivateThread) RemoveThreadMember(userID discord.Snowflake, opts ...rest.RequestOpt) error {
	return removeThreadMember(c.Bot, c.ID(), userID, opts...)
}

func (c *GuildPrivateThread) GetThreadMember(userID discord.Snowflake, opts ...rest.RequestOpt) (*ThreadMember, error) {
	return getThreadMember(c.Bot, c.ID(), userID, opts...)
}

func (c *GuildPrivateThread) GetThreadMembers(opts ...rest.RequestOpt) ([]*ThreadMember, error) {
	return getThreadMembers(c.Bot, c.ID(), opts...)
}

var (
	_ Channel           = (*GuildStageVoiceChannel)(nil)
	_ GuildChannel      = (*GuildStageVoiceChannel)(nil)
	_ GuildAudioChannel = (*GuildStageVoiceChannel)(nil)
)

type GuildStageVoiceChannel struct {
	discord.GuildStageVoiceChannel
	Bot                *Bot
	StageInstanceID    *discord.Snowflake
	ConnectedMemberIDs map[discord.Snowflake]struct{}
}

func (c *GuildStageVoiceChannel) set(channel Channel) Channel {
	switch ch := channel.(type) {
	case *GuildStageVoiceChannel:
		*c = *ch
		return c

	default:
		return c
	}
}

func (c *GuildStageVoiceChannel) UpdateGuildChannel(guildChannelUpdate discord.GuildChannelUpdate, opts ...rest.RequestOpt) (GuildChannel, error) {
	channel, err := updateChannel(c.Bot, c.ID(), guildChannelUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return channel.(GuildChannel), nil
}

// Update updates the GuildNewsChannel which can return either a GuildNewsChannel or a GuildTextChannel
func (c *GuildStageVoiceChannel) Update(channelUpdate discord.GuildStageVoiceChannelUpdate, opts ...rest.RequestOpt) (*GuildStageVoiceChannel, error) {
	channel, err := c.UpdateGuildChannel(channelUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return channel.(*GuildStageVoiceChannel), nil
}

func (c *GuildStageVoiceChannel) Delete(opts ...rest.RequestOpt) error {
	return deleteChannel(c.Bot, c.ID(), opts...)
}

func (c *GuildStageVoiceChannel) PermissionOverwrites() []discord.PermissionOverwrite {
	return c.ChannelPermissionOverwrites
}

func (c *GuildStageVoiceChannel) PermissionOverwrite(overwriteType discord.PermissionOverwriteType, id discord.Snowflake) discord.PermissionOverwrite {
	return getPermissionOverwrite(c, overwriteType, id)
}

func (c *GuildStageVoiceChannel) RolePermissionOverwrite(id discord.Snowflake) *discord.RolePermissionOverwrite {
	return getPermissionOverwrite(c, discord.PermissionOverwriteTypeRole, id).(*discord.RolePermissionOverwrite)
}

func (c *GuildStageVoiceChannel) MemberPermissionOverwrite(id discord.Snowflake) *discord.MemberPermissionOverwrite {
	return getPermissionOverwrite(c, discord.PermissionOverwriteTypeMember, id).(*discord.MemberPermissionOverwrite)
}

func (c *GuildStageVoiceChannel) SetPermissionOverwrite(overwriteType discord.PermissionOverwriteType, id discord.Snowflake, allow discord.Permissions, deny discord.Permissions, opts ...rest.RequestOpt) error {
	return setPermissionOverwrite(c.Bot, c.ID(), overwriteType, id, allow, deny, opts...)
}

func (c *GuildStageVoiceChannel) UpdatePermissionOverwrite(overwriteType discord.PermissionOverwriteType, id discord.Snowflake, allow discord.Permissions, deny discord.Permissions, opts ...rest.RequestOpt) error {
	return updatePermissionOverwrite(c.Bot, c, overwriteType, id, allow, deny, opts...)
}

func (c *GuildStageVoiceChannel) DeletePermissionOverwrite(id discord.Snowflake, opts ...rest.RequestOpt) error {
	return deletePermissionOverwrite(c.Bot, c.ID(), id, opts...)
}

func (c *GuildStageVoiceChannel) Guild() *Guild {
	return channelGuild(c.Bot, c.GuildID())
}

func (c *GuildStageVoiceChannel) Parent() *GuildCategoryChannel {
	if c.ParentID == nil {
		return nil
	}
	return c.Bot.Caches.Channels().Get(*c.ParentID).(*GuildCategoryChannel)
}

func (c *GuildStageVoiceChannel) Connect() error {
	return c.Bot.AudioController.Connect(c.GuildID(), c.ID())
}

func (c *GuildStageVoiceChannel) Members() []*Member {
	return connectedMembers(c.Bot, c)
}

func (c *GuildStageVoiceChannel) connectedMembers() map[discord.Snowflake]struct{} {
	return c.ConnectedMemberIDs
}

func (c *GuildStageVoiceChannel) IsModerator(member *Member) bool {
	return member.Permissions().Has(discord.PermissionsStageModerator)
}

func (c *GuildStageVoiceChannel) StageInstance() *StageInstance {
	if c.StageInstanceID == nil {
		return nil
	}
	return c.Bot.Caches.StageInstances().Get(*c.StageInstanceID)
}

func (c *GuildStageVoiceChannel) CreateStageInstance(stageInstanceCreate discord.StageInstanceCreate, opts ...rest.RequestOpt) (*StageInstance, error) {
	stageInstance, err := c.Bot.RestServices.StageInstanceService().CreateStageInstance(stageInstanceCreate, opts...)
	if err != nil {
		return nil, err
	}
	return c.Bot.EntityBuilder.CreateStageInstance(*stageInstance, CacheStrategyNoWs), nil
}

func (c *GuildStageVoiceChannel) UpdateStageInstance(stageInstanceUpdate discord.StageInstanceUpdate, opts ...rest.RequestOpt) (*StageInstance, error) {
	stageInstance, err := c.Bot.RestServices.StageInstanceService().UpdateStageInstance(c.ID(), stageInstanceUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return c.Bot.EntityBuilder.CreateStageInstance(*stageInstance, CacheStrategyNoWs), nil
}

func (c *GuildStageVoiceChannel) DeleteStageInstance(opts ...rest.RequestOpt) error {
	return c.Bot.RestServices.StageInstanceService().DeleteStageInstance(c.ID(), opts...)
}

//--------------------------------------------

func getPermissionOverwrite(channel GuildChannel, overwriteType discord.PermissionOverwriteType, id discord.Snowflake) discord.PermissionOverwrite {
	for _, overwrite := range channel.PermissionOverwrites() {
		if overwrite.Type() == overwriteType && overwrite.ID() == id {
			return overwrite
		}
	}
	return nil
}

func setPermissionOverwrite(bot *Bot, channelID discord.Snowflake, overwriteType discord.PermissionOverwriteType, id discord.Snowflake, allow discord.Permissions, deny discord.Permissions, opts ...rest.RequestOpt) error {
	var overwrite discord.PermissionOverwrite
	switch overwriteType {
	case discord.PermissionOverwriteTypeRole:
		overwrite = discord.RolePermissionOverwrite{
			RoleID: id,
			Allow:  allow,
			Deny:   deny,
		}

	case discord.PermissionOverwriteTypeMember:
		overwrite = discord.MemberPermissionOverwrite{
			UserID: id,
			Allow:  allow,
			Deny:   deny,
		}

	default:
		return errors.New("unknown permission overwrite type")
	}
	return bot.RestServices.ChannelService().UpdatePermissionOverwrite(channelID, id, overwrite, opts...)
}

func updatePermissionOverwrite(bot *Bot, channel GuildChannel, overwriteType discord.PermissionOverwriteType, id discord.Snowflake, allow discord.Permissions, deny discord.Permissions, opts ...rest.RequestOpt) error {
	var overwriteUpdate discord.PermissionOverwriteUpdate
	overwrite := getPermissionOverwrite(channel, overwriteType, id)
	switch overwriteType {
	case discord.PermissionOverwriteTypeRole:
		if overwrite != nil {
			o := overwrite.(discord.RolePermissionOverwrite)
			allow = o.Allow.Add(allow)
			deny = o.Deny.Add(deny)
		}
		overwriteUpdate = discord.RolePermissionOverwriteUpdate{
			Allow: allow,
			Deny:  deny,
		}

	case discord.PermissionOverwriteTypeMember:
		if overwrite != nil {
			o := overwrite.(discord.MemberPermissionOverwrite)
			allow = o.Allow.Add(allow)
			deny = o.Deny.Add(deny)
		}
		overwriteUpdate = discord.MemberPermissionOverwriteUpdate{
			Allow: allow,
			Deny:  deny,
		}

	default:
		return errors.New("unknown permission overwrite type")
	}

	return bot.RestServices.ChannelService().UpdatePermissionOverwrite(channel.ID(), id, overwriteUpdate, opts...)
}

func deletePermissionOverwrite(bot *Bot, channelID discord.Snowflake, id discord.Snowflake, opts ...rest.RequestOpt) error {
	return bot.RestServices.ChannelService().DeletePermissionOverwrite(channelID, id, opts...)
}

func channelGuild(bot *Bot, guildID discord.Snowflake) *Guild {
	return bot.Caches.Guilds().Get(guildID)
}

func createThread(bot *Bot, channelID discord.Snowflake, threadCreate discord.ThreadCreate, opts ...rest.RequestOpt) (GuildThread, error) {
	channel, err := bot.RestServices.ThreadService().CreateThread(channelID, threadCreate, opts...)
	if err != nil {
		return nil, err
	}
	return bot.EntityBuilder.CreateChannel(channel, CacheStrategyNo).(GuildThread), nil
}

func createThreadWithMessage(bot *Bot, channelID discord.Snowflake, messageID discord.Snowflake, threadCreateWithMessage discord.ThreadCreateWithMessage, opts ...rest.RequestOpt) (GuildThread, error) {
	channel, err := bot.RestServices.ThreadService().CreateThreadWithMessage(channelID, messageID, threadCreateWithMessage, opts...)
	if err != nil {
		return nil, err
	}
	return bot.EntityBuilder.CreateChannel(channel, CacheStrategyNo).(GuildThread), nil
}

func join(bot *Bot, threadID discord.Snowflake, opts ...rest.RequestOpt) error {
	return bot.RestServices.ThreadService().JoinThread(threadID, opts...)
}

func leave(bot *Bot, threadID discord.Snowflake, opts ...rest.RequestOpt) error {
	return bot.RestServices.ThreadService().LeaveThread(threadID, opts...)
}

func addThreadMember(bot *Bot, threadID discord.Snowflake, userID discord.Snowflake, opts ...rest.RequestOpt) error {
	return bot.RestServices.ThreadService().AddThreadMember(threadID, userID, opts...)
}

func removeThreadMember(bot *Bot, threadID discord.Snowflake, userID discord.Snowflake, opts ...rest.RequestOpt) error {
	return bot.RestServices.ThreadService().RemoveThreadMember(threadID, userID, opts...)
}

func getThreadMember(bot *Bot, threadID discord.Snowflake, userID discord.Snowflake, opts ...rest.RequestOpt) (*ThreadMember, error) {
	threadMember, err := bot.RestServices.ThreadService().GetThreadMember(threadID, userID, opts...)
	if err != nil {
		return nil, err
	}
	return bot.EntityBuilder.CreateThreadMember(*threadMember, CacheStrategyNo), nil
}

func getThreadMembers(bot *Bot, threadID discord.Snowflake, opts ...rest.RequestOpt) ([]*ThreadMember, error) {
	members, err := bot.RestServices.ThreadService().GetThreadMembers(threadID, opts...)
	if err != nil {
		return nil, err
	}
	threadMembers := make([]*ThreadMember, len(members))
	for i := range members {
		threadMembers[i] = bot.EntityBuilder.CreateThreadMember(members[i], CacheStrategyNo)
	}
	return threadMembers, nil
}

func getPublicArchivedThreads(bot *Bot, channelID discord.Snowflake, before discord.Time, limit int, opts ...rest.RequestOpt) ([]GuildThread, map[discord.Snowflake]*ThreadMember, bool, error) {
	getThreads, err := bot.RestServices.ThreadService().GetPublicArchivedThreads(channelID, before, limit, opts...)
	if err != nil {
		return nil, nil, false, err
	}

	threads := make([]GuildThread, len(getThreads.Threads))
	for i := range getThreads.Threads {
		threads[i] = bot.EntityBuilder.CreateChannel(getThreads.Threads[i], CacheStrategyNo).(GuildThread)
	}

	return threads, createThreadMembers(bot, getThreads.Members), getThreads.HasMore, nil
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

func sendTying(bot *Bot, channelID discord.Snowflake, opts ...rest.RequestOpt) error {
	return bot.RestServices.ChannelService().SendTyping(channelID, opts...)
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

func deleteChannel(bot *Bot, channelID discord.Snowflake, opts ...rest.RequestOpt) error {
	return bot.RestServices.ChannelService().DeleteChannel(channelID, opts...)
}

func addReaction(bot *Bot, channelID discord.Snowflake, messageID discord.Snowflake, emoji string, opts ...rest.RequestOpt) error {
	return bot.RestServices.ChannelService().AddReaction(channelID, messageID, emoji, opts...)
}

func removeOwnReaction(bot *Bot, channelID discord.Snowflake, messageID discord.Snowflake, emoji string, opts ...rest.RequestOpt) error {
	return bot.RestServices.ChannelService().RemoveOwnReaction(channelID, messageID, emoji, opts...)
}

func removeUserReaction(bot *Bot, channelID discord.Snowflake, messageID discord.Snowflake, emoji string, userID discord.Snowflake, opts ...rest.RequestOpt) error {
	return bot.RestServices.ChannelService().RemoveUserReaction(channelID, messageID, emoji, userID, opts...)
}

func removeAllReactions(bot *Bot, channelID discord.Snowflake, messageID discord.Snowflake, opts ...rest.RequestOpt) error {
	return bot.RestServices.ChannelService().RemoveAllReactions(channelID, messageID, opts...)
}

func removeAllReactionsForEmoji(bot *Bot, channelID discord.Snowflake, messageID discord.Snowflake, emoji string, opts ...rest.RequestOpt) error {
	return bot.RestServices.ChannelService().RemoveAllReactionsForEmoji(channelID, messageID, emoji, opts...)
}

func getWebhooks(bot *Bot, channelID discord.Snowflake, opts ...rest.RequestOpt) ([]Webhook, error) {
	webhooks, err := bot.RestServices.ChannelService().GetWebhooks(channelID, opts...)
	if err != nil {
		return nil, err
	}
	coreWebhooks := make([]Webhook, len(webhooks))
	for i := range webhooks {
		coreWebhooks[i] = bot.EntityBuilder.CreateWebhook(webhooks[i], CacheStrategyNoWs)
	}
	return coreWebhooks, nil
}

func createWebhook(bot *Bot, channelID discord.Snowflake, webhookCreate discord.WebhookCreate, opts ...rest.RequestOpt) (Webhook, error) {
	webhook, err := bot.RestServices.ChannelService().CreateWebhook(channelID, webhookCreate, opts...)
	if err != nil {
		return nil, err
	}
	return bot.EntityBuilder.CreateWebhook(webhook, CacheStrategyNoWs), nil
}

func deleteWebhook(bot *Bot, webhookID discord.Snowflake, opts ...rest.RequestOpt) error {
	return bot.RestServices.WebhookService().DeleteWebhook(webhookID, opts...)
}

func viewMembers(bot *Bot, guildChannel GuildChannel) []*Member {
	return bot.Caches.Members().FindAll(func(member *Member) bool {
		return member.ChannelPermissions(guildChannel).Has(discord.PermissionViewChannel)
	})
}

func connectedMembers(bot *Bot, audioChannel GuildAudioChannel) []*Member {
	return bot.Caches.Members().FindAll(func(member *Member) bool {
		_, ok := audioChannel.connectedMembers()[member.User.ID]
		return ok
	})
}

func LastPinTimestamp(channel MessageChannel) *discord.Time {
	if channel == nil {
		return nil
	}
	switch ch := channel.(type) {
	case *GuildTextChannel:
		return ch.LastPinTimestamp

	case *DMChannel:
		return ch.LastPinTimestamp

	case *GuildNewsChannel:
		return ch.LastPinTimestamp

	case *GuildNewsThread:
		return ch.LastPinTimestamp

	case *GuildPrivateThread:
		return ch.LastPinTimestamp

	case *GuildPublicThread:
		return ch.LastPinTimestamp

	default:
		panic("unknown channel type")
	}
}
