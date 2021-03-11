package internal

import (
	"strings"

	"github.com/DiscoOrg/disgo/api"
)

func newCacheImpl(memberCachePolicy api.MemberCachePolicy) api.Cache {
	return &CacheImpl{
		memberCachePolicy: memberCachePolicy,
		guilds:            map[api.Snowflake]api.Guild{},
		members:           map[api.Snowflake]api.Member{},
		users:             map[api.Snowflake]api.User{},
		channel:           map[api.Snowflake]api.Channel{},
		emotes:            map[api.Snowflake]api.Emote{},
	}
}

type CacheImpl struct {
	memberCachePolicy api.MemberCachePolicy
	guilds            map[api.Snowflake]api.Guild
	members           map[api.Snowflake]api.Member
	users             map[api.Snowflake]api.User
	channel           map[api.Snowflake]api.Channel
	emotes            map[api.Snowflake]api.Emote
}

func (c CacheImpl) GetGuildById(id api.Snowflake) api.Guild {
	return c.guilds[id]
}
func (c CacheImpl) GetGuildsByName(name string, ignoreCase bool) []api.Guild {
	if ignoreCase {
		name = strings.ToLower(name)
	}
	guilds := make([]api.Guild, 1)
	for _, guild := range c.guilds {
		if ignoreCase && strings.ToLower(guild.Name) == name || !ignoreCase && guild.Name == name {
			guilds = append(guilds, guild)
		}
	}
	return guilds
}
func (c CacheImpl) GetGuildsCache() map[api.Snowflake]api.Guild {
	return c.guilds
}
func (c CacheImpl) GetGuilds() []api.Guild {
	guilds := make([]api.Guild, len(c.guilds))
	i := 0
	for _, guild := range c.guilds {
		guilds[i] = guild
		i++
	}
	return guilds
}
func (c CacheImpl) CacheGuild(guild api.Guild) {
	c.guilds[guild.ID] = guild
}

func (c CacheImpl) GetUserById(id api.Snowflake) api.User {
	return c.users[id]
}
func (c CacheImpl) GetUsersByName(id api.Snowflake, b bool) []api.User {

}
func (c CacheImpl) GetUsersCache() map[api.Snowflake]api.User {

}
func (c CacheImpl) GetUsers() []api.User {
	users := make([]api.User, len(c.users))
	i := 0
	for _, user := range c.users {
		users[i] = user
		i++
	}
	return users
}
func (c CacheImpl) CacheUser(user api.User) {

}

func (c CacheImpl) GetMemberById(snowflake api.Snowflake) api.Member {

}
func (c CacheImpl) GetMemberByName(s string, b bool) []api.Member {

}
func (c CacheImpl) GetMembersCache() map[api.Snowflake]api.Member {

}
func (c CacheImpl) GetMembers() []api.Member {

}
func (c CacheImpl) CacheMember(member api.Member) {

}

func (c CacheImpl) GetChannelById(snowflake api.Snowflake) api.Channel {

}
func (c CacheImpl) GetChannelsByName(s string, b bool) []api.Channel {

}
func (c CacheImpl) GetChannelsCache() map[api.Snowflake]api.Channel {

}
func (c CacheImpl) GetChannels() []api.Channel {

}
func (c CacheImpl) CacheChannel(channel api.Channel) {

}

func (c CacheImpl) GetPrivateChannelById(snowflake api.Snowflake) api.Channel {

}
func (c CacheImpl) GetPrivateChannelsByName(s string, b bool) []api.Channel {

}
func (c CacheImpl) GetPrivateChannelsCache() map[api.Snowflake]api.Channel {

}
func (c CacheImpl) GetPrivateChannels() []api.Channel {

}
func (c CacheImpl) CachePrivateChannel(channel api.Channel) {

}

func (c CacheImpl) GetGuildChannelById(snowflake api.Snowflake) api.GuildChannel {

}
func (c CacheImpl) GetGuildChannelsByName(s string, b bool) []api.GuildChannel {

}
func (c CacheImpl) GetGuildChannelsCache() map[api.Snowflake]api.GuildChannel {

}
func (c CacheImpl) GetGuildChannels() []api.GuildChannel {

}
func (c CacheImpl) CacheGuildChannel(channel api.GuildChannel) {

}

func (c CacheImpl) GetTextChannelById(snowflake api.Snowflake) api.TextChannel {

}
func (c CacheImpl) GetTextChannelsByName(s string, b bool) []api.CategoryChannel {

}
func (c CacheImpl) GetTextChannelsCache() map[api.Snowflake]api.TextChannel {

}
func (c CacheImpl) GetTextChannels() []api.TextChannel {

}
func (c CacheImpl) CacheTextChannel(channel api.TextChannel) {

}

func (c CacheImpl) GetNewsChannelById(snowflake api.Snowflake) api.NewsChannel {

}
func (c CacheImpl) GetNewsChannelsByName(s string, b bool) []api.NewsChannel {

}
func (c CacheImpl) GetNewsChannelsCache() map[api.Snowflake]api.NewsChannel {

}
func (c CacheImpl) GetNewsChannels() []api.NewsChannel {

}
func (c CacheImpl) CacheNewsChannel(channel api.NewsChannel) {

}

func (c CacheImpl) GetStoreChannelById(snowflake api.Snowflake) api.StoreChannel {

}
func (c CacheImpl) GetStoreChannelsByName(s string, b bool) []api.StoreChannel {

}
func (c CacheImpl) GetStoreChannelsCache() map[api.Snowflake]api.StoreChannel {

}
func (c CacheImpl) GetStoreChannels() []api.StoreChannel {

}
func (c CacheImpl) CacheStoreChannel(channel api.StoreChannel) {

}

func (c CacheImpl) GetVoiceChannelById(snowflake api.Snowflake) api.VoiceChannel {

}
func (c CacheImpl) GetVoiceChannelsByName(s string, b bool) []api.CategoryChannel {

}
func (c CacheImpl) GetVoiceChannelsCache() map[api.Snowflake]api.VoiceChannel {

}
func (c CacheImpl) GetVoiceChannels() []api.VoiceChannel {

}
func (c CacheImpl) CacheVoiceChannel(channel api.VoiceChannel) {

}

func (c CacheImpl) GetCategoryById(snowflake api.Snowflake) api.CategoryChannel {

}
func (c CacheImpl) GetCategoriesByName(s string, b bool) []api.CategoryChannel {

}
func (c CacheImpl) GetCategoriesCache() map[api.Snowflake]api.CategoryChannel {

}
func (c CacheImpl) GetCategories() []api.CategoryChannel {

}

func (c CacheImpl) CacheEmoteCategory(channel api.CategoryChannel) {

}
func (c CacheImpl) GetEmoteById(snowflake api.Snowflake) api.Emote {

}
func (c CacheImpl) GetEmotesByName(s string, b bool) []api.Emote {

}
func (c CacheImpl) GetEmotesCache() map[api.Snowflake]api.Emote {

}
func (c CacheImpl) GetEmotes() []api.Emote {

}
func (c CacheImpl) CacheEmote(emote api.Emote) {

}
