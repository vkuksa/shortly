package inmem

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/vkuksa/shortly/internal/domain"
	"github.com/vkuksa/shortly/internal/interface/repository"
)

type Storage struct {
	mut sync.RWMutex
	m   map[string][]byte
}

func NewStorage() *Storage {
	return &Storage{m: make(map[string][]byte)}
}

func (r *Storage) GetLink(_ context.Context, uuid string) (*domain.Link, error) {
	if err := repository.ValidateKey(uuid); err != nil {
		return nil, fmt.Errorf("GetLink: %w", err)
	}

	r.mut.RLock()
	data, found := r.m[uuid]
	r.mut.RUnlock()
	if !found {
		return nil, repository.ErrValueNotFound
	}

	var result domain.Link
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("GetLink: %w", err)
	}

	return &result, nil
}

func (r *Storage) StoreLink(_ context.Context, link *domain.Link) error {
	if err := repository.ValidateKey(link.UUID); err != nil {
		return fmt.Errorf("StoreLink: %w", err)
	}

	data, err := json.Marshal(link)
	if err != nil {
		return fmt.Errorf("set: %w", err)
	}

	r.mut.Lock()
	defer r.mut.Unlock()
	r.m[link.UUID] = data
	return nil
}

// func (s *Storage[V]) Delete(k string) error {
// 	if err := storage.ValidateKey(k); err != nil {
// 		return fmt.Errorf("delete: %w", err)
// 	}

// 	s.mut.Lock()
// 	defer s.mut.Unlock()
// 	delete(s.m, k)
// 	return nil
// }
