package cache

import (
	"sync"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

func NewGuildCache(flags Flags, neededFlags Flags, policy Policy[discord.Guild]) GuildCache {
	return &GuildCacheImpl{
		Cache:             NewCache[discord.Guild](flags, neededFlags, policy),
		unreadyGuilds:     map[int]map[snowflake.ID]struct{}{},
		unavailableGuilds: map[snowflake.ID]struct{}{},
	}
}

type GuildCache interface {
	Cache[discord.Guild]

	SetReady(shardID int, guildID snowflake.ID)
	SetUnready(shardID int, guildID snowflake.ID)
	IsUnready(shardID int, guildID snowflake.ID) bool
	UnreadyGuilds(shardID int) []snowflake.ID

	SetUnavailable(guildID snowflake.ID)
	SetAvailable(guildID snowflake.ID)
	IsUnavailable(guildID snowflake.ID) bool
	UnavailableGuilds() []snowflake.ID
}

type GuildCacheImpl struct {
	Cache[discord.Guild]
	unreadyGuildsMu sync.RWMutex
	unreadyGuilds   map[int]map[snowflake.ID]struct{}

	unavailableGuildsMu sync.RWMutex
	unavailableGuilds   map[snowflake.ID]struct{}
}

func (c *GuildCacheImpl) SetReady(shardID int, guildID snowflake.ID) {
	c.unreadyGuildsMu.Lock()
	defer c.unreadyGuildsMu.Unlock()
	if _, ok := c.unreadyGuilds[shardID]; !ok {
		return
	}
	delete(c.unreadyGuilds[shardID], guildID)
}

func (c *GuildCacheImpl) SetUnready(shardID int, guildID snowflake.ID) {
	c.unreadyGuildsMu.Lock()
	defer c.unreadyGuildsMu.Unlock()
	if _, ok := c.unreadyGuilds[shardID]; !ok {
		c.unreadyGuilds[shardID] = map[snowflake.ID]struct{}{}
	}
	c.unreadyGuilds[shardID][guildID] = struct{}{}
}

func (c *GuildCacheImpl) IsUnready(shardID int, guildID snowflake.ID) bool {
	c.unreadyGuildsMu.RLock()
	defer c.unreadyGuildsMu.RUnlock()
	if _, ok := c.unreadyGuilds[shardID]; !ok {
		return false
	}
	_, ok := c.unreadyGuilds[shardID][guildID]
	return ok
}

func (c *GuildCacheImpl) UnreadyGuilds(shardID int) []snowflake.ID {
	c.unreadyGuildsMu.RLock()
	defer c.unreadyGuildsMu.RUnlock()
	if _, ok := c.unreadyGuilds[shardID]; !ok {
		return nil
	}
	guilds := make([]snowflake.ID, len(c.unreadyGuilds[shardID]))
	var i int
	for guildID := range c.unreadyGuilds[shardID] {
		guilds[i] = guildID
		i++
	}
	return guilds
}

func (c *GuildCacheImpl) SetUnavailable(id snowflake.ID) {
	c.unavailableGuildsMu.Lock()
	defer c.unavailableGuildsMu.Unlock()
	c.Remove(id)
	c.unavailableGuilds[id] = struct{}{}
}

func (c *GuildCacheImpl) SetAvailable(guildID snowflake.ID) {
	c.unavailableGuildsMu.Lock()
	defer c.unavailableGuildsMu.Unlock()
	delete(c.unavailableGuilds, guildID)
}

func (c *GuildCacheImpl) IsUnavailable(guildID snowflake.ID) bool {
	c.unavailableGuildsMu.RLock()
	defer c.unavailableGuildsMu.RUnlock()
	_, ok := c.unavailableGuilds[guildID]
	return ok
}

func (c *GuildCacheImpl) UnavailableGuilds() []snowflake.ID {
	c.unavailableGuildsMu.RLock()
	defer c.unavailableGuildsMu.RUnlock()
	guilds := make([]snowflake.ID, len(c.unavailableGuilds))
	var i int
	for guildID := range c.unavailableGuilds {
		guilds[i] = guildID
		i++
	}
	return guilds
}
