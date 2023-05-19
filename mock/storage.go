package mock

import (
	"errors"

	shortly "github.com/vkuksa/shortly/internal/domain"
)

// Ensure service implements interface.
var _ shortly.LinkStorage = (*LinkStorage)(nil)

type LinkStorage struct {
	Stor map[string]shortly.Link
}

func NewLinkStorage() *LinkStorage {
	return &LinkStorage{}
}

func (s *LinkStorage) Set(key string, link shortly.Link) error {
	if key == "" {
		return errors.New("non-valid key")
	}

	s.Stor[link.UUID] = link
	return nil
}

func (s *LinkStorage) Get(key string, link *shortly.Link) (bool, error) {
	if key == "" {
		return false, errors.New("non-valid key")
	}

	data, found := s.Stor[key]
	if !found {
		return false, nil
	}

	*link = data

	return true, nil
}

func (s *LinkStorage) Delete(key string) error {
	if key == "" {
		return errors.New("non-valid key")
	}

	delete(s.Stor, key)
	return nil
}
