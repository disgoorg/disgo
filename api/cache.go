package api

type MemberCachePolicy func(Member) bool

var (
	MemberCachePolicyNone MemberCachePolicy = func(_ Member) bool {return false}
	MemberCachePolicyOwner MemberCachePolicy = func(member Member) bool {return member.isOwner()}
	MemberCachePolicyONLINE MemberCachePolicy = func(_ Member) bool {return false }
	MemberCachePolicyVOICE MemberCachePolicy = func(member Member) bool { return false }
	MemberCachePolicyPENDING MemberCachePolicy = func(member Member) bool {return member.IsPending}
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
	GetGuildById(Snowflake) Guild
	GetGuildsByName(string, bool) []Guild
	GetGuildsCache() map[Snowflake]Guild
	GetGuilds() []Guild
	CacheGuild(Guild)

	GetUserById(Snowflake) User
	GetUsersByName(Snowflake, bool) []User
	GetUsersCache() map[Snowflake]User
	GetUsers() []User
	CacheUser(User)

	GetMemberById(Snowflake) Member
	GetMemberByName(string, bool) []Member
	GetMembersCache() map[Snowflake]Member
	GetMembers() []Member
	CacheMember(Member)

	GetChannelById(Snowflake) Channel
	GetChannelsByName(string, bool) []Channel
	GetChannelsCache() map[Snowflake]Channel
	GetChannels() []Channel
	CacheChannel(Channel)

	GetPrivateChannelById(Snowflake) Channel
	GetPrivateChannelsByName(string, bool) []Channel // not sure if we need this lul
	GetPrivateChannelsCache() map[Snowflake]Channel
	GetPrivateChannels() []Channel
	CachePrivateChannel(Channel)

	GetGuildChannelById(Snowflake) GuildChannel
	GetGuildChannelsByName(string, bool) []GuildChannel
	GetGuildChannelsCache() map[Snowflake]GuildChannel
	GetGuildChannels() []GuildChannel
	CacheGuildChannel(GuildChannel)

	GetTextChannelById(Snowflake) TextChannel
	GetTextChannelsByName(string, bool) []CategoryChannel
	GetTextChannelsCache() map[Snowflake]TextChannel
	GetTextChannels() []TextChannel
	CacheTextChannel(TextChannel)

	GetNewsChannelById(Snowflake) NewsChannel
	GetNewsChannelsByName(string, bool) []NewsChannel
	GetNewsChannelsCache() map[Snowflake]NewsChannel
	GetNewsChannels() []NewsChannel
	CacheNewsChannel(NewsChannel)

	GetStoreChannelById(Snowflake) StoreChannel
	GetStoreChannelsByName(string, bool) []StoreChannel
	GetStoreChannelsCache() map[Snowflake]StoreChannel
	GetStoreChannels() []StoreChannel
	CacheStoreChannel(StoreChannel)

	GetVoiceChannelById(Snowflake) VoiceChannel
	GetVoiceChannelsByName(string, bool) []CategoryChannel
	GetVoiceChannelsCache() map[Snowflake]VoiceChannel
	GetVoiceChannels() []VoiceChannel
	CacheVoiceChannel(VoiceChannel)

	GetCategoryById(Snowflake) CategoryChannel
	GetCategoriesByName(string, bool) []CategoryChannel
	GetCategoriesCache() map[Snowflake]CategoryChannel
	GetCategories() []CategoryChannel
	CacheEmoteCategory(CategoryChannel)

	GetEmoteById(Snowflake) Emote
	GetEmotesByName(string, bool) []Emote
	GetEmotesCache() map[Snowflake]Emote
	GetEmotes() []Emote
	CacheEmote(Emote)
}
