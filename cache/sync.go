package cache

import (
	"context"

	"github.com/sasha-s/go-csync"
)

// RWMutex is a read-write mutex that can be used to synchronize access to a shared resource.
// Right now it is just a wrapper around the csync.Mutex, but we can make it a proper read-write mutex in the future.
type RWMutex struct {
	mu csync.Mutex
}

func (m *RWMutex) RLock(ctx context.Context) error {
	return m.mu.CLock(ctx)
}

func (m *RWMutex) RUnlock() {
	m.mu.Unlock()
}

func (m *RWMutex) Lock(ctx context.Context) error {
	return m.mu.CLock(ctx)
}

func (m *RWMutex) Unlock() {
	m.mu.Unlock()
}
