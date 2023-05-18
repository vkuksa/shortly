package business

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/vkuksa/shortly"
)

var _ shortly.LinkService = (*LinkService)(nil)

type LinkService struct {
	storage shortly.LinkStorage
}

func NewLinkService(s shortly.LinkStorage) *LinkService {
	return &LinkService{storage: s}
}

func (s *LinkService) GenerateShortenedLink(ctx context.Context, url string) (*shortly.Link, error) {
	if url == "" {
		return nil, shortly.Errorf(shortly.ERRINVALID, "No url value provided for shortening")
	}

	uuid := base64.URLEncoding.EncodeToString([]byte(url))
	now := time.Now()
	expires := time.Now().AddDate(0, 1, 0)

	if l, _ := s.storage.FindLink(ctx, uuid); l != nil {
		// Link with this UUID already stored. Update it and return back
		upd := shortly.LinkUpdate{ExpiresAt: &expires}
		link, err := s.storage.UpdateLink(ctx, uuid, upd)
		if err != nil {
			return nil, fmt.Errorf("GenerateShortenedLink: %w", err)
		}
		return link, nil
	}

	var link = &shortly.Link{
		UUID:      uuid,
		URL:       url,
		CreatedAt: now,
		ExpiresAt: expires,
	}
	// Store new link object
	if err := s.storage.SaveLink(ctx, link); err != nil {
		return nil, fmt.Errorf("GenerateShortenedLink: %w", err)
	}

	return link, nil
}

func (s *LinkService) GetOriginalLink(ctx context.Context, uuid string) (*shortly.Link, error) {
	link, err := s.storage.FindLink(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("GetOriginalLink: %w", err)
	}

	return link, nil
}

func (s *LinkService) AddHit(ctx context.Context, uuid string) error {
	link, err := s.storage.FindLink(ctx, uuid)
	if err != nil {
		return fmt.Errorf("AddHit: %w", err)
	}

	count := link.Count + 1
	upd := shortly.LinkUpdate{Count: &count}

	_, err = s.storage.UpdateLink(ctx, uuid, upd)
	if err != nil {
		return fmt.Errorf("AddHit: %w", err)
	}

	return nil
}
