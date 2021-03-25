package orderedset

import (
	"fmt"
	"github.com/openmesh/flow/pkg/orderedmap"
	"sync"
)

// OrderedSet represents a dynamic, insertion-ordered, set abstract data type.
type OrderedSet struct {
	// currentIndex keeps track of the keys of the underlying store.
	currentIndex int

	// mu Mutex protects data structures below.
	mu sync.Mutex

	// index is the Set list of keys.
	index map[interface{}]int

	// store is the Set underlying store of values.
	store *orderedmap.OrderedMap
}

// NewOrderedSet creates a new empty OrderedSet.
func NewOrderedSet() *OrderedSet {
	orderedSet := &OrderedSet{
		index: make(map[interface{}]int),
		store: orderedmap.NewOrderedMap(),
	}

	return orderedSet
}

// Add adds items to the set.
//
// If an item is found in the set it replaces it.
func (s *OrderedSet) Add(items ...interface{}) {
	for _, item := range items {
		if _, found := s.index[item]; found {
			continue
		}

		s.put(item)
	}
}

// Remove deletes items from the set.
//
// If an item is not found in the set it doesn't fails, just does nothing.
func (s *OrderedSet) Remove(items ...interface{}) {
	for _, item := range items {
		index, found := s.index[item]
		if !found {
			return
		}

		s.remove(index, item)
	}
}

// Contains return if set contains the specified items or not.
func (s *OrderedSet) Contains(items ...interface{}) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, item := range items {
		if _, found := s.index[item]; !found {
			return false
		}
	}
	return true
}

// Empty return if the set in empty or not.
func (s *OrderedSet) Empty() bool {
	return s.store.Empty()
}

// Values return the set values in insertion order.
func (s *OrderedSet) Values() []interface{} {
	return s.store.Values()
}

// Size return the set number of elements.
func (s *OrderedSet) Size() int {
	return s.store.Size()
}

// String implements Stringer interface.
//
// Prints the set string representation, a concatenated string of all its
// string representation values in insertion order.
func (s *OrderedSet) String() string {
	return fmt.Sprintf("%s", s.Values())
}

// Put adds a single item into the set
func (s *OrderedSet) put(item interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.store.Put(s.currentIndex, item)
	s.index[item] = s.currentIndex
	s.currentIndex++
}

// remove deletes a single item from the test given its index
func (s *OrderedSet) remove(index int, item interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.store.Remove(index)
	delete(s.index, item)
}
