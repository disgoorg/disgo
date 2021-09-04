package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

type Caches interface {
	Disgo() Disgo
	Close()
	DoCleanup()
	Config() CacheConfig

	UserCache() UserCache
	RoleCache() RoleCache
	MemberCache() MemberCache
	VoiceStateCache() VoiceStateCache
	MessageCache() MessageCache
	EmojiCache() EmojiCache
	GuildCache() GuildCache
	ChannelCache() ChannelCache
	TextChannelCache() TextChannelCache
	VoiceChannelCache() VoiceChannelCache
	DMChannelCache() DMChannelCache
	NewsChannelCache() NewsChannelCache
	CategoryCache() CategoryCache
	StoreChannelCache() StoreChannelCache
	StageChannelCache() StageChannelCache
	StageInstanceCache() StageInstanceCache
}

type BaseCache interface {
	Disgo() Disgo
	DoCleanup()
}

type UserCache interface {
	BaseCache
	Get(userID discord.Snowflake) *User
	GetCopy(id discord.Snowflake) *User
	GetFirstByTag(string) *User
	GetAllByName(string, bool) []*User
	All() []*User
	UserCache() map[discord.Snowflake]*User
	Cache(*User) *User
	UncacheUser(discord.Snowflake)
	FindUser(func(*User) bool) *User
	FindUsers(func(*User) bool) []*User
}

type RoleCache interface {
	BaseCache
	Get(roleID discord.Snowflake) *Role
	GetCopy(id discord.Snowflake) *Role
	RolesByName(discord.Snowflake, string, bool) []*Role
	All(discord.Snowflake) []*Role
	AllRoles() []*Role
	RoleCache(discord.Snowflake) map[discord.Snowflake]*Role
	AllRoleCache() map[discord.Snowflake]map[discord.Snowflake]*Role
	Cache(*Role) *Role
	Uncache(discord.Snowflake, discord.Snowflake)
	FindRole(discord.Snowflake, func(*Role) bool) *Role
	FindRoles(discord.Snowflake, func(*Role) bool) []*Role
}

type MemberCache interface {
	BaseCache
	Get(guildId discord.Snowflake, userID discord.Snowflake) *Member
	GetCopy(guildId discord.Snowflake, userID discord.Snowflake) *Member
	MemberByTag(discord.Snowflake, string) *Member
	MembersByName(discord.Snowflake, string, bool) []*Member
	Members(discord.Snowflake) []*Member
	AllMembers() []*Member
	MemberCache(discord.Snowflake) map[discord.Snowflake]*Member
	AllMemberCache() map[discord.Snowflake]map[discord.Snowflake]*Member
	Cache(member *Member) *Member
	Uncache(discord.Snowflake, discord.Snowflake)
	FindMember(discord.Snowflake, func(*Member) bool) *Member
	FindMembers(discord.Snowflake, func(*Member) bool) []*Member
	CachePolicy(m *Message) MessageCachePolicy
}

type VoiceStateCache interface {
	BaseCache
	Get(guildID discord.Snowflake, userID discord.Snowflake) *VoiceState
	GetCopy(id discord.Snowflake) *VoiceState
	VoiceStates(guildID discord.Snowflake) []*VoiceState
	VoiceStateCache(guildID discord.Snowflake) map[discord.Snowflake]*VoiceState
	Cache(voiceState *VoiceState) *VoiceState
	Uncache(guildID discord.Snowflake, userID discord.Snowflake)
}

type MessageCache interface {
	BaseCache
	Get(channelID discord.Snowflake, messageID discord.Snowflake) *Message
	GetCopy(channelID discord.Snowflake, messageID discord.Snowflake) *Message
	GetAll(channelID discord.Snowflake) []*Message
	GetCache(channelID discord.Snowflake) map[discord.Snowflake]*Message
	GetAllCache() map[discord.Snowflake]map[discord.Snowflake]*Message
	Cache(message *Message) *Message
	Uncache(channelID discord.Snowflake, messageID discord.Snowflake)
}

type EmojiCache interface {
	BaseCache
	Get(EmojiID discord.Snowflake) *Emoji
	GetCopy(id discord.Snowflake) *Emoji
	EmojisByName(guildID discord.Snowflake, name string, ignoreCase bool) []*Emoji
	Emojis(guildID discord.Snowflake) []*Emoji
	EmojiCache(guildID discord.Snowflake) map[discord.Snowflake]*Emoji
	AllEmojiCache() map[discord.Snowflake]map[discord.Snowflake]*Emoji
	Cache(*Emoji) *Emoji
	Uncache(guildID discord.Snowflake, emojiID discord.Snowflake)
}

type ChannelCache interface {
	BaseCache
	GetChannel(id discord.Snowflake) Channel
	GetChannelCopy(id discord.Snowflake) Channel

	GetMessageChannel(id discord.Snowflake) MessageChannel
	GetMessageChannelCopy(id discord.Snowflake) MessageChannel

	GetGuildChannel(id discord.Snowflake) GuildChannel
	GetGuildChannelCopy(id discord.Snowflake) GuildChannel
}

type TextChannelCache interface {
	BaseCache
	Get(discord.Snowflake) TextChannel
	GetCopy(id discord.Snowflake) TextChannel
	GetByName(discord.Snowflake, string, bool) []TextChannel
	GetAll(discord.Snowflake) []TextChannel
	GetCache(discord.Snowflake) map[discord.Snowflake]TextChannel
	Cache(TextChannel) TextChannel
	Uncache(discord.Snowflake, discord.Snowflake)
	FindFirst(discord.Snowflake, func(TextChannel) bool) TextChannel
	FindAll(discord.Snowflake, func(TextChannel) bool) []TextChannel
}

type VoiceChannelCache interface {
	BaseCache
	Get(discord.Snowflake) VoiceChannel
	GetCopy(id discord.Snowflake) VoiceChannel
	GetByName(discord.Snowflake, string, bool) []VoiceChannel
	GetAll(discord.Snowflake) []VoiceChannel
	GetCache(discord.Snowflake) map[discord.Snowflake]VoiceChannel
	Cache(VoiceChannel) VoiceChannel
	Uncache(discord.Snowflake, discord.Snowflake)
	FindFirst(discord.Snowflake, func(VoiceChannel) bool) VoiceChannel
	FindAll(discord.Snowflake, func(VoiceChannel) bool) []VoiceChannel
}

type DMChannelCache interface {
	BaseCache
	Get(id discord.Snowflake) DMChannel
	GetCopy(id discord.Snowflake) DMChannel
	GetAll() []DMChannel
	GetCache() map[discord.Snowflake]DMChannel
	Cache(DMChannel) DMChannel
	Uncache(dmChannelID discord.Snowflake)
	FindFirst(func(DMChannel) bool) DMChannel
	FindAll(func(DMChannel) bool) []DMChannel
}

type NewsChannelCache interface {
	BaseCache
	Get(discord.Snowflake) NewsChannel
	GetCopy(id discord.Snowflake) NewsChannel
	GetByName(discord.Snowflake, string, bool) []NewsChannel
	GetAll(discord.Snowflake) []NewsChannel
	GetCache(discord.Snowflake) map[discord.Snowflake]NewsChannel
	Cache(NewsChannel) NewsChannel
	Uncache(discord.Snowflake, discord.Snowflake)
	FindFirst(discord.Snowflake, func(NewsChannel) bool) NewsChannel
	FindAll(discord.Snowflake, func(NewsChannel) bool) []NewsChannel
}

type CategoryCache interface {
	BaseCache
	Get(discord.Snowflake) Category
	GetCopy(id discord.Snowflake) Category
	GetByName(discord.Snowflake, string, bool) []Category
	GetAll() []Category
	GetCache() map[discord.Snowflake]map[discord.Snowflake]Category
	Cache(Category) Category
	Uncache(discord.Snowflake, discord.Snowflake)
	FindFirst(discord.Snowflake, func(Category) bool) Category
	FindAll(discord.Snowflake, func(Category) bool) []Category
}

type StoreChannelCache interface {
	BaseCache
	Get(discord.Snowflake) StoreChannel
	GetCopy(id discord.Snowflake) StoreChannel
	GetByName(discord.Snowflake, string, bool) []StoreChannel
	GetAll() []Category
	GetCache() map[discord.Snowflake]map[discord.Snowflake]StoreChannel
	Cache(StoreChannel) StoreChannel
	Uncache(discord.Snowflake, discord.Snowflake)
	FindFirst(discord.Snowflake, func(StoreChannel) bool) StoreChannel
	FindAll(discord.Snowflake, func(StoreChannel) bool) []StoreChannel
}

type StageChannelCache interface {
	BaseCache
	Get(discord.Snowflake) StageChannel
	GetCopy(id discord.Snowflake) StageChannel
	GetByName(discord.Snowflake, string, bool) []StageChannel
	GetAll() []StageChannelCache
	GetCache() map[discord.Snowflake]map[discord.Snowflake]StageChannel
	Cache(StageChannel) StageChannel
	Uncache(discord.Snowflake, discord.Snowflake)
	FindFirst(discord.Snowflake, func(StageChannel) bool) StageChannel
	FindAll(discord.Snowflake, func(StageChannel) bool) []StageChannel
}

type StageInstanceCache interface {
	BaseCache
	Get(id discord.Snowflake) *StageInstance
	GetCopy(id discord.Snowflake) StageInstance
	GetByName(guildID discord.Snowflake, name string, ignoreCase bool) *StageInstance
	GetAll() []StageInstance
	GetCache() map[discord.Snowflake]map[discord.Snowflake]StageInstance
	Cache(*StageInstance) *StageInstance
	Uncache(discord.Snowflake, discord.Snowflake)
	FindFirst(discord.Snowflake, func(StageInstance) bool) *StageInstance
	FindAll(discord.Snowflake, func(StageInstance) bool) []StageInstance
}
