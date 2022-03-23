package cache

import (
	"sync"

	"github.com/DisgoOrg/snowflake"
	"github.com/disgoorg/disgo/discord"
)

func NewGuildCache(flags Flags, neededFlags Flags, policy Policy[discord.Guild]) GuildCache {
	return &GuildCacheImpl{
		Cache:             NewCache[discord.Guild](flags, neededFlags, policy),
		unreadyGuilds:     map[int]map[snowflake.Snowflake]struct{}{},
		unavailableGuilds: map[snowflake.Snowflake]struct{}{},
	}
}

type GuildCache interface {
	Cache[discord.Guild]

	SetReady(shardID int, guildID snowflake.Snowflake)
	SetUnready(shardID int, guildID snowflake.Snowflake)
	IsUnready(shardID int, guildID snowflake.Snowflake) bool
	UnreadyGuilds(shardID int) []snowflake.Snowflake

	SetUnavailable(guildID snowflake.Snowflake)
	SetAvailable(guildID snowflake.Snowflake)
	IsUnavailable(guildID snowflake.Snowflake) bool
	UnavailableGuilds() []snowflake.Snowflake
}

type GuildCacheImpl struct {
	Cache[discord.Guild]
	unreadyGuildsMu sync.RWMutex
	unreadyGuilds   map[int]map[snowflake.Snowflake]struct{}

	unavailableGuildsMu sync.RWMutex
	unavailableGuilds   map[snowflake.Snowflake]struct{}
}

func (c *GuildCacheImpl) SetReady(shardID int, guildID snowflake.Snowflake) {
	c.unreadyGuildsMu.Lock()
	defer c.unreadyGuildsMu.Unlock()
	if _, ok := c.unreadyGuilds[shardID]; !ok {
		return
	}
	delete(c.unreadyGuilds[shardID], guildID)
}

func (c *GuildCacheImpl) SetUnready(shardID int, guildID snowflake.Snowflake) {
	c.unreadyGuildsMu.Lock()
	defer c.unreadyGuildsMu.Unlock()
	if _, ok := c.unreadyGuilds[shardID]; !ok {
		c.unreadyGuilds[shardID] = map[snowflake.Snowflake]struct{}{}
	}
	c.unreadyGuilds[shardID][guildID] = struct{}{}
}

func (c *GuildCacheImpl) IsUnready(shardID int, guildID snowflake.Snowflake) bool {
	c.unreadyGuildsMu.RLock()
	defer c.unreadyGuildsMu.RUnlock()
	if _, ok := c.unreadyGuilds[shardID]; !ok {
		return false
	}
	_, ok := c.unreadyGuilds[shardID][guildID]
	return ok
}

func (c *GuildCacheImpl) UnreadyGuilds(shardID int) []snowflake.Snowflake {
	c.unreadyGuildsMu.RLock()
	defer c.unreadyGuildsMu.RUnlock()
	if _, ok := c.unreadyGuilds[shardID]; !ok {
		return nil
	}
	guilds := make([]snowflake.Snowflake, len(c.unreadyGuilds[shardID]))
	var i int
	for guildID := range c.unreadyGuilds[shardID] {
		guilds[i] = guildID
		i++
	}
	return guilds
}

func (c *GuildCacheImpl) SetUnavailable(id snowflake.Snowflake) {
	c.unavailableGuildsMu.Lock()
	defer c.unavailableGuildsMu.Unlock()
	c.Remove(id)
	c.unavailableGuilds[id] = struct{}{}
}

func (c *GuildCacheImpl) SetAvailable(guildID snowflake.Snowflake) {
	c.unavailableGuildsMu.Lock()
	defer c.unavailableGuildsMu.Unlock()
	delete(c.unavailableGuilds, guildID)
}

func (c *GuildCacheImpl) IsUnavailable(guildID snowflake.Snowflake) bool {
	c.unavailableGuildsMu.RLock()
	defer c.unavailableGuildsMu.RUnlock()
	_, ok := c.unavailableGuilds[guildID]
	return ok
}

func (c *GuildCacheImpl) UnavailableGuilds() []snowflake.Snowflake {
	c.unavailableGuildsMu.RLock()
	defer c.unavailableGuildsMu.RUnlock()
	guilds := make([]snowflake.Snowflake, len(c.unavailableGuilds))
	var i int
	for guildID := range c.unavailableGuilds {
		guilds[i] = guildID
		i++
	}
	return guilds
}
