package link

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/vkuksa/shortly/internal/domain"
)

type Repository interface {
	GetLink(ctx context.Context, uuid domain.UUID) (*domain.Link, error)
	StoreLink(ctx context.Context, link *domain.Link) error
	IncHit(ctx context.Context, uuid domain.UUID) error
}

type UseCase struct {
	repo Repository
}

func NewUseCase(r Repository) *UseCase {
	return &UseCase{repo: r}
}

func (uc *UseCase) Shorten(ctx context.Context, url string) (*domain.Link, error) {
	if url == "" {
		return nil, ErrBadInput
	}

	uuid := domain.NewUUIDfromString(url)
	// Update if we already have link with given uuid
	link, err := uc.repo.GetLink(ctx, uuid)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return nil, fmt.Errorf("get link: %w", err)
	} else if (err != nil && errors.Is(err, ErrNotFound)) || link == nil {
		link = domain.NewLink(url)
	} else {
		link.ResetExpiration()
	}

	if err := uc.repo.StoreLink(ctx, link); err != nil {
		return nil, fmt.Errorf("GenerateShortenedLink: %w", err)
	}

	return link, nil
}

func (uc *UseCase) Retrieve(ctx context.Context, uuid string) (*domain.Link, error) {
	link, err := uc.repo.GetLink(ctx, domain.UUID(uuid))
	if err != nil {
		return nil, fmt.Errorf("GetOriginalLink: %w", err)
	} else if link == nil {
		return nil, ErrNotFound
	}

	if err = uc.repo.IncHit(ctx, domain.UUID(uuid)); err != nil {
		slog.Error("hit incrementation failed", slog.Any("error", err))
	}

	return link, nil
}
