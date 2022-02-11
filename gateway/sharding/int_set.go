package sharding

import (
	"strconv"
	"strings"
	"sync"
)

func NewIntSet(ints ...int) *IntSet {
	set := &IntSet{
		set: make(map[int]struct{}, len(ints)),
	}
	for _, i := range ints {
		set.set[i] = struct{}{}
	}
	return set
}

type IntSet struct {
	mu  sync.RWMutex
	set map[int]struct{}
}

func (s *IntSet) Add(i int) {
	s.mu.Lock()
	s.set[i] = struct{}{}
	s.mu.Unlock()
}

func (s *IntSet) Has(i int) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.set[i]
	return ok
}

func (s *IntSet) Delete(i int) {
	s.mu.RLock()
	_, ok := s.set[i]
	s.mu.RUnlock()
	if ok {
		s.mu.Lock()
		delete(s.set, i)
		s.mu.Unlock()
	}
}

func (s *IntSet) Len() int {
	return len(s.set)
}

func (s *IntSet) String() string {
	var builder strings.Builder
	builder.WriteString("[")
	s.mu.RLock()
	for i := range s.set {
		builder.WriteString(strconv.Itoa(i))
		if i < len(s.set)-1 {
			builder.WriteString(", ")
		}
	}
	s.mu.RUnlock()
	builder.WriteString("]")
	return builder.String()
}
