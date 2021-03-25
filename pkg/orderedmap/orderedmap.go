package orderedmap

import (
	"fmt"
	"sync"
)

// OrderedMap represents an associative array or map abstract data type.
type OrderedMap struct {
	// mu Mutex protects data structures below.
	mu sync.Mutex

	// keys is the Set list of keys.
	keys []interface{}

	// store is the Set underlying store of values.
	store map[interface{}]interface{}
}

// NewOrderedMap creates a new empty OrderedMap.
func NewOrderedMap() *OrderedMap {
	m := &OrderedMap{
		keys:  make([]interface{}, 0),
		store: make(map[interface{}]interface{}),
	}

	return m
}

// Put adds items to the map.
//
// If a key is found in the map it replaces it value.
func (m *OrderedMap) Put(key interface{}, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.store[key]; !ok {
		m.keys = append(m.keys, key)
	}

	m.store[key] = value
}

// Get returns the value of a key from the OrderedMap.
func (m *OrderedMap) Get(key interface{}) (value interface{}, found bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	value, found = m.store[key]
	return value, found
}

// Remove deletes a key-value pair from the OrderedMap.
//
// If a key is not found in the map it doesn't fail, just does nothing.
func (m *OrderedMap) Remove(key interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check key exists
	if _, found := m.store[key]; !found {
		return
	}

	// Remove the value from the store
	delete(m.store, key)

	// Remove the key
	for i := range m.keys {
		if m.keys[i] == key {
			m.keys = append(m.keys[:i], m.keys[i+1:]...)
			break
		}
	}
}

// Size returns the map's number of key-value pairs.
func (m *OrderedMap) Size() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return len(m.store)
}

// Empty return if the map in empty or not.
func (m *OrderedMap) Empty() bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	return len(m.store) == 0
}

// Keys return the keys in the map in insertion order.
func (m *OrderedMap) Keys() []interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.keys
}

// Values return the values in the map in insertion order.
func (m *OrderedMap) Values() []interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()

	values := make([]interface{}, len(m.store))
	for i, key := range m.keys {
		values[i] = m.store[key]
	}
	return values
}

// Contains returns a bool indicating whether or not the map contains a
// given value
func (m *OrderedMap) Contains(v interface{}) bool {
	for _, key := range m.keys {
		if m.store[key] == v {
			return true
		}
	}
	return false
}

// String implements Stringer interface.
//
// Prints the map string representation, a concatenated string of all its
// string representation values in insertion order.
func (m *OrderedMap) String() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	var result []string
	for i, key := range m.keys {
		result = append(result, fmt.Sprintf("%d:%s", m.keys[i].(int), m.store[key]))
	}

	return fmt.Sprintf("%s", result)
}
