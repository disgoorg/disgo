package set

import "sync"

type Set[T comparable] struct {
	mu    sync.RWMutex
	items map[T]struct{}
}

func New[T comparable]() *Set[T] {
	return &Set[T]{
		items: make(map[T]struct{}),
	}
}

func (s *Set[T]) Add(item T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[item] = struct{}{}
}

func (s *Set[T]) Remove(item T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.items, item)
}

func (s *Set[T]) Has(item T) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.items[item]
	return ok
}

func (s *Set[T]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.items)
}

func (s *Set[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = make(map[T]struct{})
}

func (s *Set[T]) ForEach(f func(item T)) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for item := range s.items {
		f(item)
	}
}
