package api

// Cache allows you to access the objects that are stored in-memory by Discord
type Cache interface {
	Close()
	DoCleanup()

	Guild(Snowflake) *Guild
	GuildsByName(string, bool) []*Guild
	Guilds() []*Guild
	GuildCache() map[Snowflake]*Guild
	CacheGuild(*Guild)
	UncacheGuild(Snowflake)

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

	VoiceState(guildID Snowflake, userID Snowflake) *VoiceState
	VoiceStates(guildID Snowflake) []*VoiceState
	VoiceStateCache(guildID Snowflake) map[Snowflake]*VoiceState
	CacheVoiceState(voiceState *VoiceState)
	UncacheVoiceState(guildID Snowflake, userID Snowflake)

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

	Category(Snowflake) *Category
	CategoriesByName(Snowflake, string, bool) []*Category
	Categories(Snowflake) []*Category
	AllCategories() []*Category
	CategoryCache(Snowflake) map[Snowflake]*Category
	AllCategoryCache() map[Snowflake]map[Snowflake]*Category
	CacheCategory(*Category)
	UncacheCategory(Snowflake, Snowflake)
	FindCategory(Snowflake, func(*Category) bool) *Category
	FindCategories(Snowflake, func(*Category) bool) []*Category

	/*Emote(Snowflake) *Emote
	EmotesByName(string, bool) []*Emote
	Emotes() []*Emote
	EmoteCache() map[Snowflake]*Emote
	CacheEmote(*Emote)
	UncacheEmote(Snowflake)*/
}
