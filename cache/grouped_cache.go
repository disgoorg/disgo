package cache

import (
	"iter"

	"github.com/disgoorg/snowflake/v2"
)

type GroupedEntity[T any] struct {
	GroupID snowflake.ID
	Entity  T
}

// GroupedFilterFunc is used to filter grouped cached entities.
type GroupedFilterFunc[T any] func(groupID snowflake.ID, entity T) bool

// GroupedCache is a simple key value store grouped by a snowflake.ID. They key is always a snowflake.ID.
// The cache provides a simple way to store and retrieve entities. But is not guaranteed to be thread safe as this depends on the underlying implementation.
type GroupedCache[T any] interface {
	// Get returns a copy of the entity with the given groupID and ID. Returns ErrNotFound if the entity is not found.
	Get(groupID snowflake.ID, id snowflake.ID, opts ...AccessOpt) (T, error)

	// Put stores the given entity with the given groupID and ID as key. If the entity is already present, it will be overwritten.
	Put(groupID snowflake.ID, id snowflake.ID, entity T, opts ...AccessOpt) error

	// Remove removes the entity with the given groupID and ID as key and returns a copy of the entity. Returns ErrNotFound if the entity is not found.
	Remove(groupID snowflake.ID, id snowflake.ID, opts ...AccessOpt) (T, error)

	// GroupRemove removes all entities in the given groupID.
	GroupRemove(groupID snowflake.ID, opts ...AccessOpt) error

	// RemoveIf removes all entities that pass the given GroupedFilterFunc.
	RemoveIf(filterFunc GroupedFilterFunc[T], opts ...AccessOpt) error

	// GroupRemoveIf removes all entities that pass the given GroupedFilterFunc within the groupID.
	GroupRemoveIf(groupID snowflake.ID, filterFunc GroupedFilterFunc[T], opts ...AccessOpt) error

	// Len returns the total number of entities in the cache.
	Len(opts ...AccessOpt) (int, error)

	// GroupLen returns the number of entities in the cache within the groupID.
	GroupLen(groupID snowflake.ID, opts ...AccessOpt) (int, error)

	// All returns an [iter.Seq2] of all entities in the cache.
	All(opts ...AccessOpt) (iter.Seq2[GroupedEntity[T], error], error)

	// GroupAll returns an [iter.Seq] of all entities in the cache within the groupID.
	GroupAll(groupID snowflake.ID, opts ...AccessOpt) (iter.Seq[T], error)
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
	mu          RWMutex
	flags       Flags
	neededFlags Flags
	policy      Policy[T]
	cache       map[snowflake.ID]map[snowflake.ID]T
}

func (c *defaultGroupedCache[T]) Get(groupID snowflake.ID, id snowflake.ID, opts ...AccessOpt) (T, error) {
	var zero T
	cfg := resolveAccessConfig(opts)

	if err := c.mu.RLock(cfg.Ctx); err != nil {
		return zero, err
	}
	defer c.mu.RUnlock()

	if groupEntities, ok := c.cache[groupID]; ok {
		if entity, ok := groupEntities[id]; ok {
			return entity, nil
		}
	}

	return zero, ErrNotFound
}

func (c *defaultGroupedCache[T]) Put(groupID snowflake.ID, id snowflake.ID, entity T, opts ...AccessOpt) error {
	cfg := resolveAccessConfig(opts)

	if c.flags.Missing(c.neededFlags) {
		return nil
	}
	if c.policy != nil && !c.policy(entity) {
		return nil
	}
	if err := c.mu.Lock(cfg.Ctx); err != nil {
		return err
	}
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
	return nil
}

func (c *defaultGroupedCache[T]) Remove(groupID snowflake.ID, id snowflake.ID, opts ...AccessOpt) (T, error) {
	var zero T
	cfg := resolveAccessConfig(opts)

	if err := c.mu.Lock(cfg.Ctx); err != nil {
		return zero, err
	}
	defer c.mu.Unlock()

	if groupEntities, ok := c.cache[groupID]; ok {
		if entity, ok := groupEntities[id]; ok {
			delete(groupEntities, id)
			return entity, nil
		}
	}
	return zero, ErrNotFound
}

func (c *defaultGroupedCache[T]) GroupRemove(groupID snowflake.ID, opts ...AccessOpt) error {
	cfg := resolveAccessConfig(opts)

	if err := c.mu.Lock(cfg.Ctx); err != nil {
		return err
	}
	defer c.mu.Unlock()

	delete(c.cache, groupID)
	return nil
}

func (c *defaultGroupedCache[T]) RemoveIf(filterFunc GroupedFilterFunc[T], opts ...AccessOpt) error {
	cfg := resolveAccessConfig(opts)

	if err := c.mu.Lock(cfg.Ctx); err != nil {
		return err
	}
	defer c.mu.Unlock()

	for groupID := range c.cache {
		for id, entity := range c.cache[groupID] {
			if filterFunc(groupID, entity) {
				delete(c.cache[groupID], id)
			}
		}
	}
	return nil
}

func (c *defaultGroupedCache[T]) GroupRemoveIf(groupID snowflake.ID, filterFunc GroupedFilterFunc[T], opts ...AccessOpt) error {
	cfg := resolveAccessConfig(opts)

	if err := c.mu.Lock(cfg.Ctx); err != nil {
		return err
	}
	defer c.mu.Unlock()

	if groupEntities, ok := c.cache[groupID]; ok {
		for id, entity := range groupEntities {
			if filterFunc(groupID, entity) {
				delete(c.cache[groupID], id)
			}
		}
	}
	return nil
}

func (c *defaultGroupedCache[T]) Len(opts ...AccessOpt) (int, error) {
	cfg := resolveAccessConfig(opts)

	if err := c.mu.RLock(cfg.Ctx); err != nil {
		return 0, err
	}
	defer c.mu.RUnlock()

	var totalLen int
	for _, groupEntities := range c.cache {
		totalLen += len(groupEntities)
	}
	return totalLen, nil
}

func (c *defaultGroupedCache[T]) GroupLen(groupID snowflake.ID, opts ...AccessOpt) (int, error) {
	cfg := resolveAccessConfig(opts)

	if err := c.mu.RLock(cfg.Ctx); err != nil {
		return 0, err
	}
	defer c.mu.RUnlock()
	if groupEntities, ok := c.cache[groupID]; ok {
		return len(groupEntities), nil
	}
	return 0, nil
}

func (c *defaultGroupedCache[T]) All(opts ...AccessOpt) (iter.Seq2[GroupedEntity[T], error], error) {
	cfg := resolveAccessConfig(opts)

	return func(yield func(GroupedEntity[T], error) bool) {
		if err := c.mu.RLock(cfg.Ctx); err != nil {
			return
		}
		defer c.mu.RUnlock()

		for groupID, groupEntities := range c.cache {
			for _, entity := range groupEntities {
				if !yield(GroupedEntity[T]{GroupID: groupID, Entity: entity}, nil) {
					return
				}
			}
		}
	}, nil
}

func (c *defaultGroupedCache[T]) GroupAll(groupID snowflake.ID, opts ...AccessOpt) (iter.Seq[T], error) {
	cfg := resolveAccessConfig(opts)

	if err := c.mu.RLock(cfg.Ctx); err != nil {
		return nil, err
	}

	return func(yield func(T) bool) {
		defer c.mu.RUnlock()

		if groupEntities, ok := c.cache[groupID]; ok {
			for _, entity := range groupEntities {
				if !yield(entity) {
					return
				}
			}
		}
	}, nil
}
