package disgo

import (
	"strings"

	"github.com/DiscoOrg/disgo/models"
)

type Cache interface {
	GetGuildById(models.Snowflake) models.Guild
	GetGuildsByName(string, bool) []models.Guild
	GetGuildsCache() map[models.Snowflake]models.Guild
	GetGuilds() []models.Guild

	GetUserById(models.Snowflake) models.User
	GetUsersByName(models.Snowflake, bool) []models.User
	GetUsersCache() map[models.Snowflake]models.User
	GetUsers() []models.User

	GetMemberById(models.Snowflake) models.Member
	GetMemberByName(string, bool) []models.Member
	GetMembersCache() map[models.Snowflake]models.Member
	GetMembers() []models.Member

	GetChannelById(models.Snowflake) Channel
	GetChannelsByName(string, bool) []Channel
	GetChannelsCache() map[models.Snowflake]Channel
	GetChannels() []Channel

	GetPrivateChannelById(models.Snowflake) Channel
	GetPrivateChannelsByName(string, bool) []Channel // not sure if we need this lul
	GetPrivateChannelsCache() map[models.Snowflake]Channel
	GetPrivateChannels() []Channel

	GetGuildChannelById(models.Snowflake) GuildChannel
	GetGuildChannelsByName(string, bool) []GuildChannel
	GetGuildChannelsCache() map[models.Snowflake]GuildChannel
	GetGuildChannels() []GuildChannel

	GetTextChannelById(models.Snowflake) TextChannel
	GetTextChannelsByName(string, bool) []CategoryChannel
	GetTextChannelsCache() map[models.Snowflake]TextChannel
	GetTextChannels() []TextChannel

	GetNewsChannelById(models.Snowflake) NewsChannel
	GetNewsChannelsByName(string, bool) []NewsChannel
	GetNewsChannelsCache() map[models.Snowflake]NewsChannel
	GetNewsChannels() []NewsChannel

	GetStoreChannelById(models.Snowflake) StoreChannel
	GetStoreChannelsByName(string, bool) []StoreChannel
	GetStoreChannelsCache() map[models.Snowflake]StoreChannel
	GetStoreChannels() []StoreChannel

	GetVoiceChannelById(models.Snowflake) VoiceChannel
	GetVoiceChannelsByName(string, bool) []CategoryChannel
	GetVoiceChannelsCache() map[models.Snowflake]VoiceChannel
	GetVoiceChannels() []VoiceChannel

	GetCategoryById(models.Snowflake) CategoryChannel
	GetCategoriesByName(string, bool) []CategoryChannel
	GetCategoriesCache() map[models.Snowflake]CategoryChannel
	GetCategories() []CategoryChannel

	GetEmoteById(models.Snowflake) models.Emote
	GetEmotesByName(string, bool) []models.Emote
	GetEmotesCache() map[models.Snowflake]models.Emote
	GetEmotes() []models.Emote
	cacheEmote(models.Emote)
}

type CacheImpl struct {
	guilds  map[models.Snowflake]models.Guild
	members map[models.Snowflake]models.Member
	users   map[models.Snowflake]models.User
	channel map[models.Snowflake]Channel
	emotes  map[models.Snowflake]models.Emote
}

func (c CacheImpl) GetGuildById(id models.Snowflake) models.Guild {
	return c.guilds[id]
}
func (c CacheImpl) GetGuildsByName(name string, ignoreCase bool) []models.Guild {
	if ignoreCase {
		name = strings.ToLower(name)
	}
	guilds := make([]models.Guild, 1)
	for _, guild := range c.guilds {
		if ignoreCase && strings.ToLower(guild.Name) == name || !ignoreCase && guild.Name == name{
			guilds = append(guilds, guild)
		}
	}
	return guilds
}
func (c CacheImpl) GetGuildsCache() map[models.Snowflake]models.Guild {
	return c.guilds
}
func (c CacheImpl) GetGuilds() []models.Guild {
	guilds := make([]models.Guild, len(c.guilds))
	i := 0
	for _, guild := range c.guilds {
		guilds[i] = guild
		i++
	}
	return guilds
}



func (c CacheImpl) GetUserById(id models.Snowflake) models.User {
	return c.users[id]
}
