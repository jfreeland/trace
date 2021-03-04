package storage

import (
	"sync"
)

// InMemory is an in memory storage.
type InMemory struct {
	mu   sync.RWMutex
	data map[string]interface{}
}

// NewInMemory creates a new in memory storage.
func NewInMemory() *InMemory {
	data := make(map[string]interface{})

	return &InMemory{
		data: data,
	}
}

// Store a payload.
func (s *InMemory) Store(data *map[string]interface{}) {
	for k, v := range *data {
		s.mu.Lock()
		s.data[k] = v
		s.mu.Unlock()
	}
}

// Get a single key.
func (s *InMemory) Get(key string) (interface{}, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]

	return val, ok
}

// GetAll returns all values in storage.
func (s *InMemory) GetAll() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.data
}
