package cache

import (
	"iter"
)

// Set is a collection of unique items. It should be thread-safe.
type Set[T comparable] interface {
	// Add adds an item to the set.
	Add(item T, opts ...AccessOpt) error
	// Remove removes an item from the set.
	Remove(item T, opts ...AccessOpt) error
	// Has returns true if the item is in the set, false otherwise.
	Has(item T, opts ...AccessOpt) (bool, error)
	// Len returns the number of items in the set.
	Len(opts ...AccessOpt) (int, error)
	// Clear removes all items from the set.
	Clear(opts ...AccessOpt) error
	// All returns an iterator over all items in the set.
	All(opts ...AccessOpt) (iter.Seq2[T, error], error)
}

type set[T comparable] struct {
	mu    RWMutex
	items map[T]struct{}
}

// NewSet returns a thread-safe in memory implementation of a set.
func NewSet[T comparable]() Set[T] {
	return &set[T]{
		items: make(map[T]struct{}),
	}
}

func (s *set[T]) Add(item T, opts ...AccessOpt) error {
	cfg := resolveAccessConfig(opts)

	if err := s.mu.Lock(cfg.Ctx); err != nil {
		return err
	}

	defer s.mu.Unlock()
	s.items[item] = struct{}{}
	return nil
}

func (s *set[T]) Remove(item T, opts ...AccessOpt) error {
	cfg := resolveAccessConfig(opts)

	if err := s.mu.Lock(cfg.Ctx); err != nil {
		return err
	}

	defer s.mu.Unlock()
	delete(s.items, item)
	return nil
}

func (s *set[T]) Has(item T, opts ...AccessOpt) (bool, error) {
	cfg := resolveAccessConfig(opts)

	if err := s.mu.RLock(cfg.Ctx); err != nil {
		return false, err
	}

	defer s.mu.RUnlock()
	_, ok := s.items[item]
	return ok, nil
}

func (s *set[T]) Len(opts ...AccessOpt) (int, error) {
	cfg := resolveAccessConfig(opts)

	if err := s.mu.RLock(cfg.Ctx); err != nil {
		return 0, err
	}

	defer s.mu.RUnlock()
	return len(s.items), nil
}

func (s *set[T]) Clear(opts ...AccessOpt) error {
	cfg := resolveAccessConfig(opts)

	if err := s.mu.Lock(cfg.Ctx); err != nil {
		return err
	}

	defer s.mu.Unlock()
	s.items = make(map[T]struct{})
	return nil
}

func (s *set[T]) All(opts ...AccessOpt) (iter.Seq2[T, error], error) {
	cfg := resolveAccessConfig(opts)

	if err := s.mu.RLock(cfg.Ctx); err != nil {
		return nil, err
	}

	return func(yield func(T, error) bool) {
		defer s.mu.RUnlock()

		for item := range s.items {
			if !yield(item, nil) {
				return
			}
		}
	}, nil
}
