package core

import (
	"sync"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/internal/rwsync"
)

type CacheFindFunc[T any] func(T) bool

type Cache[T any] interface {
	RWLocker() rwsync.RWLocker
	Get(id discord.Snowflake) (T, bool)
	Put(id discord.Snowflake, entity T) T
	Remove(id discord.Snowflake) *T

	Cache() map[discord.Snowflake]T
	All() []T

	FindFirst(cacheFindFunc CacheFindFunc[T]) *T
	FindAll(cacheFindFunc CacheFindFunc[T]) []T
}

var _ Cache[any] = (*DefaultCache[any])(nil)

func NewCache[T any](flags CacheFlags, neededFlags CacheFlags) Cache[T] {
	return &DefaultCache[T]{
		flags:       flags,
		neededFlags: neededFlags,
		mu:          &sync.RWMutex{},
		entities:    make(map[discord.Snowflake]T),
	}
}

func NewCacheWithPolicy[T any](policy CachePolicy[T]) Cache[T] {
	return &DefaultCache[T]{
		policy:   policy,
		mu:       &sync.RWMutex{},
		entities: make(map[discord.Snowflake]T),
	}
}

type DefaultCache[T any, C CacheFlags] struct {
	flags       CacheFlags
	neededFlags CacheFlags
	policy      CachePolicy[T]
	mu          *sync.RWMutex
	entities    map[discord.Snowflake]T
}

func (c *DefaultCache[T]) RWLocker() rwsync.RWLocker {
	return c.mu
}

func (c *DefaultCache[T]) Get(id discord.Snowflake) (T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entity, ok := c.entities[id]
	return entity, ok
}

func (c *DefaultCache[T, C]) Put(id discord.Snowflake, entity T) T {
	if c.neededFlags != CacheFlagsNone && c.flags.Missing(c.neededFlags) {
		return entity
	}
	if c.policy != nil && !c.policy(entity) {
		return entity
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entities[id] = entity
	return entity
}

func (c *DefaultCache[T]) Remove(id discord.Snowflake) *T {
	c.mu.Lock()
	defer c.mu.Unlock()
	entity := c.entities[id]
	delete(c.entities, id)
	return &entity
}

func (c *DefaultCache[T]) Cache() map[discord.Snowflake]T {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.entities
}

func (c *DefaultCache[T]) All() []T {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entities := make([]T, 0, len(c.entities))
	for _, entity := range c.entities {
		entities = append(entities, entity)
	}
	return entities
}

func (c *DefaultCache[T]) FindFirst(cacheFindFunc CacheFindFunc[T]) *T {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, entity := range c.entities {
		if cacheFindFunc(entity) {
			return &entity
		}
	}
	return nil
}

func (c *DefaultCache[T]) FindAll(cacheFindFunc CacheFindFunc[T]) []T {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entities := make([]T, 0, len(c.entities))
	for _, entity := range c.entities {
		if cacheFindFunc(entity) {
			entities = append(entities, entity)
		}
	}
	return entities
}
