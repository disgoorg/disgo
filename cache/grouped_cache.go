package cache

import (
	"sync"

	"github.com/disgoorg/snowflake/v2"
)

// GroupedFilterFunc is used to filter grouped cached entities.
type GroupedFilterFunc[T any] func(groupID snowflake.ID, entity T) bool

// GroupedCache is a simple key value store grouped by a snowflake.ID. They key is always a snowflake.ID.
// The cache provides a simple way to store and retrieve entities. But is not guaranteed to be thread safe as this depends on the underlying implementation.
type GroupedCache[T any] interface {
	// Get returns a copy of the entity with the given groupID and ID and a bool wheaten it was found or not.
	Get(groupID snowflake.ID, id snowflake.ID) (T, bool)

	// Put stores the given entity with the given groupID and ID as key. If the entity is already present, it will be overwritten.
	Put(groupID snowflake.ID, id snowflake.ID, entity T)

	// Remove removes the entity with the given groupID and ID as key and returns a copy of the entity and a bool whether it was removed or not.
	Remove(groupID snowflake.ID, id snowflake.ID) (T, bool)

	// GroupRemove removes all entities in the given groupID.
	GroupRemove(groupID snowflake.ID)

	// RemoveIf removes all entities that pass the given GroupedFilterFunc.
	RemoveIf(filterFunc GroupedFilterFunc[T])

	// GroupRemoveIf removes all entities that pass the given GroupedFilterFunc within the groupID.
	GroupRemoveIf(groupID snowflake.ID, filterFunc GroupedFilterFunc[T])

	// Len returns the total number of entities in the cache.
	Len() int

	// GroupLen returns the number of entities in the cache within the groupID.
	GroupLen(groupID snowflake.ID) int

	// ForEach calls the given function for each entity in the cache.
	ForEach(func(groupID snowflake.ID, entity T))

	// GroupForEach calls the given function for each entity in the cache within the groupID.
	GroupForEach(groupID snowflake.ID, forEachFunc func(entity T))
}

var _ GroupedCache[any] = (*defaultGroupedCache[any])(nil)

// NewGroupedCache returns a new default GroupedCache with the provided flags, neededFlags and policy.
func NewGroupedCache[T any](flags Flags, neededFlags Flags, policy Policy[T]) GroupedCache[T] {
	return &defaultGroupedCache[T]{
		flags:       flags,
		neededFlags: neededFlags,
		policy:      policy,
		cache:       make(map[snowflake.ID]map[snowflake.ID]T),
	}
}

type defaultGroupedCache[T any] struct {
	mu          sync.RWMutex
	flags       Flags
	neededFlags Flags
	policy      Policy[T]
	cache       map[snowflake.ID]map[snowflake.ID]T
}

func (c *defaultGroupedCache[T]) Get(groupID snowflake.ID, id snowflake.ID) (T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if groupEntities, ok := c.cache[groupID]; ok {
		if entity, ok := groupEntities[id]; ok {
			return entity, true
		}
	}

	var entity T
	return entity, false
}

func (c *defaultGroupedCache[T]) Put(groupID snowflake.ID, id snowflake.ID, entity T) {
	if c.flags.Missing(c.neededFlags) {
		return
	}
	if c.policy != nil && !c.policy(entity) {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cache == nil {
		c.cache = make(map[snowflake.ID]map[snowflake.ID]T)
	}

	if groupEntities, ok := c.cache[groupID]; ok {
		groupEntities[id] = entity
	} else {
		groupEntities = make(map[snowflake.ID]T)
		groupEntities[id] = entity
		c.cache[groupID] = groupEntities
	}
}

func (c *defaultGroupedCache[T]) Remove(groupID snowflake.ID, id snowflake.ID) (entity T, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if groupEntities, ok := c.cache[groupID]; ok {
		if entity, ok := groupEntities[id]; ok {
			delete(groupEntities, id)
			return entity, ok
		}
	}
	ok = false
	return
}

func (c *defaultGroupedCache[T]) GroupRemove(groupID snowflake.ID) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.cache, groupID)
}

func (c *defaultGroupedCache[T]) RemoveIf(filterFunc GroupedFilterFunc[T]) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for groupID := range c.cache {
		for id, entity := range c.cache[groupID] {
			if filterFunc(groupID, entity) {
				delete(c.cache[groupID], id)
			}
		}
	}
}

func (c *defaultGroupedCache[T]) GroupRemoveIf(groupID snowflake.ID, filterFunc GroupedFilterFunc[T]) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if groupEntities, ok := c.cache[groupID]; ok {
		for id, entity := range groupEntities {
			if filterFunc(groupID, entity) {
				delete(c.cache[groupID], id)
			}
		}
	}
}

func (c *defaultGroupedCache[T]) Len() int {
	var totalLen int
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, groupEntities := range c.cache {
		totalLen += len(groupEntities)
	}
	return totalLen
}

func (c *defaultGroupedCache[T]) GroupLen(groupID snowflake.ID) int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if groupEntities, ok := c.cache[groupID]; ok {
		return len(groupEntities)
	}
	return 0
}

func (c *defaultGroupedCache[T]) ForEach(forEachFunc func(groupID snowflake.ID, entity T)) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for groupID, groupEntities := range c.cache {
		for _, entity := range groupEntities {
			forEachFunc(groupID, entity)
		}
	}
}
func (c *defaultGroupedCache[T]) GroupForEach(groupID snowflake.ID, forEachFunc func(entity T)) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, entity := range c.cache[groupID] {
		forEachFunc(entity)
	}
}
