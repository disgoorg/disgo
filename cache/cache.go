package cache

import (
	"sync"

	"github.com/disgoorg/snowflake/v2"
)

type FilterFunc[T any] func(T) bool

type Cache[T any] interface {
	Get(id snowflake.ID) (T, bool)
	Put(id snowflake.ID, entity T)
	Remove(id snowflake.ID) (T, bool)
	RemoveIf(filterFunc FilterFunc[T])

	All() []T
	MapAll() map[snowflake.ID]T

	FindFirst(cacheFindFunc FilterFunc[T]) (T, bool)
	FindAll(cacheFindFunc FilterFunc[T]) []T

	ForEach(func(entity T))
}

var _ Cache[any] = (*DefaultCache[any])(nil)

func NewCache[T any](flags Flags, neededFlags Flags, policy Policy[T]) Cache[T] {
	return &DefaultCache[T]{
		flags:       flags,
		neededFlags: neededFlags,
		policy:      policy,
		cache:       make(map[snowflake.ID]T),
	}
}

type DefaultCache[T any] struct {
	mu          sync.RWMutex
	flags       Flags
	neededFlags Flags
	policy      Policy[T]
	cache       map[snowflake.ID]T
}

func (c *DefaultCache[T]) Get(id snowflake.ID) (T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entity, ok := c.cache[id]
	return entity, ok
}

func (c *DefaultCache[T]) Put(id snowflake.ID, entity T) {
	if c.neededFlags != FlagsNone && c.flags.Missing(c.neededFlags) {
		return
	}
	if c.policy != nil && !c.policy(entity) {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[id] = entity
}

func (c *DefaultCache[T]) Remove(id snowflake.ID) (T, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entity, ok := c.cache[id]
	if ok {
		delete(c.cache, id)
	}
	return entity, ok
}

func (c *DefaultCache[T]) RemoveIf(filterFunc FilterFunc[T]) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for id, entity := range c.cache {
		if filterFunc(entity) {
			delete(c.cache, id)
		}
	}
}

func (c *DefaultCache[T]) All() []T {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entities := make([]T, len(c.cache))
	i := 0
	for _, entity := range c.cache {
		entities[i] = entity
		i++
	}
	return entities
}

func (c *DefaultCache[T]) MapAll() map[snowflake.ID]T {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entities := make(map[snowflake.ID]T, len(c.cache))
	for entityID, entity := range c.cache {
		entities[entityID] = entity
	}
	return entities
}

func (c *DefaultCache[T]) FindFirst(cacheFindFunc FilterFunc[T]) (T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, entity := range c.cache {
		if cacheFindFunc(entity) {
			return entity, true
		}
	}
	var entity T
	return entity, false
}

func (c *DefaultCache[T]) FindAll(cacheFindFunc FilterFunc[T]) []T {
	c.mu.RLock()
	defer c.mu.RUnlock()
	var entities []T
	for _, entity := range c.cache {
		if cacheFindFunc(entity) {
			entities = append(entities, entity)
		}
	}
	return entities
}

func (c *DefaultCache[T]) ForEach(forEachFunc func(entity T)) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, entity := range c.cache {
		forEachFunc(entity)
	}
}
