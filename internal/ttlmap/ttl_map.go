package ttlmap

import (
	"sync"
	"time"
)

type value[V any] struct {
	value      V
	insertedAt int64
}

func New[K comparable, V any](maxTTL time.Duration) *Map[K, V] {
	m := &Map[K, V]{
		maxTTL: maxTTL,
		m:      map[K]value[V]{},
	}

	if maxTTL > 0 {
		go func() {
			ticker := time.NewTicker(10 * time.Second)
			defer ticker.Stop()
			for now := range ticker.C {
				m.mu.Lock()
				for k, v := range m.m {
					if now.Unix()-v.insertedAt > int64(m.maxTTL) {
						delete(m.m, k)
					}
				}
				m.mu.Unlock()
			}
		}()
	}

	return m
}

type Map[K comparable, V any] struct {
	maxTTL time.Duration
	m      map[K]value[V]
	mu     sync.Mutex
}

func (m *Map[K, V]) Len() int {
	return len(m.m)
}

func (m *Map[K, V]) Put(k K, v V) {
	m.mu.Lock()
	m.m[k] = value[V]{v, time.Now().Unix()}
	m.mu.Unlock()
}

func (m *Map[K, V]) Get(k K) (V, bool) {
	m.mu.Lock()
	v, ok := m.m[k]
	m.mu.Unlock()
	if ok {
		return v.value, true
	}
	var empty V
	return empty, false
}

func (m *Map[K, V]) Delete(k K) {
	m.mu.Lock()
	delete(m.m, k)
	m.mu.Unlock()
}
