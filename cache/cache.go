package cache

import (
	"sync"

	"github.com/disgoorg/snowflake/v2"
)

// FilterFunc is used to filter cached entities.
type FilterFunc[T any] func(T) bool

// Cache is a simple key value store. They key is always a snowflake.ID.
// The cache provides a simple way to store and retrieve entities. But is not guaranteed to be thread safe as this depends on the underlying implementation.
type Cache[T any] interface {
	// Get returns a copy of the entity with the given snowflake and a bool whether it was found or not.
	Get(id snowflake.ID) (T, bool)

	// Put stores the given entity with the given snowflake as key. If the entity is already present, it will be overwritten.
	Put(id snowflake.ID, entity T)

	// Remove removes the entity with the given snowflake as key and returns a copy of the entity and a bool whether it was removed or not.
	Remove(id snowflake.ID) (T, bool)

	// RemoveIf removes all entities that pass the given FilterFunc
	RemoveIf(filterFunc FilterFunc[T])

	// Len returns the number of entities in the cache.
	Len() int

	// ForEach calls the given function for each entity in the cache.
	ForEach(func(entity T))
}

var _ Cache[any] = (*DefaultCache[any])(nil)

// NewCache returns a new DefaultCache implementation which filter the entities after the gives Flags and Policy.
// This cache implementation is thread safe and can be used in multiple goroutines without any issues.
// It also only hands out copies to the entities. Regardless these entities should be handles as immutable.
func NewCache[T any](flags Flags, neededFlags Flags, policy Policy[T]) Cache[T] {
	return &DefaultCache[T]{
		flags:       flags,
		neededFlags: neededFlags,
		policy:      policy,
		cache:       make(map[snowflake.ID]T),
	}
}

// DefaultCache is a simple thread safe cache key value store.
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
	if c.flags.Missing(c.neededFlags) {
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

func (c *DefaultCache[T]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.cache)
}

func (c *DefaultCache[T]) ForEach(forEachFunc func(entity T)) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, entity := range c.cache {
		forEachFunc(entity)
	}
}
