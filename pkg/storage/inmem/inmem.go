package inmem

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/vkuksa/shortly/pkg/storage"
)

type Storage[V any] struct {
	mut sync.RWMutex
	m   map[string][]byte
}

func NewStorage[V any]() *Storage[V] {
	return &Storage[V]{m: make(map[string][]byte)}
}

// Saves data into in-memory storage
// Returns an error, if given key or value are not valid
func (s *Storage[V]) Set(key string, value V) error {
	if s.m == nil {
		panic("Trying to set a value in a nil map")
	}

	if err := storage.ValidateKey(key); err != nil {
		return fmt.Errorf("set: %w", err)
	}

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("set: %w", err)
	}

	s.mut.Lock()
	defer s.mut.Unlock()
	s.m[key] = data
	return nil
}

// Retrieves a stored value by a given key
// Returns an error, if no value have been found for a given key
func (s *Storage[V]) Get(k string, v *V) (bool, error) {
	if s.m == nil {
		panic("Trying to get a value from a nil map")
	}

	if err := storage.ValidateKey(k); err != nil {
		return false, fmt.Errorf("get: %w", err)
	}

	s.mut.RLock()
	data, found := s.m[k]
	s.mut.RUnlock()
	if !found {
		return false, nil
	}

	return true, json.Unmarshal(data, v)
}

// Deletes a key-value pair from a storage
// If there's no key stored, delete is no-op
// Returns an error if given key is not valid
func (s *Storage[V]) Delete(k string) error {
	if s.m == nil {
		panic("Trying to delete a value in a nil map")
	}

	if err := storage.ValidateKey(k); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	s.mut.Lock()
	defer s.mut.Unlock()
	delete(s.m, k)
	return nil
}

func (s *Storage[V]) Close() error {
	s.mut.Lock()
	defer s.mut.Unlock()
	s.m = nil
	return nil
}
