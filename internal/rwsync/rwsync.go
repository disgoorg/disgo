package rwsync

import "sync"

// RWLocker is a read-write mutex interface
type RWLocker interface {
	sync.Locker
	RLock()
	RUnlock()
}
