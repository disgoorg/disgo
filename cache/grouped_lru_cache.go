package cache

import (
	"iter"
	"sync"
	"time"

	"github.com/disgoorg/snowflake/v2"
)

// GroupedLRUCacheOpt is used to functionally configure a GroupedLRUCache.
type GroupedLRUCacheOpt[T any] func(c *groupedLRUConfig[T])

type groupedLRUConfig[T any] struct {
	ttl       time.Duration
	maxGroups int
	onEvict   func(groupID snowflake.ID, id snowflake.ID, entity T)
}

// WithGroupedLRUTTL sets the time-to-live for items in the GroupedLRUCache.
// If ttl <= 0, TTL expiration is disabled and entries do not expire.
func WithGroupedLRUTTL[T any](ttl time.Duration) GroupedLRUCacheOpt[T] {
	return func(c *groupedLRUConfig[T]) {
		c.ttl = ttl
	}
}

// WithGroupedLRUMaxGroups sets the maximum number of groups retained in the cache.
// When exceeded, the least recently accessed group is evicted.
func WithGroupedLRUMaxGroups[T any](maxGroups int) GroupedLRUCacheOpt[T] {
	return func(c *groupedLRUConfig[T]) {
		c.maxGroups = maxGroups
	}
}

// WithGroupedLRUOnEvict sets a callback function that is called when an entity is evicted from the cache.
// Callbacks are always executed outside the cache mutex, allowing reentrant operations on the cache.
func WithGroupedLRUOnEvict[T any](onEvict func(groupID snowflake.ID, id snowflake.ID, entity T)) GroupedLRUCacheOpt[T] {
	return func(c *groupedLRUConfig[T]) {
		c.onEvict = onEvict
	}
}

type groupedLRUEntry[T any] struct {
	groupID   snowflake.ID
	id        snowflake.ID
	entity    T
	expiresAt time.Time
	prev      *groupedLRUEntry[T]
	next      *groupedLRUEntry[T]
}

type lruGroup[T any] struct {
	id    snowflake.ID
	cache map[snowflake.ID]*groupedLRUEntry[T]
	head  *groupedLRUEntry[T]
	tail  *groupedLRUEntry[T]
	prev  *lruGroup[T]
	next  *lruGroup[T]
}

func newLRUGroup[T any](id snowflake.ID) *lruGroup[T] {
	g := &lruGroup[T]{
		id:    id,
		cache: make(map[snowflake.ID]*groupedLRUEntry[T]),
		head:  new(groupedLRUEntry[T]),
		tail:  new(groupedLRUEntry[T]),
	}
	g.head.next = g.tail
	g.tail.prev = g.head
	return g
}

func (g *lruGroup[T]) insertHead(entry *groupedLRUEntry[T]) {
	entry.next = g.head.next
	entry.prev = g.head
	g.head.next.prev = entry
	g.head.next = entry
}

func (g *lruGroup[T]) removeEntry(entry *groupedLRUEntry[T]) {
	entry.prev.next = entry.next
	entry.next.prev = entry.prev
	entry.prev = nil
	entry.next = nil
}

func (g *lruGroup[T]) moveToHead(entry *groupedLRUEntry[T]) {
	g.removeEntry(entry)
	g.insertHead(entry)
}

func (g *lruGroup[T]) evictOldestLocked() (snowflake.ID, T, bool) {
	oldest := g.tail.prev
	if oldest == g.head {
		var zero T
		return 0, zero, false
	}
	g.removeEntry(oldest)
	delete(g.cache, oldest.id)
	return oldest.id, oldest.entity, true
}

type evictedGroupedItem[T any] struct {
	groupID snowflake.ID
	id      snowflake.ID
	entity  T
}

var _ GroupedCache[any] = (*GroupedLRUCache[any])(nil)

// NewGroupedLRUCache returns a new GroupedLRUCache with per-group LRU capacity and optional TTL.
// If maxItemsPerGroup <= 0, groups will have no item limit (can be used as pure grouped TTL cache with WithGroupedLRUTTL).
func NewGroupedLRUCache[T any](flags Flags, neededFlags Flags, policy Policy[T], maxItemsPerGroup int, opts ...GroupedLRUCacheOpt[T]) GroupedCache[T] {
	var cfg groupedLRUConfig[T]
	for _, opt := range opts {
		opt(&cfg)
	}

	c := &GroupedLRUCache[T]{
		flags:            flags,
		neededFlags:      neededFlags,
		policy:           policy,
		maxItemsPerGroup: maxItemsPerGroup,
		maxGroups:        cfg.maxGroups,
		ttl:              cfg.ttl,
		onEvict:          cfg.onEvict,
		groups:           make(map[snowflake.ID]*lruGroup[T]),
		groupHead:        new(lruGroup[T]),
		groupTail:        new(lruGroup[T]),
	}
	c.groupHead.next = c.groupTail
	c.groupTail.prev = c.groupHead
	return c
}

// NewGroupedTTLCache returns an unbounded grouped cache where items expire after the given TTL.
func NewGroupedTTLCache[T any](flags Flags, neededFlags Flags, policy Policy[T], ttl time.Duration, opts ...GroupedLRUCacheOpt[T]) GroupedCache[T] {
	return NewGroupedLRUCache(flags, neededFlags, policy, 0, append([]GroupedLRUCacheOpt[T]{WithGroupedLRUTTL[T](ttl)}, opts...)...)
}

// GroupedLRUCache is a thread-safe grouped LRU and TTL key-value cache implementing [GroupedCache].
type GroupedLRUCache[T any] struct {
	mu               sync.RWMutex
	flags            Flags
	neededFlags      Flags
	policy           Policy[T]
	maxItemsPerGroup int
	maxGroups        int
	ttl              time.Duration
	onEvict          func(groupID snowflake.ID, id snowflake.ID, entity T)
	groups           map[snowflake.ID]*lruGroup[T]
	groupHead        *lruGroup[T]
	groupTail        *lruGroup[T]
}

func (c *GroupedLRUCache[T]) insertGroupHead(g *lruGroup[T]) {
	g.next = c.groupHead.next
	g.prev = c.groupHead
	c.groupHead.next.prev = g
	c.groupHead.next = g
}

func (c *GroupedLRUCache[T]) removeGroup(g *lruGroup[T]) {
	g.prev.next = g.next
	g.next.prev = g.prev
	g.prev = nil
	g.next = nil
}

func (c *GroupedLRUCache[T]) moveGroupToHead(g *lruGroup[T]) {
	c.removeGroup(g)
	c.insertGroupHead(g)
}

func (c *GroupedLRUCache[T]) isExpired(entry *groupedLRUEntry[T], now time.Time) bool {
	return c.ttl > 0 && !entry.expiresAt.IsZero() && now.After(entry.expiresAt)
}

func (c *GroupedLRUCache[T]) Get(groupID snowflake.ID, id snowflake.ID) (T, bool) {
	c.mu.Lock()
	group, ok := c.groups[groupID]
	if !ok {
		c.mu.Unlock()
		var zero T
		return zero, false
	}

	entry, ok := group.cache[id]
	if !ok {
		c.mu.Unlock()
		var zero T
		return zero, false
	}

	now := time.Now()
	if c.isExpired(entry, now) {
		group.removeEntry(entry)
		delete(group.cache, id)
		entity := entry.entity
		c.mu.Unlock()

		if c.onEvict != nil {
			c.onEvict(groupID, id, entity)
		}
		var zero T
		return zero, false
	}

	if c.maxItemsPerGroup > 0 {
		group.moveToHead(entry)
	}
	if c.maxGroups > 0 {
		c.moveGroupToHead(group)
	}
	entity := entry.entity
	c.mu.Unlock()

	return entity, true
}

// Put stores the given entity with the given groupID and snowflake as key. If the entity is
// already present, its value and TTL expiration are updated, and it is marked as most-recently-used.
func (c *GroupedLRUCache[T]) Put(groupID snowflake.ID, id snowflake.ID, entity T) {
	if c.flags.Missing(c.neededFlags) {
		return
	}
	if c.policy != nil && !c.policy(entity) {
		return
	}

	c.mu.Lock()
	now := time.Now()
	var expiresAt time.Time
	if c.ttl > 0 {
		expiresAt = now.Add(c.ttl)
	}

	var evicted []evictedGroupedItem[T]

	group, ok := c.groups[groupID]
	if !ok {
		if c.maxGroups > 0 && len(c.groups) >= c.maxGroups {
			// Evict oldest group
			oldestGroup := c.groupTail.prev
			if oldestGroup != c.groupHead {
				c.removeGroup(oldestGroup)
				delete(c.groups, oldestGroup.id)
				if c.onEvict != nil {
					for _, e := range oldestGroup.cache {
						evicted = append(evicted, evictedGroupedItem[T]{
							groupID: oldestGroup.id,
							id:      e.id,
							entity:  e.entity,
						})
					}
				}
			}
		}
		group = newLRUGroup[T](groupID)
		c.groups[groupID] = group
		c.insertGroupHead(group)
	} else if c.maxGroups > 0 {
		c.moveGroupToHead(group)
	}

	if entry, ok := group.cache[id]; ok {
		entry.entity = entity
		entry.expiresAt = expiresAt
		if c.maxItemsPerGroup > 0 {
			group.moveToHead(entry)
		}
		c.mu.Unlock()
		return
	}

	if c.maxItemsPerGroup > 0 && len(group.cache) >= c.maxItemsPerGroup {
		if eID, eEntity, ok := group.evictOldestLocked(); ok && c.onEvict != nil {
			evicted = append(evicted, evictedGroupedItem[T]{
				groupID: groupID,
				id:      eID,
				entity:  eEntity,
			})
		}
	}

	entry := &groupedLRUEntry[T]{
		groupID:   groupID,
		id:        id,
		entity:    entity,
		expiresAt: expiresAt,
	}
	group.cache[id] = entry
	group.insertHead(entry)
	c.mu.Unlock()

	for _, item := range evicted {
		c.onEvict(item.groupID, item.id, item.entity)
	}
}

// Remove removes the entity with the given groupID and snowflake as key and returns a copy of the
// entity and a bool whether it was removed or not. If an eviction callback is configured, it will be
// invoked outside the lock after removal.
func (c *GroupedLRUCache[T]) Remove(groupID snowflake.ID, id snowflake.ID) (T, bool) {
	c.mu.Lock()
	group, ok := c.groups[groupID]
	if !ok {
		c.mu.Unlock()
		var zero T
		return zero, false
	}

	entry, ok := group.cache[id]
	if !ok {
		c.mu.Unlock()
		var zero T
		return zero, false
	}

	group.removeEntry(entry)
	delete(group.cache, id)
	expired := c.isExpired(entry, time.Now())
	entity := entry.entity
	c.mu.Unlock()

	if c.onEvict != nil {
		c.onEvict(groupID, id, entity)
	}

	if expired {
		var zero T
		return zero, false
	}
	return entity, true
}

func (c *GroupedLRUCache[T]) GroupRemove(groupID snowflake.ID) {
	c.mu.Lock()
	group, ok := c.groups[groupID]
	if !ok {
		c.mu.Unlock()
		return
	}

	c.removeGroup(group)
	delete(c.groups, groupID)

	var evicted []evictedGroupedItem[T]
	if c.onEvict != nil {
		for _, e := range group.cache {
			evicted = append(evicted, evictedGroupedItem[T]{
				groupID: groupID,
				id:      e.id,
				entity:  e.entity,
			})
		}
	}
	c.mu.Unlock()

	for _, item := range evicted {
		c.onEvict(item.groupID, item.id, item.entity)
	}
}

func (c *GroupedLRUCache[T]) RemoveIf(filterFunc GroupedFilterFunc[T]) {
	c.mu.Lock()
	now := time.Now()
	var evicted []evictedGroupedItem[T]
	for groupID, group := range c.groups {
		curr := group.head.next
		for curr != group.tail {
			next := curr.next
			if c.isExpired(curr, now) || filterFunc(groupID, curr.entity) {
				group.removeEntry(curr)
				delete(group.cache, curr.id)
				if c.onEvict != nil {
					evicted = append(evicted, evictedGroupedItem[T]{
						groupID: groupID,
						id:      curr.id,
						entity:  curr.entity,
					})
				}
			}
			curr = next
		}
	}
	c.mu.Unlock()

	for _, item := range evicted {
		c.onEvict(item.groupID, item.id, item.entity)
	}
}

func (c *GroupedLRUCache[T]) GroupRemoveIf(groupID snowflake.ID, filterFunc GroupedFilterFunc[T]) {
	c.mu.Lock()
	group, ok := c.groups[groupID]
	if !ok {
		c.mu.Unlock()
		return
	}

	now := time.Now()
	var evicted []evictedGroupedItem[T]
	curr := group.head.next
	for curr != group.tail {
		next := curr.next
		if c.isExpired(curr, now) || filterFunc(groupID, curr.entity) {
			group.removeEntry(curr)
			delete(group.cache, curr.id)
			if c.onEvict != nil {
				evicted = append(evicted, evictedGroupedItem[T]{
					groupID: groupID,
					id:      curr.id,
					entity:  curr.entity,
				})
			}
		}
		curr = next
	}
	c.mu.Unlock()

	for _, item := range evicted {
		c.onEvict(item.groupID, item.id, item.entity)
	}
}

func (c *GroupedLRUCache[T]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.ttl <= 0 {
		total := 0
		for _, group := range c.groups {
			total += len(group.cache)
		}
		return total
	}

	now := time.Now()
	total := 0
	for _, group := range c.groups {
		for _, entry := range group.cache {
			if !c.isExpired(entry, now) {
				total++
			}
		}
	}
	return total
}

func (c *GroupedLRUCache[T]) GroupLen(groupID snowflake.ID) int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	group, ok := c.groups[groupID]
	if !ok {
		return 0
	}

	if c.ttl <= 0 {
		return len(group.cache)
	}

	now := time.Now()
	count := 0
	for _, entry := range group.cache {
		if !c.isExpired(entry, now) {
			count++
		}
	}
	return count
}

func (c *GroupedLRUCache[T]) All() iter.Seq2[snowflake.ID, T] {
	return func(yield func(snowflake.ID, T) bool) {
		c.mu.RLock()
		now := time.Now()
		type pair struct {
			groupID snowflake.ID
			entity  T
		}
		var items []pair
		gCurr := c.groupHead.next
		for gCurr != c.groupTail {
			curr := gCurr.head.next
			for curr != gCurr.tail {
				if !c.isExpired(curr, now) {
					items = append(items, pair{groupID: curr.groupID, entity: curr.entity})
				}
				curr = curr.next
			}
			gCurr = gCurr.next
		}
		c.mu.RUnlock()

		for _, item := range items {
			if !yield(item.groupID, item.entity) {
				return
			}
		}
	}
}

func (c *GroupedLRUCache[T]) GroupAll(groupID snowflake.ID) iter.Seq[T] {
	return func(yield func(T) bool) {
		c.mu.RLock()
		group, ok := c.groups[groupID]
		if !ok {
			c.mu.RUnlock()
			return
		}

		now := time.Now()
		var items []T
		curr := group.head.next
		for curr != group.tail {
			if !c.isExpired(curr, now) {
				items = append(items, curr.entity)
			}
			curr = curr.next
		}
		c.mu.RUnlock()

		for _, item := range items {
			if !yield(item) {
				return
			}
		}
	}
}
