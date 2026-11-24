package cache

import (
	"sync"
	"testing"
	"time"

	"github.com/disgoorg/snowflake/v2"
)

type testEntity struct {
	ID   snowflake.ID
	Name string
}

func TestLRUCache_Basic(t *testing.T) {
	c := NewLRUCache[testEntity](FlagsAll, FlagsAll, nil, 10)

	id1 := snowflake.ID(1)
	id2 := snowflake.ID(2)

	c.Put(id1, testEntity{ID: id1, Name: "one"})
	c.Put(id2, testEntity{ID: id2, Name: "two"})

	if c.Len() != 2 {
		t.Fatalf("expected len 2, got %d", c.Len())
	}

	e1, ok := c.Get(id1)
	if !ok || e1.Name != "one" {
		t.Fatalf("expected to find id1 with Name 'one', got ok=%v, entity=%+v", ok, e1)
	}

	removed, ok := c.Remove(id1)
	if !ok || removed.Name != "one" {
		t.Fatalf("expected removed ok with Name 'one', got ok=%v", ok)
	}

	if c.Len() != 1 {
		t.Fatalf("expected len 1, got %d", c.Len())
	}

	_, ok = c.Get(id1)
	if ok {
		t.Fatalf("expected id1 to be removed")
	}
}

func TestLRUCache_Eviction(t *testing.T) {
	var evicted []snowflake.ID
	c := NewLRUCache[testEntity](FlagsAll, FlagsAll, nil, 3, WithLRUOnEvict[testEntity](func(id snowflake.ID, entity testEntity) {
		evicted = append(evicted, id)
	}))

	c.Put(1, testEntity{ID: 1, Name: "1"})
	c.Put(2, testEntity{ID: 2, Name: "2"})
	c.Put(3, testEntity{ID: 3, Name: "3"})

	// Access 1 so it becomes most recently used: order MRU -> LRU is now: 1, 3, 2
	c.Get(1)

	// Insert 4, which should evict 2 (the oldest)
	c.Put(4, testEntity{ID: 4, Name: "4"})

	if c.Len() != 3 {
		t.Fatalf("expected len 3, got %d", c.Len())
	}

	if _, ok := c.Get(2); ok {
		t.Fatalf("expected id 2 to be evicted")
	}

	if len(evicted) != 1 || evicted[0] != 2 {
		t.Fatalf("expected evicted [2], got %v", evicted)
	}

	// Verify 1, 3, 4 are still present
	for _, id := range []snowflake.ID{1, 3, 4} {
		if _, ok := c.Get(id); !ok {
			t.Fatalf("expected id %d to be present", id)
		}
	}
}

func TestLRUCache_TTL(t *testing.T) {
	c := NewTTLCache[testEntity](FlagsAll, FlagsAll, nil, 50*time.Millisecond)

	c.Put(1, testEntity{ID: 1, Name: "temp"})

	if _, ok := c.Get(1); !ok {
		t.Fatalf("expected id 1 to be present immediately")
	}

	time.Sleep(70 * time.Millisecond)

	if _, ok := c.Get(1); ok {
		t.Fatalf("expected id 1 to be expired after TTL")
	}

	if c.Len() != 0 {
		t.Fatalf("expected len 0 after expiry, got %d", c.Len())
	}
}

func TestLRUCache_PutRefreshesTTL(t *testing.T) {
	c := NewTTLCache[testEntity](FlagsAll, FlagsAll, nil, 80*time.Millisecond)

	c.Put(1, testEntity{ID: 1, Name: "v1"})
	time.Sleep(50 * time.Millisecond)

	// Updating the entry resets expiration and updates value
	c.Put(1, testEntity{ID: 1, Name: "v2"})
	time.Sleep(50 * time.Millisecond)

	// 100ms total elapsed, but TTL was refreshed at 50ms, so it should still be alive
	e, ok := c.Get(1)
	if !ok || e.Name != "v2" {
		t.Fatalf("expected item to be alive with updated value, got ok=%v, entity=%+v", ok, e)
	}

	time.Sleep(50 * time.Millisecond)
	if _, ok := c.Get(1); ok {
		t.Fatalf("expected item to expire after refreshed TTL")
	}
}

func TestLRUCache_ZeroOrNegativeTTL(t *testing.T) {
	c := NewLRUCache[testEntity](FlagsAll, FlagsAll, nil, 10, WithLRUTTL[testEntity](-1*time.Minute))

	c.Put(1, testEntity{ID: 1, Name: "eternal"})
	time.Sleep(20 * time.Millisecond)

	if e, ok := c.Get(1); !ok || e.Name != "eternal" {
		t.Fatalf("expected item to persist with negative TTL, got ok=%v", ok)
	}
	if c.Len() != 1 {
		t.Fatalf("expected len 1, got %d", c.Len())
	}
}

func TestLRUCache_RemoveCallsOnEvict(t *testing.T) {
	var evictedID snowflake.ID
	var evictedEntity testEntity
	c := NewLRUCache[testEntity](FlagsAll, FlagsAll, nil, 10, WithLRUOnEvict[testEntity](func(id snowflake.ID, entity testEntity) {
		evictedID = id
		evictedEntity = entity
	}))

	c.Put(42, testEntity{ID: 42, Name: "answer"})
	removed, ok := c.Remove(42)
	if !ok || removed.Name != "answer" {
		t.Fatalf("expected successful remove")
	}

	if evictedID != 42 || evictedEntity.Name != "answer" {
		t.Fatalf("expected onEvict to be invoked on Remove, got id=%v entity=%+v", evictedID, evictedEntity)
	}
}

func TestLRUCache_ReentrantCallback(t *testing.T) {
	// Verify that onEvict callback can safely perform Get, Put, and Remove operations on the
	// same cache without causing a mutex deadlock.
	var c Cache[testEntity]
	var reentrantCalled bool
	done := make(chan struct{})

	c = NewLRUCache[testEntity](FlagsAll, FlagsAll, nil, 2, WithLRUOnEvict[testEntity](func(id snowflake.ID, entity testEntity) {
		if !reentrantCalled {
			reentrantCalled = true
			// Reentrant read
			_, _ = c.Get(id)
			// Reentrant write/put
			c.Put(1000, testEntity{ID: 1000, Name: "reentrant"})
			// Reentrant remove
			_, _ = c.Remove(1000)
		}
	}))

	go func() {
		c.Put(1, testEntity{ID: 1, Name: "one"})
		c.Put(2, testEntity{ID: 2, Name: "two"})
		// Trigger eviction of 1
		c.Put(3, testEntity{ID: 3, Name: "three"})
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("deadlock detected: onEvict callback failed reentrancy")
	}
}

func TestLRUCache_IteratorMutation(t *testing.T) {
	// Verify that mutating the cache inside an All() iteration loop does not deadlock.
	c := NewLRUCache[testEntity](FlagsAll, FlagsAll, nil, 10)
	for i := 1; i <= 5; i++ {
		c.Put(snowflake.ID(i), testEntity{ID: snowflake.ID(i), Name: "val"})
	}

	done := make(chan struct{})
	go func() {
		for entity := range c.All() {
			c.Remove(entity.ID)
		}
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("deadlock detected: mutating cache during All() iteration")
	}

	if c.Len() != 0 {
		t.Fatalf("expected len 0 after removing all entries, got %d", c.Len())
	}
}

func TestLRUCache_All(t *testing.T) {
	c := NewLRUCache[testEntity](FlagsAll, FlagsAll, nil, 5)

	c.Put(1, testEntity{ID: 1, Name: "one"})
	c.Put(2, testEntity{ID: 2, Name: "two"})
	c.Put(3, testEntity{ID: 3, Name: "three"})

	count := 0
	for entity := range c.All() {
		count++
		if entity.ID == 0 {
			t.Fatalf("unexpected empty entity in All()")
		}
	}
	if count != 3 {
		t.Fatalf("expected 3 items in All(), got %d", count)
	}
}

func TestLRUCache_Concurrency(t *testing.T) {
	c := NewLRUCache[testEntity](FlagsAll, FlagsAll, nil, 100, WithLRUTTL[testEntity](200*time.Millisecond))

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(2)
		val := i
		go func() {
			defer wg.Done()
			id := snowflake.ID(val)
			c.Put(id, testEntity{ID: id, Name: "concurrency"})
		}()
		go func() {
			defer wg.Done()
			id := snowflake.ID(val)
			c.Get(id)
		}()
	}
	wg.Wait()
}

func TestGroupedLRUCache_Basic(t *testing.T) {
	c := NewGroupedLRUCache[testEntity](FlagsAll, FlagsAll, nil, 2)

	guild1 := snowflake.ID(100)
	guild2 := snowflake.ID(200)

	c.Put(guild1, 1, testEntity{ID: 1, Name: "g1-m1"})
	c.Put(guild1, 2, testEntity{ID: 2, Name: "g1-m2"})
	c.Put(guild2, 1, testEntity{ID: 1, Name: "g2-m1"})

	if c.Len() != 3 {
		t.Fatalf("expected total len 3, got %d", c.Len())
	}
	if c.GroupLen(guild1) != 2 {
		t.Fatalf("expected guild1 len 2, got %d", c.GroupLen(guild1))
	}
	if c.GroupLen(guild2) != 1 {
		t.Fatalf("expected guild2 len 1, got %d", c.GroupLen(guild2))
	}

	e, ok := c.Get(guild1, 1)
	if !ok || e.Name != "g1-m1" {
		t.Fatalf("expected g1-m1, got ok=%v, %+v", ok, e)
	}

	// Insert third into guild1 (limit is 2) -> oldest (2, since 1 was accessed) will be evicted
	c.Put(guild1, 3, testEntity{ID: 3, Name: "g1-m3"})

	if _, ok := c.Get(guild1, 2); ok {
		t.Fatalf("expected guild1 user 2 to be evicted")
	}

	if _, ok := c.Get(guild1, 1); !ok {
		t.Fatalf("expected guild1 user 1 to remain")
	}
	if _, ok := c.Get(guild1, 3); !ok {
		t.Fatalf("expected guild1 user 3 to remain")
	}
}

func TestGroupedLRUCache_RemoveCallsOnEvict(t *testing.T) {
	var evicted []snowflake.ID
	c := NewGroupedLRUCache[testEntity](FlagsAll, FlagsAll, nil, 10,
		WithGroupedLRUOnEvict[testEntity](func(groupID snowflake.ID, id snowflake.ID, entity testEntity) {
			evicted = append(evicted, id)
		}),
	)

	guild := snowflake.ID(99)
	c.Put(guild, 1, testEntity{ID: 1, Name: "one"})
	c.Put(guild, 2, testEntity{ID: 2, Name: "two"})

	// Remove single entry
	c.Remove(guild, 1)
	if len(evicted) != 1 || evicted[0] != 1 {
		t.Fatalf("expected evicted [1], got %v", evicted)
	}

	// GroupRemove
	c.GroupRemove(guild)
	if len(evicted) != 2 || evicted[1] != 2 {
		t.Fatalf("expected evicted [1, 2], got %v", evicted)
	}
}

func TestGroupedLRUCache_ReentrantCallback(t *testing.T) {
	// Verify that grouped onEvict callback can safely call cache operations without deadlocking.
	var c GroupedCache[testEntity]
	var reentrantCalled bool
	done := make(chan struct{})

	c = NewGroupedLRUCache[testEntity](FlagsAll, FlagsAll, nil, 2,
		WithGroupedLRUOnEvict[testEntity](func(groupID snowflake.ID, id snowflake.ID, entity testEntity) {
			if !reentrantCalled {
				reentrantCalled = true
				_, _ = c.Get(groupID, id)
				c.Put(groupID, 1000, testEntity{ID: 1000})
				_, _ = c.Remove(groupID, 1000)
			}
		}),
	)

	go func() {
		guild := snowflake.ID(1)
		c.Put(guild, 1, testEntity{ID: 1})
		c.Put(guild, 2, testEntity{ID: 2})
		// Trigger eviction of 1
		c.Put(guild, 3, testEntity{ID: 3})
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("deadlock detected: grouped onEvict callback failed reentrancy")
	}
}

func TestGroupedLRUCache_IteratorMutation(t *testing.T) {
	c := NewGroupedLRUCache[testEntity](FlagsAll, FlagsAll, nil, 10)
	guild := snowflake.ID(10)
	for i := 1; i <= 5; i++ {
		c.Put(guild, snowflake.ID(i), testEntity{ID: snowflake.ID(i)})
	}

	done := make(chan struct{})
	go func() {
		for entity := range c.GroupAll(guild) {
			c.Remove(guild, entity.ID)
		}
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("deadlock detected: mutating grouped cache during GroupAll() iteration")
	}

	if c.GroupLen(guild) != 0 {
		t.Fatalf("expected group len 0, got %d", c.GroupLen(guild))
	}
}

func TestGroupedLRUCache_GroupRemove(t *testing.T) {
	c := NewGroupedLRUCache[testEntity](FlagsAll, FlagsAll, nil, 10)

	guild := snowflake.ID(999)
	c.Put(guild, 1, testEntity{ID: 1})
	c.Put(guild, 2, testEntity{ID: 2})

	if c.GroupLen(guild) != 2 {
		t.Fatalf("expected len 2")
	}

	c.GroupRemove(guild)

	if c.GroupLen(guild) != 0 {
		t.Fatalf("expected len 0 after GroupRemove")
	}
	if _, ok := c.Get(guild, 1); ok {
		t.Fatalf("expected user 1 to be removed")
	}
}

func TestGroupedLRUCache_MaxGroups(t *testing.T) {
	var evictedGroups []snowflake.ID
	c := NewGroupedLRUCache[testEntity](FlagsAll, FlagsAll, nil, 10,
		WithGroupedLRUMaxGroups[testEntity](2),
		WithGroupedLRUOnEvict[testEntity](func(groupID snowflake.ID, id snowflake.ID, entity testEntity) {
			evictedGroups = append(evictedGroups, groupID)
		}),
	)

	c.Put(10, 1, testEntity{ID: 1})
	c.Put(20, 1, testEntity{ID: 1})

	// Access group 10 so it becomes MRU; order is 10, 20
	c.Get(10, 1)

	// Add group 30 -> group 20 should be evicted
	c.Put(30, 1, testEntity{ID: 1})

	if c.GroupLen(20) != 0 {
		t.Fatalf("expected group 20 to be evicted")
	}
	if c.GroupLen(10) != 1 {
		t.Fatalf("expected group 10 to remain")
	}
	if c.GroupLen(30) != 1 {
		t.Fatalf("expected group 30 to remain")
	}
}

func TestGroupedLRUCache_TTL(t *testing.T) {
	c := NewGroupedTTLCache[testEntity](FlagsAll, FlagsAll, nil, 50*time.Millisecond)

	guild := snowflake.ID(1)
	c.Put(guild, 1, testEntity{ID: 1, Name: "expiring"})

	if _, ok := c.Get(guild, 1); !ok {
		t.Fatalf("expected item to be present immediately")
	}

	time.Sleep(70 * time.Millisecond)

	if _, ok := c.Get(guild, 1); ok {
		t.Fatalf("expected item to expire after TTL")
	}
	if c.GroupLen(guild) != 0 {
		t.Fatalf("expected group len 0 after expiry")
	}
}

func TestGroupedLRUCache_Concurrency(t *testing.T) {
	c := NewGroupedLRUCache[testEntity](FlagsAll, FlagsAll, nil, 50, WithGroupedLRUTTL[testEntity](200*time.Millisecond))

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(2)
		val := i
		go func() {
			defer wg.Done()
			groupID := snowflake.ID(val % 5)
			id := snowflake.ID(val)
			c.Put(groupID, id, testEntity{ID: id, Name: "threaded"})
		}()
		go func() {
			defer wg.Done()
			groupID := snowflake.ID(val % 5)
			id := snowflake.ID(val)
			c.Get(groupID, id)
		}()
	}
	wg.Wait()
}
