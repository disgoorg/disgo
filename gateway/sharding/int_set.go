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

func (s *IntSet) Values() []int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	values := make([]int, len(s.set))
	i := 0
	for ii := range s.set {
		values[i] = ii
		i++
	}
	return values
}

func (s *IntSet) Add(i int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.set[i] = struct{}{}
}

func (s *IntSet) Has(i int) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.set[i]
	return ok
}

func (s *IntSet) Delete(i int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.set, i)
}

func (s *IntSet) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.set)
}

func (s *IntSet) String() string {
	var builder strings.Builder
	builder.WriteString("[")
	values := s.Values()
	for i := range values {
		builder.WriteString(strconv.Itoa(i))
		if i < len(values)-1 {
			builder.WriteString(", ")
		}
	}
	builder.WriteString("]")
	return builder.String()
}
