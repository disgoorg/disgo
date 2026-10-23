package cache

import (
	"iter"
	"sync"
	"time"

	"github.com/disgoorg/snowflake/v2"
)

// LRUCacheOpt is used to functionally configure an LRUCache.
type LRUCacheOpt[T any] func(c *lruConfig[T])

type lruConfig[T any] struct {
	ttl     time.Duration
	onEvict func(id snowflake.ID, entity T)
}

// WithLRUTTL sets the time-to-live for items in the LRUCache.
func WithLRUTTL[T any](ttl time.Duration) LRUCacheOpt[T] {
	return func(c *lruConfig[T]) {
		c.ttl = ttl
	}
}

// WithLRUOnEvict sets a callback function that is called when an entity is evicted from the cache.
func WithLRUOnEvict[T any](onEvict func(id snowflake.ID, entity T)) LRUCacheOpt[T] {
	return func(c *lruConfig[T]) {
		c.onEvict = onEvict
	}
}

type lruEntry[T any] struct {
	id        snowflake.ID
	entity    T
	expiresAt time.Time
	prev      *lruEntry[T]
	next      *lruEntry[T]
}

var _ Cache[any] = (*LRUCache[any])(nil)

// NewLRUCache returns a new LRUCache implementation with the given capacity and options.
// If capacity <= 0, the cache will have no size limit (can be used as pure TTL cache with WithLRUTTL).
func NewLRUCache[T any](flags Flags, neededFlags Flags, policy Policy[T], capacity int, opts ...LRUCacheOpt[T]) Cache[T] {
	var cfg lruConfig[T]
	for _, opt := range opts {
		opt(&cfg)
	}

	c := &LRUCache[T]{
		flags:       flags,
		neededFlags: neededFlags,
		policy:      policy,
		capacity:    capacity,
		ttl:         cfg.ttl,
		onEvict:     cfg.onEvict,
		cache:       make(map[snowflake.ID]*lruEntry[T]),
		head:        new(lruEntry[T]),
		tail:        new(lruEntry[T]),
	}
	c.head.next = c.tail
	c.tail.prev = c.head
	return c
}

// NewTTLCache returns an unbounded cache where items expire after the given TTL.
func NewTTLCache[T any](flags Flags, neededFlags Flags, policy Policy[T], ttl time.Duration, opts ...LRUCacheOpt[T]) Cache[T] {
	return NewLRUCache(flags, neededFlags, policy, 0, append([]LRUCacheOpt[T]{WithLRUTTL[T](ttl)}, opts...)...)
}

// LRUCache is a thread-safe LRU and TTL key-value cache implementing [Cache].
type LRUCache[T any] struct {
	mu          sync.RWMutex
	flags       Flags
	neededFlags Flags
	policy      Policy[T]
	capacity    int
	ttl         time.Duration
	onEvict     func(id snowflake.ID, entity T)
	cache       map[snowflake.ID]*lruEntry[T]
	head        *lruEntry[T]
	tail        *lruEntry[T]
}

func (c *LRUCache[T]) insertHead(entry *lruEntry[T]) {
	entry.next = c.head.next
	entry.prev = c.head
	c.head.next.prev = entry
	c.head.next = entry
}

func (c *LRUCache[T]) removeEntry(entry *lruEntry[T]) {
	entry.prev.next = entry.next
	entry.next.prev = entry.prev
	entry.prev = nil
	entry.next = nil
}

func (c *LRUCache[T]) moveToHead(entry *lruEntry[T]) {
	c.removeEntry(entry)
	c.insertHead(entry)
}

func (c *LRUCache[T]) isExpired(entry *lruEntry[T], now time.Time) bool {
	return c.ttl > 0 && !entry.expiresAt.IsZero() && now.After(entry.expiresAt)
}

func (c *LRUCache[T]) evictOldest() (T, bool) {
	oldest := c.tail.prev
	if oldest == c.head {
		var zero T
		return zero, false
	}
	c.removeEntry(oldest)
	delete(c.cache, oldest.id)
	if c.onEvict != nil {
		c.onEvict(oldest.id, oldest.entity)
	}
	return oldest.entity, true
}

func (c *LRUCache[T]) Get(id snowflake.ID) (T, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.cache[id]
	if !ok {
		var zero T
		return zero, false
	}

	now := time.Now()
	if c.isExpired(entry, now) {
		c.removeEntry(entry)
		delete(c.cache, id)
		if c.onEvict != nil {
			c.onEvict(id, entry.entity)
		}
		var zero T
		return zero, false
	}

	if c.capacity > 0 {
		c.moveToHead(entry)
	}

	return entry.entity, true
}

func (c *LRUCache[T]) Put(id snowflake.ID, entity T) {
	if c.flags.Missing(c.neededFlags) {
		return
	}
	if c.policy != nil && !c.policy(entity) {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	var expiresAt time.Time
	if c.ttl > 0 {
		expiresAt = now.Add(c.ttl)
	}

	if entry, ok := c.cache[id]; ok {
		entry.entity = entity
		entry.expiresAt = expiresAt
		if c.capacity > 0 {
			c.moveToHead(entry)
		}
		return
	}

	if c.capacity > 0 && len(c.cache) >= c.capacity {
		c.evictOldest()
	}

	entry := &lruEntry[T]{
		id:        id,
		entity:    entity,
		expiresAt: expiresAt,
	}
	c.cache[id] = entry
	c.insertHead(entry)
}

func (c *LRUCache[T]) Remove(id snowflake.ID) (T, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.cache[id]
	if !ok {
		var zero T
		return zero, false
	}

	c.removeEntry(entry)
	delete(c.cache, id)

	if c.isExpired(entry, time.Now()) {
		var zero T
		return zero, false
	}

	return entry.entity, true
}

func (c *LRUCache[T]) RemoveIf(filterFunc FilterFunc[T]) {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	curr := c.head.next
	for curr != c.tail {
		next := curr.next
		if c.isExpired(curr, now) || filterFunc(curr.entity) {
			c.removeEntry(curr)
			delete(c.cache, curr.id)
			if c.onEvict != nil {
				c.onEvict(curr.id, curr.entity)
			}
		}
		curr = next
	}
}

func (c *LRUCache[T]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.ttl == 0 {
		return len(c.cache)
	}

	now := time.Now()
	count := 0
	for _, entry := range c.cache {
		if !c.isExpired(entry, now) {
			count++
		}
	}
	return count
}

func (c *LRUCache[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		c.mu.RLock()
		defer c.mu.RUnlock()

		now := time.Now()
		curr := c.head.next
		for curr != c.tail {
			if !c.isExpired(curr, now) {
				if !yield(curr.entity) {
					break
				}
			}
			curr = curr.next
		}
	}
}
