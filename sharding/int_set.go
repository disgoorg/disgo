package sharding

import (
	"strconv"
	"strings"
	"sync"
)

func NewIntSet(ints ...int) *IntSet {
	set := &IntSet{
		Set: make(map[int]struct{}, len(ints)),
	}
	for _, i := range ints {
		set.Set[i] = struct{}{}
	}
	return set
}

type IntSet struct {
	sync.RWMutex
	Set map[int]struct{}
}

func (s *IntSet) Add(i int) {
	s.Lock()
	s.Set[i] = struct{}{}
	s.Unlock()
}

func (s *IntSet) Has(i int) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.Set[i]
	return ok
}

func (s *IntSet) Delete(i int) {
	s.RLock()
	_, ok := s.Set[i]
	s.RUnlock()
	if ok {
		s.Lock()
		delete(s.Set, i)
		s.Unlock()
	}
}

func (s *IntSet) Len() int {
	return len(s.Set)
}

func (s *IntSet) String() string {
	var builder strings.Builder
	builder.WriteString("[")
	s.RLock()
	for i := range s.Set {
		builder.WriteString(strconv.Itoa(i))
		if i < len(s.Set)-1 {
			builder.WriteString(", ")
		}
	}
	s.RUnlock()
	builder.WriteString("]")
	return builder.String()
}
