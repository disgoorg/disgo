package core

import (
	"sync"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/internal/rwsync"
)

type GroupedCache[T any] interface {
	RWLocker() rwsync.RWLocker
	Get(groupID discord.Snowflake, id discord.Snowflake) (T, bool)
	Put(groupID discord.Snowflake, id discord.Snowflake, entity T) T
	Remove(groupID discord.Snowflake, id discord.Snowflake) (T, bool)

	Cache() map[discord.Snowflake]map[discord.Snowflake]T
	All() map[discord.Snowflake][]T

	GroupCache(groupID discord.Snowflake) map[discord.Snowflake]T
	GroupAll(groupID discord.Snowflake) []T

	FindFirst(cacheFindFunc CacheFindFunc[T]) T
	FindAll(cacheFindFunc CacheFindFunc[T]) []T
}

var _ GroupedCache[any] = (*DefaultGroupedCache[any])(nil)

func NewGroupedCache[T any](flags CacheFlags, neededFlags CacheFlags) GroupedCache[T] {
	return &DefaultGroupedCache[T]{
		flags:       flags,
		neededFlags: neededFlags,
		mu:          &sync.RWMutex{},
		cache:       make(map[discord.Snowflake]map[discord.Snowflake]T),
	}
}

func NewGroupedCacheWithPolicy[T any](policy CachePolicy[T]) Cache[T] {
	return &DefaultCache[T]{
		policy:   policy,
		mu:       &sync.RWMutex{},
		entities: make(map[discord.Snowflake]T),
	}
}

type DefaultGroupedCache[T any] struct {
	flags       CacheFlags
	neededFlags CacheFlags
	policy      CachePolicy[T]
	mu          *sync.RWMutex
	cache       map[discord.Snowflake]map[discord.Snowflake]T
}

func (c *DefaultGroupedCache[T]) RWLocker() rwsync.RWLocker {
	return c.mu
}

func (c *DefaultGroupedCache[T]) Get(groupID discord.Snowflake, id discord.Snowflake) (entity T, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if guildEntities, ok := c.cache[groupID]; ok {
		if entity, ok := guildEntities[id]; ok {
			return entity, true
		}
	}

	ok = false
	return
}

func (c *DefaultGroupedCache[T]) Put(groupID discord.Snowflake, id discord.Snowflake, entity T) T {
	if c.neededFlags != CacheFlagsNone && c.flags.Missing(c.neededFlags) {
		return entity
	}
	if c.policy != nil && !c.policy(entity) {
		return entity
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cache == nil {
		c.cache = make(map[discord.Snowflake]map[discord.Snowflake]T)
	}

	if guildEntities, ok := c.cache[groupID]; ok {
		guildEntities[id] = entity
	} else {
		guildEntities = make(map[discord.Snowflake]T)
		guildEntities[id] = entity
		c.cache[groupID] = guildEntities
	}

	return entity
}

func (c *DefaultGroupedCache[T]) Remove(groupID discord.Snowflake, id discord.Snowflake) (entity T, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if guildEntities, ok := c.cache[groupID]; ok {
		if entity, ok := guildEntities[id]; ok {
			delete(guildEntities, id)
			return entity, ok
		}
	}
	ok = false
	return
}

func (c *DefaultGroupedCache[T]) Cache() map[discord.Snowflake]map[discord.Snowflake]T {
	return c.cache
}

func (c *DefaultGroupedCache[T]) All() map[discord.Snowflake][]T {
	c.mu.RLock()
	defer c.mu.RUnlock()

	all := make(map[discord.Snowflake][]T)
	for groupID, guildEntities := range c.cache {
		all[groupID] = make([]T, 0, len(guildEntities))
		for _, entity := range guildEntities {
			all[groupID] = append(all[groupID], entity)
		}
	}

	return all
}

func (c *DefaultGroupedCache[T]) GroupCache(groupID discord.Snowflake) map[discord.Snowflake]T {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if guildEntities, ok := c.cache[groupID]; ok {
		return guildEntities
	}

	return nil
}

func (c *DefaultGroupedCache[T]) GroupAll(groupID discord.Snowflake) []T {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if guildEntities, ok := c.cache[groupID]; ok {
		all := make([]T, 0, len(guildEntities))
		for _, entity := range guildEntities {
			all = append(all, entity)
		}

		return all
	}

	return nil
}

func (c *DefaultGroupedCache[T]) FindFirst(cacheFindFunc CacheFindFunc[T]) (entity T) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, guildEntities := range c.cache {
		for _, entity = range guildEntities {
			if cacheFindFunc(entity) {
				return
			}
		}
	}

	return
}

func (c *DefaultGroupedCache[T]) FindAll(cacheFindFunc CacheFindFunc[T]) []T {
	c.mu.RLock()
	defer c.mu.RUnlock()

	all := make([]T, 0)
	for _, guildEntities := range c.cache {
		for _, entity := range guildEntities {
			if cacheFindFunc(entity) {
				all = append(all, entity)
			}
		}
	}

	return all
}
