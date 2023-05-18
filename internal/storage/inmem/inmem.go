package inmem

import (
	"context"
	"fmt"
	"sync"

	"github.com/vkuksa/shortly"
)

// Ensure service implements interface.
var _ shortly.LinkStorage = (*LinkStorage)(nil)

type LinkStorage struct {
	mu   sync.RWMutex
	stor map[string]*shortly.Link
}

func NewLinkStorage() *LinkStorage {
	return &LinkStorage{stor: make(map[string]*shortly.Link)}
}

// Saves given link into in-memory storage
// Returns an error, if given link value is nil or there's aalready a link with a specified uuid
func (s *LinkStorage) SaveLink(ctx context.Context, link *shortly.Link) error {
	if link == nil {
		return shortly.Errorf(shortly.ERRINVALID, "SaveLink: empty link value provided for save")
	}

	if _, err := s.findLink(link.UUID); err == nil {
		// There's a link found in datastore
		return shortly.Errorf(shortly.ERRCONFLICT, fmt.Sprintf("SaveLink: there's already a link stored for a given key %s", link.UUID))
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.stor[link.UUID] = link
	return nil
}

// Find a link by a uuid value from in-memory storage
// Returns an error, if no links have been found for a given uuid
func (s *LinkStorage) FindLink(ctx context.Context, uuid string) (*shortly.Link, error) {
	link, err := s.findLink(uuid)
	if err != nil {
		// No link found for provided uuid
		return nil, fmt.Errorf("FindLink: %w", err)
	}

	return link, nil
}

// Delete's a link from in-memory storage by provided uuid value
// If there's no value with a given uuid in collection - operation returns an error
func (s *LinkStorage) DeleteLink(ctx context.Context, uuid string) error {
	_, err := s.findLink(uuid)
	if err != nil {
		// No link found for provided uuid
		return fmt.Errorf("DeleteLink: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.stor, uuid)
	return nil
}

func (s *LinkStorage) UpdateLink(ctx context.Context, uuid string, upd shortly.LinkUpdate) (*shortly.Link, error) {
	link, err := s.findLink(uuid)
	if err != nil {
		// No link found for provided uuid
		return nil, fmt.Errorf("UpdateLink: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if upd.Count != nil {
		link.Count = *upd.Count
	}
	if upd.CreatedAt != nil {
		link.CreatedAt = *upd.CreatedAt
	}
	if upd.ExpiresAt != nil {
		link.ExpiresAt = *upd.ExpiresAt
	}
	return link, nil
}

func (s *LinkStorage) findLink(uuid string) (*shortly.Link, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	link, found := s.stor[uuid]
	if !found {
		return nil, shortly.Errorf(shortly.ERRNOTFOUND, fmt.Sprintf("findLink: no link with uuid %s found", uuid))
	}
	return link, nil
}

// Nop
func (s *LinkStorage) Open(ctx context.Context) error {
	return nil
}

// Nop
func (s *LinkStorage) Close() error {
	return nil
}
