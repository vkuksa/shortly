package mock

import (
	"context"
	"fmt"

	"github.com/vkuksa/shortly"
)

// Ensure service implements interface.
var _ shortly.LinkStorage = (*LinkStorage)(nil)

type LinkStorage struct {
	Stor map[string]*shortly.Link
}

func (s *LinkStorage) SaveLink(ctx context.Context, link *shortly.Link) error {
	if link.URL == "" {
		return fmt.Errorf("Empty url in a link given")
	}

	if _, err := s.findLink(link.UUID); err == nil {
		// There's a link found in datastore
		return shortly.Errorf(shortly.ERRCONFLICT, fmt.Sprintf("SaveLink: there's already a link stored for a given key %s", link.UUID))
	}

	s.Stor[link.UUID] = link
	return nil
}

func (s *LinkStorage) FindLink(ctx context.Context, uuid string) (*shortly.Link, error) {
	link, err := s.findLink(uuid)
	if err != nil {
		// No link found for provided uuid
		return nil, fmt.Errorf("FindLink: %w", err)
	}

	return link, nil
}

func (s *LinkStorage) DeleteLink(ctx context.Context, uuid string) error {
	_, err := s.findLink(uuid)
	if err != nil {
		// No link found for provided uuid
		return fmt.Errorf("DeleteLink: %w", err)
	}

	delete(s.Stor, uuid)
	return nil
}

func (s *LinkStorage) UpdateLink(ctx context.Context, uuid string, upd shortly.LinkUpdate) (*shortly.Link, error) {
	link, err := s.findLink(uuid)
	if err != nil {
		// No link found for provided uuid
		return nil, fmt.Errorf("UpdateLink: %w", err)
	}

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
	link, found := s.Stor[uuid]
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
