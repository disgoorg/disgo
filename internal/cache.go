package internal

import (
	"strings"

	"github.com/DiscoOrg/disgo/api"
)

func newCacheImpl() api.Cache {
	return &CacheImpl{
		guilds:  map[api.Snowflake]api.Guild{},
		members: map[api.Snowflake]api.Member{},
		users:   map[api.Snowflake]api.User{},
		channel: map[api.Snowflake]api.Channel{},
		emotes:  map[api.Snowflake]api.Emote{},
	}
}

type CacheImpl struct {
	guilds  map[api.Snowflake]api.Guild
	members map[api.Snowflake]api.Member
	users   map[api.Snowflake]api.User
	channel map[api.Snowflake]api.Channel
	emotes  map[api.Snowflake]api.Emote
}

func (c CacheImpl) GetUsersByName(snowflake api.Snowflake, b bool) []api.User {
	panic("implement me")
}

func (c CacheImpl) GetUsersCache() map[api.Snowflake]api.User {
	panic("implement me")
}

func (c CacheImpl) GetUsers() []api.User {
	panic("implement me")
}

func (c CacheImpl) GetMemberById(snowflake api.Snowflake) api.Member {
	panic("implement me")
}

func (c CacheImpl) GetMemberByName(s string, b bool) []api.Member {
	panic("implement me")
}

func (c CacheImpl) GetMembersCache() map[api.Snowflake]api.Member {
	panic("implement me")
}

func (c CacheImpl) GetMembers() []api.Member {
	panic("implement me")
}

func (c CacheImpl) GetChannelById(snowflake api.Snowflake) api.Channel {
	panic("implement me")
}

func (c CacheImpl) GetChannelsByName(s string, b bool) []api.Channel {
	panic("implement me")
}

func (c CacheImpl) GetChannelsCache() map[api.Snowflake]api.Channel {
	panic("implement me")
}

func (c CacheImpl) GetChannels() []api.Channel {
	panic("implement me")
}

func (c CacheImpl) GetPrivateChannelById(snowflake api.Snowflake) api.Channel {
	panic("implement me")
}

func (c CacheImpl) GetPrivateChannelsByName(s string, b bool) []api.Channel {
	panic("implement me")
}

func (c CacheImpl) GetPrivateChannelsCache() map[api.Snowflake]api.Channel {
	panic("implement me")
}

func (c CacheImpl) GetPrivateChannels() []api.Channel {
	panic("implement me")
}

func (c CacheImpl) GetGuildChannelById(snowflake api.Snowflake) api.GuildChannel {
	panic("implement me")
}

func (c CacheImpl) GetGuildChannelsByName(s string, b bool) []api.GuildChannel {
	panic("implement me")
}

func (c CacheImpl) GetGuildChannelsCache() map[api.Snowflake]api.GuildChannel {
	panic("implement me")
}

func (c CacheImpl) GetGuildChannels() []api.GuildChannel {
	panic("implement me")
}

func (c CacheImpl) GetTextChannelById(snowflake api.Snowflake) api.TextChannel {
	panic("implement me")
}

func (c CacheImpl) GetTextChannelsByName(s string, b bool) []api.CategoryChannel {
	panic("implement me")
}

func (c CacheImpl) GetTextChannelsCache() map[api.Snowflake]api.TextChannel {
	panic("implement me")
}

func (c CacheImpl) GetTextChannels() []api.TextChannel {
	panic("implement me")
}

func (c CacheImpl) GetNewsChannelById(snowflake api.Snowflake) api.NewsChannel {
	panic("implement me")
}

func (c CacheImpl) GetNewsChannelsByName(s string, b bool) []api.NewsChannel {
	panic("implement me")
}

func (c CacheImpl) GetNewsChannelsCache() map[api.Snowflake]api.NewsChannel {
	panic("implement me")
}

func (c CacheImpl) GetNewsChannels() []api.NewsChannel {
	panic("implement me")
}

func (c CacheImpl) GetStoreChannelById(snowflake api.Snowflake) api.StoreChannel {
	panic("implement me")
}

func (c CacheImpl) GetStoreChannelsByName(s string, b bool) []api.StoreChannel {
	panic("implement me")
}

func (c CacheImpl) GetStoreChannelsCache() map[api.Snowflake]api.StoreChannel {
	panic("implement me")
}

func (c CacheImpl) GetStoreChannels() []api.StoreChannel {
	panic("implement me")
}

func (c CacheImpl) GetVoiceChannelById(snowflake api.Snowflake) api.VoiceChannel {
	panic("implement me")
}

func (c CacheImpl) GetVoiceChannelsByName(s string, b bool) []api.CategoryChannel {
	panic("implement me")
}

func (c CacheImpl) GetVoiceChannelsCache() map[api.Snowflake]api.VoiceChannel {
	panic("implement me")
}

func (c CacheImpl) GetVoiceChannels() []api.VoiceChannel {
	panic("implement me")
}

func (c CacheImpl) GetCategoryById(snowflake api.Snowflake) api.CategoryChannel {
	panic("implement me")
}

func (c CacheImpl) GetCategoriesByName(s string, b bool) []api.CategoryChannel {
	panic("implement me")
}

func (c CacheImpl) GetCategoriesCache() map[api.Snowflake]api.CategoryChannel {
	panic("implement me")
}

func (c CacheImpl) GetCategories() []api.CategoryChannel {
	panic("implement me")
}

func (c CacheImpl) GetEmoteById(snowflake api.Snowflake) api.Emote {
	panic("implement me")
}

func (c CacheImpl) GetEmotesByName(s string, b bool) []api.Emote {
	panic("implement me")
}

func (c CacheImpl) GetEmotesCache() map[api.Snowflake]api.Emote {
	panic("implement me")
}

func (c CacheImpl) GetEmotes() []api.Emote {
	panic("implement me")
}

func (c CacheImpl) CacheEmote(emote api.Emote) {
	panic("implement me")
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

func (c CacheImpl) GetUserById(id api.Snowflake) api.User {
	return c.users[id]
}
