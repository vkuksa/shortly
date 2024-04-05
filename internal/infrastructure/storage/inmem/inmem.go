package inmem

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/vkuksa/shortly/internal/domain"
	"github.com/vkuksa/shortly/internal/link"
)

type Storage struct {
	mut sync.RWMutex
	m   map[string][]byte
}

func NewStorage() *Storage {
	return &Storage{m: make(map[string][]byte)}
}

func (r *Storage) GetLink(_ context.Context, uuid string) (*domain.Link, error) {
	r.mut.RLock()
	data, found := r.m[uuid]
	r.mut.RUnlock()
	if !found {
		return nil, link.ErrNotFound
	}

	var result domain.Link
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *Storage) StoreLink(_ context.Context, link *domain.Link) error {
	data, err := json.Marshal(link)
	if err != nil {
		return err
	}

	r.mut.Lock()
	defer r.mut.Unlock()
	r.m[link.UUID] = data
	return nil
}

func (r *Storage) IncHit(ctx context.Context, uuid string) error {
	link, err := r.GetLink(ctx, uuid)
	if err != nil {
		return fmt.Errorf("get link: %w", err)
	}

	link.Count++

	err = r.StoreLink(ctx, link)
	if err != nil {
		return fmt.Errorf("store link: %w", err)
	}

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
