package internal

import (
	"strings"

	"github.com/DiscoOrg/disgo/api"
)

type CacheImpl struct {
	guilds  map[api.Snowflake]api.Guild
	members map[api.Snowflake]api.Member
	users   map[api.Snowflake]api.User
	channel map[api.Snowflake]api.Channel
	emotes  map[api.Snowflake]api.Emote
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
		if ignoreCase && strings.ToLower(guild.Name) == name || !ignoreCase && guild.Name == name{
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
