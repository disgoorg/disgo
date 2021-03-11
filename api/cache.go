package api

type MemberCachePolicy func(Member) bool

var (
	MemberCachePolicyNone    MemberCachePolicy = func(_ Member) bool { return false }
	MemberCachePolicyOwner   MemberCachePolicy = func(member Member) bool { return member.isOwner() }
	MemberCachePolicyONLINE  MemberCachePolicy = func(_ Member) bool { return false }
	MemberCachePolicyVOICE   MemberCachePolicy = func(member Member) bool { return false }
	MemberCachePolicyPENDING MemberCachePolicy = func(member Member) bool { return member.IsPending }
	MemberCachePolicyDEFAULT MemberCachePolicy = MemberCachePolicyOwner.or(MemberCachePolicyVOICE)
)

func (p MemberCachePolicy) or(policy MemberCachePolicy) MemberCachePolicy {
	return func(member Member) bool {
		return p(member) || policy(member)
	}
}

func (p MemberCachePolicy) and(policy MemberCachePolicy) MemberCachePolicy {
	return func(member Member) bool {
		return p(member) && policy(member)
	}
}

func MemberCachePolicyAny(policy MemberCachePolicy, policies ...MemberCachePolicy) MemberCachePolicy {
	for _, p := range policies {
		policy = policy.or(p)
	}
	return policy
}

func MemberCachePolicyAll(policy MemberCachePolicy, policies ...MemberCachePolicy) MemberCachePolicy {
	for _, p := range policies {
		policy = policy.and(p)
	}
	return policy
}

type Cache interface {
	Guild(Snowflake) *Guild
	GuildsByName(string, bool) []*Guild
	Guilds() []*Guild
	GuildCache() map[Snowflake]*Guild
	CacheGuild(*Guild)
	UncacheGuild(Snowflake)

	User(Snowflake) *User
	UserByTag(string) *User
	UsersByName(Snowflake, bool) []*User
	Users() []*User
	UserCache() map[Snowflake]*User
	CacheUser(*User)
	UncacheUser(Snowflake)

	Member(Snowflake, Snowflake) *Member
	MemberByTag(Snowflake, string) *Member
	MembersByName(Snowflake, string, bool) []*Member
	Members(Snowflake) []*Member
	AllMembers() []*Member
	MemberCache(Snowflake) map[Snowflake]*Member
	AllMemberCache(Snowflake) map[Snowflake]map[Snowflake]*Member
	CacheMember(*Member)
	UncacheMember(Snowflake, Snowflake)

	Role(Snowflake, Snowflake) *Role
	RolesByName(Snowflake, string, bool) []*Role
	Roles(Snowflake) []*Role
	AllRoles() []*Role
	RoleCache(Snowflake) map[Snowflake]*Role
	AllRoleCache(Snowflake) map[Snowflake]map[Snowflake]*Role
	CacheRole(*Role)
	UncacheRole(Snowflake, Snowflake)

	Channel(Snowflake) *Channel
	ChannelsByName(string, bool) []*Channel
	Channels() []*Channel
	ChannelCache() map[Snowflake]*Channel
	CacheChannel(*Channel)
	UncacheChannel(Snowflake)

	DMChannel(Snowflake) *DMChannel
	DMChannelsByName(string, bool) []*DMChannel
	DMChannels() []*DMChannel
	DMChannelCache() map[Snowflake]*DMChannel
	CacheDMChannel(*DMChannel)
	UncacheDMChannel(Snowflake)

	GuildChannel(Snowflake) *GuildChannel
	GuildChannelsByName(string, bool) []*GuildChannel
	GuildChannels() []*GuildChannel
	GuildChannelCache() map[Snowflake]*GuildChannel
	CacheGuildChannel(*GuildChannel)
	UncacheGuildChannel(Snowflake)

	TextChannel(Snowflake) TextChannel
	TextChannelsByName(string, bool) []*TextChannel
	TextChannels() []*TextChannel
	TextChannelCache() map[Snowflake]*TextChannel
	CacheTextChannel(*TextChannel)
	UncacheTextChannel(Snowflake)

	NewsChannel(Snowflake) *NewsChannel
	NewsChannelsByName(string, bool) []*NewsChannel
	NewsChannels() []*NewsChannel
	NewsChannelCache() map[Snowflake]*NewsChannel
	CacheNewsChannel(*NewsChannel)
	UncacheNewsChannel(Snowflake)

	StoreChannel(Snowflake) *StoreChannel
	StoreChannelsByName(string, bool) []*StoreChannel
	StoreChannels() []*StoreChannel
	StoreChannelCache() map[Snowflake]*StoreChannel
	CacheStoreChannel(*StoreChannel)
	UncacheStoreChannel(Snowflake)

	VoiceChannel(Snowflake) *VoiceChannel
	VoiceChannelsByName(string, bool) []*VoiceChannel
	VoiceChannels() []*VoiceChannel
	VoiceChannelCache() map[Snowflake]*VoiceChannel
	CacheVoiceChannel(*VoiceChannel)
	UncacheVoiceChannel(Snowflake)

	Category(Snowflake) *CategoryChannel
	CategoriesByName(string, bool) []*CategoryChannel
	Categories() []*CategoryChannel
	CategoryCache() map[Snowflake]*CategoryChannel
	CacheCategory(*CategoryChannel)
	UncacheCategory(Snowflake)

	Emote(Snowflake) *Emote
	EmotesByName(string, bool) []*Emote
	Emotes() []*Emote
	EmoteCache() map[Snowflake]*Emote
	CacheEmote(*Emote)
	UncacheEmote(Snowflake)
}
