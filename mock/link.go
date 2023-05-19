package mock

import (
	"context"

	shortly "github.com/vkuksa/shortly/internal/domain"
)

var _ shortly.LinkService = (*LinkService)(nil)

type LinkService struct {
	GenerateShortenedLinkFn func(ctx context.Context, url string) (*shortly.Link, error)
	GetOriginalLinkFn       func(ctx context.Context, uuid string) (*shortly.Link, error)
	AddHitFn                func(ctx context.Context, uuid string) error
}

func (s *LinkService) GenerateShortenedLink(ctx context.Context, url string) (*shortly.Link, error) {
	return s.GenerateShortenedLinkFn(ctx, url)
}

func (s *LinkService) GetOriginalLink(ctx context.Context, uuid string) (*shortly.Link, error) {
	return s.GetOriginalLinkFn(ctx, uuid)
}

func (s *LinkService) AddHit(ctx context.Context, uuid string) error {
	return s.AddHitFn(ctx, uuid)
}
