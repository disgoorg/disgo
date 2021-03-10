package api

type Cache interface {
	GetGuildById(Snowflake) Guild
	GetGuildsByName(string, bool) []Guild
	GetGuildsCache() map[Snowflake]Guild
	GetGuilds() []Guild

	GetUserById(Snowflake) User
	GetUsersByName(Snowflake, bool) []User
	GetUsersCache() map[Snowflake]User
	GetUsers() []User

	GetMemberById(Snowflake) Member
	GetMemberByName(string, bool) []Member
	GetMembersCache() map[Snowflake]Member
	GetMembers() []Member

	GetChannelById(Snowflake) Channel
	GetChannelsByName(string, bool) []Channel
	GetChannelsCache() map[Snowflake]Channel
	GetChannels() []Channel

	GetPrivateChannelById(Snowflake) Channel
	GetPrivateChannelsByName(string, bool) []Channel // not sure if we need this lul
	GetPrivateChannelsCache() map[Snowflake]Channel
	GetPrivateChannels() []Channel

	GetGuildChannelById(Snowflake) GuildChannel
	GetGuildChannelsByName(string, bool) []GuildChannel
	GetGuildChannelsCache() map[Snowflake]GuildChannel
	GetGuildChannels() []GuildChannel

	GetTextChannelById(Snowflake) TextChannel
	GetTextChannelsByName(string, bool) []CategoryChannel
	GetTextChannelsCache() map[Snowflake]TextChannel
	GetTextChannels() []TextChannel

	GetNewsChannelById(Snowflake) NewsChannel
	GetNewsChannelsByName(string, bool) []NewsChannel
	GetNewsChannelsCache() map[Snowflake]NewsChannel
	GetNewsChannels() []NewsChannel

	GetStoreChannelById(Snowflake) StoreChannel
	GetStoreChannelsByName(string, bool) []StoreChannel
	GetStoreChannelsCache() map[Snowflake]StoreChannel
	GetStoreChannels() []StoreChannel

	GetVoiceChannelById(Snowflake) VoiceChannel
	GetVoiceChannelsByName(string, bool) []CategoryChannel
	GetVoiceChannelsCache() map[Snowflake]VoiceChannel
	GetVoiceChannels() []VoiceChannel

	GetCategoryById(Snowflake) CategoryChannel
	GetCategoriesByName(string, bool) []CategoryChannel
	GetCategoriesCache() map[Snowflake]CategoryChannel
	GetCategories() []CategoryChannel

	GetEmoteById(Snowflake) Emote
	GetEmotesByName(string, bool) []Emote
	GetEmotesCache() map[Snowflake]Emote
	GetEmotes() []Emote
	cacheEmote(Emote)
}
