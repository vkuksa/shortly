package inmem

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/dolthub/swiss"

	"github.com/vkuksa/shortly/internal/link"
)

var _ link.Repository = (*Storage)(nil)

// Storage is an in-memory implementation of the link.Repository interface.
type Storage struct {
	mut           sync.RWMutex
	links         *swiss.Map[string, *link.ShortenedLink] // UUID → ShortenedLink
	shortenedToID *swiss.Map[string, string]              // shortened → UUID
}

// NewStorage constructs a new in-memory storage.
func NewStorage() *Storage {
	return &Storage{
		links:         swiss.NewMap[string, *link.ShortenedLink](0),
		shortenedToID: swiss.NewMap[string, string](0),
	}
}

// Get retrieves a ShortenedLink by its short string.
func (s *Storage) Get(ctx context.Context, shortened string) (*link.ShortenedLink, error) {
	s.mut.RLock()
	defer s.mut.RUnlock()

	id, ok := s.shortenedToID.Get(shortened)
	if !ok {
		return nil, link.ErrNotFound
	}
	lk, ok := s.links.Get(id)
	if !ok {
		return nil, link.ErrNotFound
	}
	return lk, nil
}

// Store persists a ShortenedLink in memory, updating both indexes.
func (s *Storage) Store(ctx context.Context, l *link.ShortenedLink) error {
	s.mut.Lock()
	defer s.mut.Unlock()

	s.shortenedToID.Put(l.Shortened, l.Id)
	s.links.Put(l.Id, l)

	return nil
}

// AddHit increments the Hits counter of the link identified by its UUID.
func (s *Storage) AddHit(ctx context.Context, id string) error {
	s.mut.Lock()
	defer s.mut.Unlock()

	lk, ok := s.links.Get(id)
	if !ok {
		return fmt.Errorf("addHit: %w", link.ErrNotFound)
	}

	lk.Hits++
	// Since lk is already a pointer in the map, no need to re-store it
	return nil
}

// UpdateExpiration sets a new expiration date/time for the link identified by UUID.
func (s *Storage) UpdateExpiration(ctx context.Context, id string, expiresAt time.Time) error {
	s.mut.Lock()
	defer s.mut.Unlock()

	lk, ok := s.links.Get(id)
	if !ok {
		return fmt.Errorf("updateExpiration: %w", link.ErrNotFound)
	}

	lk.ExpiresAt = expiresAt
	// No need to re-store since pointer is updated in place
	return nil
}
