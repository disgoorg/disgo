package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

type (
	GuildFindFunc func(guild *Guild) bool

	GuildCache interface {
		Get(guildID discord.Snowflake) *Guild
		GetCopy(guildID discord.Snowflake) *Guild
		Set(guild *Guild) *Guild
		Remove(guildID discord.Snowflake)

		Cache() map[discord.Snowflake]*Guild
		All() []*Guild

		FindFirst(guildFindFunc GuildFindFunc) *Guild
		FindAll(guildFindFunc GuildFindFunc) []*Guild
	}

	guildCacheImpl struct {
		cacheFlags CacheFlags
		guilds     map[discord.Snowflake]*Guild
	}
)

func NewGuildCache(cacheFlags CacheFlags) GuildCache {
	return &guildCacheImpl{
		cacheFlags: cacheFlags,
		guilds:     map[discord.Snowflake]*Guild{},
	}
}

func (c *guildCacheImpl) Get(guildID discord.Snowflake) *Guild {
	return c.guilds[guildID]
}

func (c *guildCacheImpl) GetCopy(guildID discord.Snowflake) *Guild {
	if guild := c.Get(guildID); guild != nil {
		gu := *guild
		return &gu
	}
	return nil
}

func (c *guildCacheImpl) Set(guild *Guild) *Guild {
	if c.cacheFlags.Missing(CacheFlagGuilds) {
		return guild
	}
	gui, ok := c.guilds[guild.ID]
	if ok {
		*gui = *guild
		return gui
	}
	c.guilds[guild.ID] = guild
	return guild
}

func (c *guildCacheImpl) Remove(id discord.Snowflake) {
	delete(c.guilds, id)
}

func (c *guildCacheImpl) Cache() map[discord.Snowflake]*Guild {
	return c.guilds
}

func (c *guildCacheImpl) All() []*Guild {
	guilds := make([]*Guild, len(c.guilds))
	i := 0
	for _, guild := range c.guilds {
		guilds[i] = guild
		i++
	}
	return guilds
}

func (c *guildCacheImpl) FindFirst(guildFindFunc GuildFindFunc) *Guild {
	for _, gui := range c.guilds {
		if guildFindFunc(gui) {
			return gui
		}
	}
	return nil
}

func (c *guildCacheImpl) FindAll(guildFindFunc GuildFindFunc) []*Guild {
	var guilds []*Guild
	for _, gui := range c.guilds {
		if guildFindFunc(gui) {
			guilds = append(guilds, gui)
		}
	}
	return guilds
}
