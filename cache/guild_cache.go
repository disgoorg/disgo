package cache

import (
	"sync"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

// GuildCache is a Cache for guilds.
// It also keeps track of unready and unavailable guilds.
type GuildCache interface {
	Cache[discord.Guild]

	// SetReady sets the specified guildID as ready in the specific shard.
	SetReady(shardID int, guildID snowflake.ID)

	// SetUnready sets the specified guildID as not ready in the specific shard.
	SetUnready(shardID int, guildID snowflake.ID)

	// IsUnready returns a bool indicating if the specified guildID is ready or not in the specific shard.
	IsUnready(shardID int, guildID snowflake.ID) bool

	// UnreadyGuilds returns all guildIDs that are not ready in the specific shard.
	UnreadyGuilds(shardID int) []snowflake.ID

	// SetUnavailable sets the specified guildID as unavailable.
	SetUnavailable(guildID snowflake.ID)

	// SetAvailable sets the specified guildID as available.
	SetAvailable(guildID snowflake.ID)

	// IsUnavailable returns a bool indicating if the specified guildID is unavailable or not.
	IsUnavailable(guildID snowflake.ID) bool

	// UnavailableGuilds returns all guildIDs that are unavailable.
	UnavailableGuilds() []snowflake.ID
}

// NewGuildCache a new GuildCacheImpl with the given flags and policy.
// GuildCacheImpl is thread safe and can be used in multiple goroutines.
func NewGuildCache(flags Flags, policy Policy[discord.Guild]) GuildCache {
	return &GuildCacheImpl{
		Cache:             NewCache[discord.Guild](flags, FlagGuilds, policy),
		unreadyGuilds:     map[int]map[snowflake.ID]struct{}{},
		unavailableGuilds: map[snowflake.ID]struct{}{},
	}
}

// GuildCacheImpl is a thread safe GuildCache implementation.
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
