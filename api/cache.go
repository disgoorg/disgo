package api

import "time"

// MemberCachePolicy can be used to define your own policy for caching members
type MemberCachePolicy func(*Member) bool

// Default member cache policies
var (
	MemberCachePolicyNone    MemberCachePolicy = func(_ *Member) bool { return false }
	MemberCachePolicyAll     MemberCachePolicy = func(_ *Member) bool { return true }
	MemberCachePolicyOwner   MemberCachePolicy = func(member *Member) bool { return member.isOwner() }
	MemberCachePolicyOnline  MemberCachePolicy = func(_ *Member) bool { return false }
	MemberCachePolicyVoice   MemberCachePolicy = func(member *Member) bool { return false }
	MemberCachePolicyPending MemberCachePolicy = func(member *Member) bool { return member.IsPending }
	MemberCachePolicyDefault                   = MemberCachePolicyOwner.Or(MemberCachePolicyVoice)
)

// Or allows you to combine that policy with another, meaning either needs to be true
func (p MemberCachePolicy) Or(policy MemberCachePolicy) MemberCachePolicy {
	return func(member *Member) bool {
		return p(member) || policy(member)
	}
}

// And allows you to require both policies to be true for the member to be cached
func (p MemberCachePolicy) And(policy MemberCachePolicy) MemberCachePolicy {
	return func(member *Member) bool {
		return p(member) && policy(member)
	}
}

// MemberCachePolicyAny is a shorthand for MemberCachePolicy.Or(MemberCachePolicy).Or(MemberCachePolicy) etc.
func MemberCachePolicyAnyOf(policy MemberCachePolicy, policies ...MemberCachePolicy) MemberCachePolicy {
	for _, p := range policies {
		policy = policy.Or(p)
	}
	return policy
}

// MemberCachePolicyAll is a shorthand for MemberCachePolicy.And(MemberCachePolicy).And(MemberCachePolicy) etc.
func MemberCachePolicyAllOf(policy MemberCachePolicy, policies ...MemberCachePolicy) MemberCachePolicy {
	for _, p := range policies {
		policy = policy.And(p)
	}
	return policy
}

// MessageCachePolicy can be used to define your own policy for caching messages
type MessageCachePolicy func(*Message) bool

// Default member cache policies
var (
	MessageCachePolicyNone    MessageCachePolicy = func(_ *Message) bool { return false }
	MessageCachePolicyDefault                    = MessageCachePolicyNone
)

// Or allows you to combine that policy with another, meaning either needs to be true
func (p MessageCachePolicy) Or(policy MessageCachePolicy) MessageCachePolicy {
	return func(message *Message) bool {
		return p(message) || policy(message)
	}
}

// And allows you to require both policies to be true for the member to be cached
func (p MessageCachePolicy) And(policy MessageCachePolicy) MessageCachePolicy {
	return func(message *Message) bool {
		return p(message) && policy(message)
	}
}

// MessageCachePolicyDuration creates a new MessageCachePolicy which caches messages for the give duration
func MessageCachePolicyDuration(duration time.Duration) MessageCachePolicy {
	return func(message *Message) bool {
		return message.CreatedAt.Add(duration).After(time.Now())
	}
}

// MessageCachePolicyAny is a shorthand for MessageCachePolicy.Or(MessageCachePolicy).Or(MessageCachePolicy) etc.
func MessageCachePolicyAny(policy MessageCachePolicy, policies ...MessageCachePolicy) MessageCachePolicy {
	for _, p := range policies {
		policy = policy.Or(p)
	}
	return policy
}

// MessageCachePolicyAll is a shorthand for MessageCachePolicy.And(MessageCachePolicy).And(MessageCachePolicy) etc.
func MessageCachePolicyAll(policy MessageCachePolicy, policies ...MessageCachePolicy) MessageCachePolicy {
	for _, p := range policies {
		policy = policy.And(p)
	}
	return policy
}

// Cache allows you to access the objects that are stored in-memory by Discord
type Cache interface {
	Guild(Snowflake) *Guild
	GuildsByName(string, bool) []*Guild
	Guilds() []*Guild
	GuildCache() map[Snowflake]*Guild
	CacheGuild(*Guild)
	UncacheGuild(Snowflake)

	UnavailableGuild(Snowflake) *UnavailableGuild
	CacheUnavailableGuild(*UnavailableGuild)
	UncacheUnavailableGuild(Snowflake)

	/*Message(Snowflake) *Message
	Messages(Snowflake) []*Message
	AllMessages() []*Message
	MessageCache(Snowflake) map[Snowflake]*Message
	AllMessageCache() map[Snowflake]map[Snowflake]*Message
	CacheMessage(*Message)
	UncacheMessage(Snowflake)*/

	User(Snowflake) *User
	UserByTag(string) *User
	UsersByName(string, bool) []*User
	Users() []*User
	UserCache() map[Snowflake]*User
	CacheUser(*User)
	UncacheUser(Snowflake)
	FindUser(func(*User) bool) *User
	FindUsers(func(*User) bool) []*User

	Member(Snowflake, Snowflake) *Member
	MemberByTag(Snowflake, string) *Member
	MembersByName(Snowflake, string, bool) []*Member
	Members(Snowflake) []*Member
	AllMembers() []*Member
	MemberCache(Snowflake) map[Snowflake]*Member
	AllMemberCache() map[Snowflake]map[Snowflake]*Member
	CacheMember(*Member)
	UncacheMember(Snowflake, Snowflake)
	FindMember(Snowflake, func(*Member) bool) *Member
	FindMembers(Snowflake, func(*Member) bool) []*Member

	Role(Snowflake, Snowflake) *Role
	RolesByName(Snowflake, string, bool) []*Role
	Roles(Snowflake) []*Role
	AllRoles() []*Role
	RoleCache(Snowflake) map[Snowflake]*Role
	AllRoleCache() map[Snowflake]map[Snowflake]*Role
	CacheRole(*Role)
	UncacheRole(Snowflake, Snowflake)
	FindRole(Snowflake, func(*Role) bool) *Role
	FindRoles(Snowflake, func(*Role) bool) []*Role

	DMChannel(Snowflake) *DMChannel
	DMChannels() []*DMChannel
	DMChannelCache() map[Snowflake]*DMChannel
	CacheDMChannel(*DMChannel)
	UncacheDMChannel(Snowflake)
	FindDMChannel(func(*DMChannel) bool) *DMChannel
	FindDMChannels(func(*DMChannel) bool) []*DMChannel

	Channel(Snowflake) *Channel
	MessageChannel(Snowflake) *MessageChannel
	GuildChannel(Snowflake) *GuildChannel

	TextChannel(Snowflake) *TextChannel
	TextChannelsByName(Snowflake, string, bool) []*TextChannel
	TextChannels(Snowflake) []*TextChannel
	TextChannelCache(Snowflake) map[Snowflake]*TextChannel
	CacheTextChannel(*TextChannel)
	UncacheTextChannel(Snowflake, Snowflake)
	FindTextChannel(Snowflake, func(*TextChannel) bool) *TextChannel
	FindTextChannels(Snowflake, func(*TextChannel) bool) []*TextChannel

	StoreChannel(Snowflake) *StoreChannel
	StoreChannelsByName(Snowflake, string, bool) []*StoreChannel
	StoreChannels(Snowflake) []*StoreChannel
	StoreChannelCache(Snowflake) map[Snowflake]*StoreChannel
	CacheStoreChannel(*StoreChannel)
	UncacheStoreChannel(Snowflake, Snowflake)
	FindStoreChannel(Snowflake, func(*StoreChannel) bool) *StoreChannel
	FindStoreChannels(Snowflake, func(*StoreChannel) bool) []*StoreChannel

	VoiceChannel(Snowflake) *VoiceChannel
	VoiceChannelsByName(Snowflake, string, bool) []*VoiceChannel
	VoiceChannels(Snowflake) []*VoiceChannel
	VoiceChannelCache(Snowflake) map[Snowflake]*VoiceChannel
	CacheVoiceChannel(*VoiceChannel)
	UncacheVoiceChannel(Snowflake, Snowflake)
	FindVoiceChannel(Snowflake, func(*VoiceChannel) bool) *VoiceChannel
	FindVoiceChannels(Snowflake, func(*VoiceChannel) bool) []*VoiceChannel

	Category(Snowflake) *CategoryChannel
	CategoriesByName(Snowflake, string, bool) []*CategoryChannel
	Categories(Snowflake) []*CategoryChannel
	AllCategories() []*CategoryChannel
	CategoryCache(Snowflake) map[Snowflake]*CategoryChannel
	AllCategoryCache() map[Snowflake]map[Snowflake]*CategoryChannel
	CacheCategory(*CategoryChannel)
	UncacheCategory(Snowflake, Snowflake)
	FindCategory(Snowflake, func(*CategoryChannel) bool) *CategoryChannel
	FindCategories(Snowflake, func(*CategoryChannel) bool) []*CategoryChannel

	/*Emote(Snowflake) *Emote
	EmotesByName(string, bool) []*Emote
	Emotes() []*Emote
	EmoteCache() map[Snowflake]*Emote
	CacheEmote(*Emote)
	UncacheEmote(Snowflake)*/
}
