package internal

import (
	"strings"

	"github.com/DiscoOrg/disgo/api/models"
)

type CacheImpl struct {
	guilds  map[models.Snowflake]models.Guild
	members map[models.Snowflake]models.Member
	users   map[models.Snowflake]models.User
	channel map[models.Snowflake]models.Channel
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
