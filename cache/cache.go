package cache

import (
	"errors"
	"iter"
	"sync"

	"github.com/disgoorg/snowflake/v2"
)

// ErrNotFound is returned when an entity is not found in the cache.
var ErrNotFound = errors.New("not found")

// FilterFunc is used to filter cached entities.
type FilterFunc[T any] func(T) bool

// Cache is a simple key value store. They key is always a snowflake.ID.
// The cache provides a simple way to store and retrieve entities. But is not guaranteed to be thread safe as this depends on the underlying implementation.
type Cache[T any] interface {
	// Get returns a copy of the entity with the given snowflake. Returns ErrNotFound if the entity is not found.
	Get(id snowflake.ID, opts ...AccessOpt) (T, error)

	// Put stores the given entity with the given snowflake as key. If the entity is already present, it will be overwritten.
	Put(id snowflake.ID, entity T, opts ...AccessOpt) error

	// Remove removes the entity with the given snowflake as key and returns a copy of the entity. Returns ErrNotFound if the entity is not found.
	Remove(id snowflake.ID, opts ...AccessOpt) (T, error)

	// RemoveIf removes all entities that pass the given FilterFunc
	RemoveIf(filterFunc FilterFunc[T], opts ...AccessOpt) error

	// Len returns the number of entities in the cache.
	Len(opts ...AccessOpt) (int, error)

	// All returns an [iter.Seq] of all entities in the cache.
	All(opts ...AccessOpt) (iter.Seq[T], error)
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

func (c *DefaultCache[T]) Get(id snowflake.ID, opts ...AccessOpt) (T, error) {
	var zero T
	cfg := &accessConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	if cfg.Ctx != nil {
		select {
		case <-cfg.Ctx.Done():
			return zero, cfg.Ctx.Err()
		default:
		}
	}

	c.mu.RLock()
	defer c.mu.RUnlock()
	entity, ok := c.cache[id]
	if !ok {
		return zero, ErrNotFound
	}
	return entity, nil
}

func (c *DefaultCache[T]) Put(id snowflake.ID, entity T, opts ...AccessOpt) error {
	cfg := &accessConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	if cfg.Ctx != nil {
		select {
		case <-cfg.Ctx.Done():
			return cfg.Ctx.Err()
		default:
		}
	}

	if c.flags.Missing(c.neededFlags) {
		return nil
	}
	if c.policy != nil && !c.policy(entity) {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[id] = entity
	return nil
}

func (c *DefaultCache[T]) Remove(id snowflake.ID, opts ...AccessOpt) (T, error) {
	var zero T
	cfg := &accessConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	if cfg.Ctx != nil {
		select {
		case <-cfg.Ctx.Done():
			return zero, cfg.Ctx.Err()
		default:
		}
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	entity, ok := c.cache[id]
	if !ok {
		return zero, ErrNotFound
	}
	delete(c.cache, id)
	return entity, nil
}

func (c *DefaultCache[T]) RemoveIf(filterFunc FilterFunc[T], opts ...AccessOpt) error {
	cfg := &accessConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	if cfg.Ctx != nil {
		select {
		case <-cfg.Ctx.Done():
			return cfg.Ctx.Err()
		default:
		}
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	for id, entity := range c.cache {
		if filterFunc(entity) {
			delete(c.cache, id)
		}
	}
	return nil
}

func (c *DefaultCache[T]) Len(opts ...AccessOpt) (int, error) {
	cfg := &accessConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	if cfg.Ctx != nil {
		select {
		case <-cfg.Ctx.Done():
			return 0, cfg.Ctx.Err()
		default:
		}
	}

	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.cache), nil
}

func (c *DefaultCache[T]) All(opts ...AccessOpt) (iter.Seq[T], error) {
	cfg := &accessConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	if cfg.Ctx != nil {
		select {
		case <-cfg.Ctx.Done():
			return nil, cfg.Ctx.Err()
		default:
		}
	}

	return func(yield func(T) bool) {
		c.mu.RLock()
		defer c.mu.RUnlock()
		for _, entity := range c.cache {
			if cfg.Ctx != nil {
				select {
				case <-cfg.Ctx.Done():
					return
				default:
				}
			}
			if !yield(entity) {
				break
			}
		}
	}, nil
}
