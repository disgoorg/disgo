package api

// Cache allows you to access the objects that are stored in-memory by Discord
type Cache interface {
	Disgo() Disgo
	Close()
	DoCleanup()
	CacheFlags() CacheFlags

	Command(commandID Snowflake) *Command
	GuildCommandCache(guildID Snowflake) map[Snowflake]*Command
	AllGuildCommandCache() map[Snowflake]map[Snowflake]*Command
	GlobalCommandCache() map[Snowflake]*Command
	CacheGlobalCommand(command *Command) *Command
	CacheGuildCommand(command *Command) *Command
	UncacheCommand(commandID Snowflake)

	User(Snowflake) *User
	UserByTag(string) *User
	UsersByName(string, bool) []*User
	Users() []*User
	UserCache() map[Snowflake]*User
	CacheUser(*User) *User
	UncacheUser(Snowflake)
	FindUser(func(*User) bool) *User
	FindUsers(func(*User) bool) []*User

	Guild(guildId Snowflake) *Guild
	GuildsByName(name string, ignoreCase bool) []*Guild
	Guilds() []*Guild
	GuildCache() map[Snowflake]*Guild
	CacheGuild(guild *Guild) *Guild
	UncacheGuild(guildID Snowflake)

	Message(channelID Snowflake, messageID Snowflake) *Message
	Messages(channelID Snowflake) []*Message
	MessageCache(channelID Snowflake) map[Snowflake]*Message
	AllMessageCache() map[Snowflake]map[Snowflake]*Message
	CacheMessage(message *Message) *Message
	UncacheMessage(channelID Snowflake, messageID Snowflake)

	Member(Snowflake, Snowflake) *Member
	MemberByTag(Snowflake, string) *Member
	MembersByName(Snowflake, string, bool) []*Member
	Members(Snowflake) []*Member
	AllMembers() []*Member
	MemberCache(Snowflake) map[Snowflake]*Member
	AllMemberCache() map[Snowflake]map[Snowflake]*Member
	CacheMember(member *Member) *Member
	UncacheMember(Snowflake, Snowflake)
	FindMember(Snowflake, func(*Member) bool) *Member
	FindMembers(Snowflake, func(*Member) bool) []*Member

	ThreadMember(Snowflake, Snowflake, Snowflake) *ThreadMember
	ThreadMembers(Snowflake, Snowflake) []*ThreadMember
	ThreadMemberCache(snowflake Snowflake) map[Snowflake]map[Snowflake]*ThreadMember
	AllThreadMemberCache() map[Snowflake]map[Snowflake]map[Snowflake]*ThreadMember
	CacheThreadMember(member *ThreadMember) *ThreadMember
	UncacheThreadMember(Snowflake, Snowflake, Snowflake)
	UncacheThreadMembers(guildID Snowflake)

	VoiceState(guildID Snowflake, userID Snowflake) *VoiceState
	VoiceStates(guildID Snowflake) []*VoiceState
	VoiceStateCache(guildID Snowflake) map[Snowflake]*VoiceState
	CacheVoiceState(voiceState *VoiceState) *VoiceState
	UncacheVoiceState(guildID Snowflake, userID Snowflake)

	Role(roleID Snowflake) *Role
	RolesByName(Snowflake, string, bool) []*Role
	Roles(Snowflake) []*Role
	AllRoles() []*Role
	RoleCache(Snowflake) map[Snowflake]*Role
	AllRoleCache() map[Snowflake]map[Snowflake]*Role
	CacheRole(*Role) *Role
	UncacheRole(Snowflake, Snowflake)
	FindRole(Snowflake, func(*Role) bool) *Role
	FindRoles(Snowflake, func(*Role) bool) []*Role

	DMChannel(Snowflake) DMChannel
	DMChannels() []DMChannel
	DMChannelCache() map[Snowflake]DMChannel
	CacheDMChannel(DMChannel) DMChannel
	UncacheDMChannel(dmChannelID Snowflake)
	FindDMChannel(func(DMChannel) bool) DMChannel
	FindDMChannels(func(DMChannel) bool) []DMChannel

	Channel(Snowflake) Channel
	MessageChannel(Snowflake) MessageChannel
	GuildChannel(Snowflake) GuildChannel

	TextChannel(Snowflake) TextChannel
	TextChannelsByName(Snowflake, string, bool) []TextChannel
	TextChannels(Snowflake) []TextChannel
	TextChannelCache(Snowflake) map[Snowflake]TextChannel
	CacheTextChannel(TextChannel) TextChannel
	UncacheTextChannel(Snowflake, Snowflake)
	FindTextChannel(Snowflake, func(TextChannel) bool) TextChannel
	FindTextChannels(Snowflake, func(TextChannel) bool) []TextChannel

	Thread(Snowflake) Thread
	ThreadsByName(Snowflake, string, bool) []Thread
	Threads(Snowflake) []Thread
	ThreadCache(Snowflake) map[Snowflake]Thread
	AllThreadCache() map[Snowflake]map[Snowflake]Thread
	CacheThread(Thread) Thread
	UncacheThread(Snowflake, Snowflake)
	UncacheThreads(guildID Snowflake)
	FindThread(Snowflake, func(Thread) bool) Thread
	FindThreads(Snowflake, func(Thread) bool) []Thread

	StoreChannel(Snowflake) StoreChannel
	StoreChannelsByName(Snowflake, string, bool) []StoreChannel
	StoreChannels(Snowflake) []StoreChannel
	StoreChannelCache(Snowflake) map[Snowflake]StoreChannel
	CacheStoreChannel(StoreChannel) StoreChannel
	UncacheStoreChannel(Snowflake, Snowflake)
	FindStoreChannel(Snowflake, func(StoreChannel) bool) StoreChannel
	FindStoreChannels(Snowflake, func(StoreChannel) bool) []StoreChannel

	VoiceChannel(Snowflake) VoiceChannel
	VoiceChannelsByName(Snowflake, string, bool) []VoiceChannel
	VoiceChannels(Snowflake) []VoiceChannel
	VoiceChannelCache(Snowflake) map[Snowflake]VoiceChannel
	CacheVoiceChannel(VoiceChannel) VoiceChannel
	UncacheVoiceChannel(Snowflake, Snowflake)
	FindVoiceChannel(Snowflake, func(VoiceChannel) bool) VoiceChannel
	FindVoiceChannels(Snowflake, func(VoiceChannel) bool) []VoiceChannel

	Category(Snowflake) Category
	CategoriesByName(Snowflake, string, bool) []Category
	Categories(Snowflake) []Category
	AllCategories() []Category
	CategoryCache(Snowflake) map[Snowflake]Category
	AllCategoryCache() map[Snowflake]map[Snowflake]Category
	CacheCategory(Category) Category
	UncacheCategory(Snowflake, Snowflake)
	FindCategory(Snowflake, func(Category) bool) Category
	FindCategories(Snowflake, func(Category) bool) []Category

	Emote(emoteID Snowflake) *Emoji
	EmotesByName(guildID Snowflake, name string, ignoreCase bool) []*Emoji
	Emotes(guildID Snowflake) []*Emoji
	EmoteCache(guildID Snowflake) map[Snowflake]*Emoji
	AllEmoteCache() map[Snowflake]map[Snowflake]*Emoji
	CacheEmote(*Emoji) *Emoji
	UncacheEmote(guildID Snowflake, emoteID Snowflake)
}
