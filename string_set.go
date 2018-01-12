package main

import "sync"

// StringSet is a thread-safe implementation of a generic set of strings.
// Map is chosen as the underlying structure for its retrieval complexity of O(1).
type StringSet struct {
	sync.Mutex
	set map[string]bool
}

// NewStringSet returns an empty StringSet.
func NewStringSet() *StringSet {
	return &StringSet{set: map[string]bool{}}
}

// Add value to a StringSet.
// Returns a boolean value that indicates whether the item was added.
func (s *StringSet) Add(value string) bool {
	s.Lock()
	defer s.Unlock()
	_, found := s.set[value]
	s.set[value] = true
	return !found
}

// List all values stored in StringSet.
func (s *StringSet) List() []string {
	s.Lock()
	defer s.Unlock()
	list := make([]string, 0, len(s.set))
	for k := range s.set {
		list = append(list, k)
	}
	return list
}

// Length function returns the length of a StringSet.
func (s *StringSet) Length() int {
	s.Lock()
	defer s.Unlock()
	return len(s.set)
}
