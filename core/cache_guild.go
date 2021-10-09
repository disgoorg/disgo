package core

import (
	"sync"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/internal/rwsync"
)

type (
	GuildFindFunc func(guild *Guild) bool

	GuildCache interface {
		rwsync.RWLocker

		Get(guildID discord.Snowflake) *Guild
		GetCopy(guildID discord.Snowflake) *Guild
		Set(guild *Guild) *Guild
		Remove(guildID discord.Snowflake)

		Cache() map[discord.Snowflake]*Guild
		All() []*Guild

		FindFirst(guildFindFunc GuildFindFunc) *Guild
		FindAll(guildFindFunc GuildFindFunc) []*Guild

		SetReady(shardID int, guildID discord.Snowflake)
		SetUnready(shardID int, guildID discord.Snowflake)
		IsUnready(shardID int, guildID discord.Snowflake) bool
		UnreadyGuilds(shardID int) []discord.Snowflake

		SetUnavailable(guildID discord.Snowflake)
		SetAvailable(guildID discord.Snowflake)
		IsUnavailable(guildID discord.Snowflake) bool
		UnavailableGuilds() []discord.Snowflake
	}

	guildCacheImpl struct {
		sync.RWMutex
		cacheFlags CacheFlags
		guilds     map[discord.Snowflake]*Guild

		unreadyGuildsMu sync.RWMutex
		unreadyGuilds   map[int]map[discord.Snowflake]struct{}

		unavailableGuildsMu sync.RWMutex
		unavailableGuilds   map[discord.Snowflake]struct{}
	}
)

func NewGuildCache(cacheFlags CacheFlags) GuildCache {
	return &guildCacheImpl{
		cacheFlags:        cacheFlags,
		guilds:            map[discord.Snowflake]*Guild{},
		unreadyGuilds:     map[int]map[discord.Snowflake]struct{}{},
		unavailableGuilds: map[discord.Snowflake]struct{}{},
	}
}

func (c *guildCacheImpl) Get(guildID discord.Snowflake) *Guild {
	c.RLock()
	defer c.RUnlock()
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
	c.Lock()
	defer c.Unlock()
	gui, ok := c.guilds[guild.ID]
	if ok {
		*gui = *guild
		return gui
	}
	c.guilds[guild.ID] = guild
	return guild
}

func (c *guildCacheImpl) Remove(id discord.Snowflake) {
	c.Lock()
	defer c.Unlock()
	delete(c.guilds, id)
}

func (c *guildCacheImpl) Cache() map[discord.Snowflake]*Guild {
	return c.guilds
}

func (c *guildCacheImpl) All() []*Guild {
	c.RLock()
	defer c.RUnlock()
	guilds := make([]*Guild, len(c.guilds))
	i := 0
	for _, guild := range c.guilds {
		guilds[i] = guild
		i++
	}
	return guilds
}

func (c *guildCacheImpl) FindFirst(guildFindFunc GuildFindFunc) *Guild {
	c.RLock()
	defer c.RUnlock()
	for _, gui := range c.guilds {
		if guildFindFunc(gui) {
			return gui
		}
	}
	return nil
}

func (c *guildCacheImpl) FindAll(guildFindFunc GuildFindFunc) []*Guild {
	c.RLock()
	defer c.RUnlock()
	var guilds []*Guild
	for _, gui := range c.guilds {
		if guildFindFunc(gui) {
			guilds = append(guilds, gui)
		}
	}
	return guilds
}

func (c *guildCacheImpl) SetReady(shardID int, guildID discord.Snowflake) {
	c.unreadyGuildsMu.Lock()
	defer c.unreadyGuildsMu.Unlock()
	if _, ok := c.unreadyGuilds[shardID]; !ok {
		return
	}
	delete(c.unreadyGuilds[shardID], guildID)
}

func (c *guildCacheImpl) SetUnready(shardID int, guildID discord.Snowflake) {
	c.unreadyGuildsMu.Lock()
	defer c.unreadyGuildsMu.Unlock()
	if _, ok := c.unreadyGuilds[shardID]; !ok {
		c.unreadyGuilds[shardID] = map[discord.Snowflake]struct{}{}
	}
	c.unreadyGuilds[shardID][guildID] = struct{}{}
}

func (c *guildCacheImpl) IsUnready(shardID int, guildID discord.Snowflake) bool {
	c.unreadyGuildsMu.RLock()
	defer c.unreadyGuildsMu.RUnlock()
	if _, ok := c.unreadyGuilds[shardID]; !ok {
		return false
	}
	_, ok := c.unreadyGuilds[shardID][guildID]
	return ok
}

func (c *guildCacheImpl) UnreadyGuilds(shardID int) []discord.Snowflake {
	c.unreadyGuildsMu.RLock()
	defer c.unreadyGuildsMu.RUnlock()
	if _, ok := c.unreadyGuilds[shardID]; !ok {
		return nil
	}
	guilds := make([]discord.Snowflake, len(c.unreadyGuilds[shardID]))
	var i int
	for guildID := range c.unreadyGuilds[shardID] {
		guilds[i] = guildID
		i++
	}
	return guilds
}

func (c *guildCacheImpl) SetUnavailable(id discord.Snowflake) {
	c.unavailableGuildsMu.Lock()
	defer c.unavailableGuildsMu.Unlock()
	c.Remove(id)
	c.unavailableGuilds[id] = struct{}{}
}

func (c *guildCacheImpl) SetAvailable(guildID discord.Snowflake) {
	c.unavailableGuildsMu.Lock()
	defer c.unavailableGuildsMu.Unlock()
	delete(c.unavailableGuilds, guildID)
}

func (c *guildCacheImpl) IsUnavailable(guildID discord.Snowflake) bool {
	c.unavailableGuildsMu.RLock()
	defer c.unavailableGuildsMu.RUnlock()
	_, ok := c.unavailableGuilds[guildID]
	return ok
}

func (c *guildCacheImpl) UnavailableGuilds() []discord.Snowflake {
	c.unavailableGuildsMu.RLock()
	defer c.unavailableGuildsMu.RUnlock()
	guilds := make([]discord.Snowflake, len(c.unavailableGuilds))
	var i int
	for guildID := range c.unavailableGuilds {
		guilds[i] = guildID
		i++
	}
	return guilds
}
