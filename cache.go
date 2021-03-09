package disgo

import (
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

	GetChannelById(models.Snowflake) models.Channel
	GetChannelsByName(string, bool) []models.Channel
	GetChannelsCache() map[models.Snowflake]models.Channel
	GetChannels() []models.Channel

	GetPrivateChannelById(models.Snowflake) models.Channel
	GetPrivateChannelsByName(string, bool) []models.Channel // not sure if we need this lul
	GetPrivateChannelsCache() map[models.Snowflake]models.Channel
	GetPrivateChannels() []models.Channel

	GetGuildChannelById(models.Snowflake) models.GuildChannel
	GetGuildChannelsByName(string, bool) []models.GuildChannel
	GetGuildChannelsCache() map[models.Snowflake]models.GuildChannel
	GetGuildChannels() []models.GuildChannel

	GetTextChannelById(models.Snowflake) models.TextChannel
	GetTextChannelsByName(string, bool) []models.CategoryChannel
	GetTextChannelsCache() map[models.Snowflake]models.TextChannel
	GetTextChannels() []models.TextChannel

	GetNewsChannelById(models.Snowflake) models.NewsChannel
	GetNewsChannelsByName(string, bool) []models.NewsChannel
	GetNewsChannelsCache() map[models.Snowflake]models.NewsChannel
	GetNewsChannels() []models.NewsChannel

	GetStoreChannelById(models.Snowflake) models.StoreChannel
	GetStoreChannelsByName(string, bool) []models.StoreChannel
	GetStoreChannelsCache() map[models.Snowflake]models.StoreChannel
	GetStoreChannels() []models.StoreChannel

	GetVoiceChannelById(models.Snowflake) models.VoiceChannel
	GetVoiceChannelsByName(string, bool) []models.CategoryChannel
	GetVoiceChannelsCache() map[models.Snowflake]models.VoiceChannel
	GetVoiceChannels() []models.VoiceChannel

	GetCategoryById(models.Snowflake) models.CategoryChannel
	GetCategoriesByName(string, bool) []models.CategoryChannel
	GetCategoriesCache() map[models.Snowflake]models.CategoryChannel
	GetCategories() []models.CategoryChannel

	GetEmoteById(models.Snowflake) models.Emote
	GetEmotesByName(string, bool) []models.Emote
	GetEmotesCache() map[models.Snowflake]models.Emote
	GetEmotes() []models.Emote
	cacheEmote(models.Emote)
}
