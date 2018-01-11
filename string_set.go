package main

// StringSet is an implementation of a generic set of strings.
// Map is chosed as the baseline structure for its retrieval complexity of O(1)
type StringSet struct {
	set map[string]bool
}

// NewStringSet returns an empty StringSet
func NewStringSet() StringSet {
	return StringSet{map[string]bool{}}
}

// Add value to a StringSet
// Returns a boolean value that indicates whether the item was added.
func (s StringSet) Add(value string) bool {
	_, found := s.set[value]
	s.set[value] = true
	return !found
}

// List all values stored in StringSet
func (s StringSet) List() []string {
	list := make([]string, 0, len(s.set))
	for k := range s.set {
		list = append(list, k)
	}
	return list
}

// Length returns the length of the StringSet
func (s StringSet) Length() int {
	return len(s.set)
}
