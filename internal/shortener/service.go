//nolint

package shortener

import (
	"context"
	"encoding/base64"
	"time"

	shortly "github.com/vkuksa/shortly/internal/domain"
)

type Service struct {
	stor shortly.LinkStorage
}

func NewService(s shortly.LinkStorage) *Service {
	return &Service{stor: s}
}

// nolint: revive
func (s *Service) GenerateShortenedLink(ctx context.Context, url string) (*shortly.Link, error) {
	if url == "" {
		return nil, shortly.NewError(shortly.ErrInvalid, "No url value provided for shortening")
	}

	link := &shortly.Link{}
	uuid := base64.URLEncoding.EncodeToString([]byte(url))
	expires := time.Now().AddDate(0, 1, 0)

	// Get link from a storage
	if found, err := s.stor.Get(uuid, link); err != nil {
		// If there's an error occured with retrieving a link
		return nil, shortly.NewError(shortly.ErrInternal, "GenerateShortenedLink: %s", err.Error())
	} else if !found {
		// If there's no link in a storage - fill the fields
		link.UUID = uuid
		link.URL = url
		link.CreatedAt = time.Now()
	}
	// If there's a link in a storage and there were no errors in retrieving it: it will be in a link variable

	// Update expiration date of a link
	link.ExpiresAt = expires

	// Store new link object
	if err := s.stor.Set(uuid, *link); err != nil {
		return nil, shortly.NewError(shortly.ErrInternal, "GenerateShortenedLink: %s", err.Error())
	}

	return link, nil
}

// nolint: revive
func (s *Service) GetOriginalLink(ctx context.Context, uuid string) (*shortly.Link, error) {
	link := &shortly.Link{}

	// Get link from a storage
	if found, err := s.stor.Get(uuid, link); err != nil {
		return nil, shortly.NewError(shortly.ErrInternal, "GetOriginalLink: %s", err.Error())
	} else if !found {
		return nil, shortly.NewError(shortly.ErrNotFound, "No link stored for a given uuid")
	}

	return link, nil
}

// nolint: revive
func (s *Service) AddHit(ctx context.Context, uuid string) error {
	link := &shortly.Link{}

	// Get link from a storage
	if found, err := s.stor.Get(uuid, link); err != nil {
		return shortly.NewError(shortly.ErrInternal, "AddHit: %s", err.Error())
	} else if !found {
		return shortly.NewError(shortly.ErrNotFound, "No link stored for a given uuid")
	}

	// Update count
	link.Count++

	// Update stored link
	if err := s.stor.Set(uuid, *link); err != nil {
		return shortly.NewError(shortly.ErrInternal, "AddHit: %s", err.Error())
	}

	return nil
}
