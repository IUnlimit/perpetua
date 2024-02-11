package concurrent

import "sync"

type Set struct {
	mu    sync.RWMutex
	items map[interface{}]bool
}

func NewSet() *Set {
	return &Set{
		items: make(map[interface{}]bool),
	}
}

func (s *Set) Add(item interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[item] = true
}

func (s *Set) Remove(item interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.items, item)
}

func (s *Set) Contains(item interface{}) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.items[item]
}

func (s *Set) Iterator() []interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make([]interface{}, 0, len(s.items))
	for item := range s.items {
		result = append(result, item)
	}

	return result
}
