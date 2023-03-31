package set

import "sync"

type Set[T comparable] interface {
	Add(item T)
	Remove(item T)
	Has(item T) bool
	Len() int
	Clear()
	ForEach(f func(item T))
}
type set[T comparable] struct {
	mu    sync.RWMutex
	items map[T]struct{}
}

func New[T comparable]() Set[T] {
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
