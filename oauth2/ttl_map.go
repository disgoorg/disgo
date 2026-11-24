package oauth2

import (
	"sync"
	"time"
)

type value struct {
	value      string
	insertedAt int64
}

func newTTLMap(maxTTL time.Duration) *ttlMap {
	m := &ttlMap{
		maxTTL: maxTTL,
		m:      map[string]value{},
		stop:   make(chan struct{}),
	}

	if maxTTL > 0 {
		go func() {
			ticker := time.NewTicker(10 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-m.stop:
					return
				case <-ticker.C:
					m.mu.Lock()
					for k, v := range m.m {
						if time.Since(time.Unix(v.insertedAt, 0)) > m.maxTTL {
							delete(m.m, k)
						}
					}
					m.mu.Unlock()
				}
			}
		}()
	}

	return m
}

type ttlMap struct {
	maxTTL   time.Duration
	m        map[string]value
	mu       sync.Mutex
	stop     chan struct{}
	stopOnce sync.Once
}

func (m *ttlMap) close() {
	m.stopOnce.Do(func() {
		close(m.stop)
	})
}

func (m *ttlMap) put(k string, v string) {
	m.mu.Lock()
	m.m[k] = value{v, time.Now().Unix()}
	m.mu.Unlock()
}

func (m *ttlMap) get(k string) string {
	m.mu.Lock()
	v, ok := m.m[k]
	m.mu.Unlock()
	if ok {
		return v.value
	}
	return ""
}

func (m *ttlMap) delete(k string) {
	m.mu.Lock()
	delete(m.m, k)
	m.mu.Unlock()
}
