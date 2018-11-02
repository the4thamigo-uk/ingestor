package bimap

import (
	"sync"
)

// BiMap defines an interface for a bi-directional map of strings
type BiMap interface {
	Put(key string, val string)
	Get(key string) (string, bool)
	GetByVal(val string) (string, bool)
	Delete(key string) (string, bool)
	DeleteByVal(val string) (string, bool)
	Iterate(f IterateFunc)
}

type bimap struct {
	keyToVal map[string]string
	valToKey map[string]string
	mtx      sync.RWMutex
}

// IterateFunc is a callback interface used in the call to Iterate()
type IterateFunc func(key string, val string)

// New creates a new one-to-one bidirectional map of key <-> value strings
// It is safe for concurrent use by multiple goroutines
func New() BiMap {
	return &bimap{
		keyToVal: map[string]string{},
		valToKey: map[string]string{},
	}
}

// Put inserts an (id, val) pair into the map replacing any items with the same id or value)
func (m *bimap) Put(key string, val string) {
	defer m.mtx.Unlock()
	m.mtx.Lock()

	m.valToKey[val] = key
	m.keyToVal[key] = val
}

// Get retrieves the value for the given key
func (m *bimap) Get(key string) (string, bool) {
	defer m.mtx.RUnlock()
	m.mtx.RLock()

	val, ok := m.keyToVal[key]
	return val, ok
}

// Get retrieves the key for the given value
func (m *bimap) GetByVal(val string) (string, bool) {
	defer m.mtx.RUnlock()
	m.mtx.RLock()

	key, ok := m.valToKey[val]
	return key, ok
}

// Delete removes the value for the given key
func (m *bimap) Delete(key string) (string, bool) {
	defer m.mtx.RUnlock()
	m.mtx.RLock()

	val, ok := m.keyToVal[key]
	if ok {
		delete(m.keyToVal, key)
		delete(m.valToKey, val)
	}
	return val, ok
}

// Get retrieves the key for the given value
func (m *bimap) DeleteByVal(val string) (string, bool) {
	defer m.mtx.RUnlock()
	m.mtx.RLock()

	key, ok := m.valToKey[val]
	if ok {
		delete(m.keyToVal, key)
		delete(m.valToKey, val)
	}
	return key, ok
}

// Iterate calls the given function for each item in the map
// Calling another method of 'm' in the implementation of 'f' will cause a deadlock, so this is not supported.
func (m *bimap) Iterate(f IterateFunc) {
	defer m.mtx.RUnlock()
	m.mtx.RLock()

	for k, v := range m.keyToVal {
		f(k, v)
	}
}
