package oauth2

import (
	"sync"
	"time"
)

type value struct {
	value      string
	insertedAt int64
}

func NewTTLMap(maxTTL time.Duration) *TTLMap {
	m := &TTLMap{
		maxTTL: maxTTL,
		m:      map[string]value{},
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

type TTLMap struct {
	maxTTL time.Duration
	m      map[string]value
	mu     sync.Mutex
}

func (m *TTLMap) Len() int {
	return len(m.m)
}

func (m *TTLMap) Put(k string, v string) {
	m.mu.Lock()
	m.m[k] = value{v, time.Now().Unix()}
	m.mu.Unlock()
}

func (m *TTLMap) Get(k string) string {
	m.mu.Lock()
	v, ok := m.m[k]
	m.mu.Unlock()
	if ok {
		return v.value
	}
	return ""
}

func (m *TTLMap) Delete(k string) {
	m.mu.Lock()
	delete(m.m, k)
	m.mu.Unlock()
}
