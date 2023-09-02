package cache

import "sync"

// Set is a collection of unique items. It should be thread-safe.
type Set[T comparable] interface {
	// Add adds an item to the set.
	Add(item T)
	// Remove removes an item from the set.
	Remove(item T)
	// Has returns true if the item is in the set, false otherwise.
	Has(item T) bool
	// Len returns the number of items in the set.
	Len() int
	// Clear removes all items from the set.
	Clear()
	// ForEach iterates over the items in the set.
	ForEach(f func(item T))
}

type set[T comparable] struct {
	mu    sync.RWMutex
	items map[T]struct{}
}

// NewSet returns a thread-safe in memory implementation of a set.
func NewSet[T comparable]() Set[T] {
	return &set[T]{
		items: make(map[T]struct{}),
	}
}

func (s *set[T]) Add(item T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[item] = struct{}{}
}

func (s *set[T]) Remove(item T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.items, item)
}

func (s *set[T]) Has(item T) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.items[item]
	return ok
}

func (s *set[T]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.items)
}

func (s *set[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = make(map[T]struct{})
}

func (s *set[T]) ForEach(f func(item T)) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for item := range s.items {
		f(item)
	}
}
