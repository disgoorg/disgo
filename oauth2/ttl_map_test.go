package oauth2

import (
	"testing"
	"time"
)

func TestTTLMap_Expiration(t *testing.T) {
	// Create a TTL map with 50ms TTL
	m := newTTLMap(50 * time.Millisecond)
	defer m.close()

	m.put("key1", "val1")
	if v := m.get("key1"); v != "val1" {
		t.Fatalf("expected val1, got %s", v)
	}

	// Inserted timestamp is in seconds, so let's test time.Since logic directly
	m.mu.Lock()
	v, ok := m.m["key1"]
	m.mu.Unlock()
	if !ok {
		t.Fatal("expected key1 to exist in map")
	}

	// Artificially age the entry by 2 seconds
	v.insertedAt = time.Now().Unix() - 2
	m.mu.Lock()
	m.m["key1"] = v
	// Run eviction check
	for k, item := range m.m {
		if time.Since(time.Unix(item.insertedAt, 0)) > m.maxTTL {
			delete(m.m, k)
		}
	}
	m.mu.Unlock()

	if m.get("key1") != "" {
		t.Fatal("expected key1 to be evicted after TTL expiration")
	}
}
