package usecase

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"

	"github.com/vkuksa/shortly/internal/domain"
	"github.com/vkuksa/shortly/internal/interface/repository"
)

type LinkRepository interface {
	GetLink(ctx context.Context, uuid string) (*domain.Link, error)
	StoreLink(ctx context.Context, link *domain.Link) error
}

type LinkUseCase struct {
	repo LinkRepository
}

func NewLinkUseCase(r LinkRepository) *LinkUseCase {
	return &LinkUseCase{repo: r}
}

func (uc *LinkUseCase) GenerateShortenedLink(ctx context.Context, url string) (*domain.Link, error) {
	if url == "" {
		return nil, NewError(ErrInvalid, "invalid url value provided for shortening")
	}

	uuid := base64.URLEncoding.EncodeToString([]byte(url))

	// Check if we elready have a link for a given uuid
	link, err := uc.repo.GetLink(ctx, uuid)
	if err != nil && !errors.Is(err, repository.ErrValueNotFound) {
		return nil, fmt.Errorf("GenerateShortenedLink: %w", err)
	} else if err != nil && errors.Is(err, repository.ErrValueNotFound) || link == nil {
		link = domain.NewLink(uuid, url)
	} else {
		link.ResetExpiration()
	}

	if err := uc.repo.StoreLink(ctx, link); err != nil {
		return nil, fmt.Errorf("GenerateShortenedLink: %w", err)
	}

	return link, nil
}

func (uc *LinkUseCase) GetOriginalLink(ctx context.Context, uuid string) (*domain.Link, error) {
	link, err := uc.repo.GetLink(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("GetOriginalLink: %w", err)
	} else if link == nil {
		return nil, NewError(ErrNotFound, "No link stored for a given uuid")
	}

	if err = uc.addHit(ctx, *link); err != nil {
		log.Printf("GetOriginalLink: %s", err.Error())
	}

	return link, nil
}

// Passing link by value to add a hit to copy of original link object, and store it
func (uc *LinkUseCase) addHit(ctx context.Context, link domain.Link) error {
	link.Count++

	if err := uc.repo.StoreLink(ctx, &link); err != nil {
		return fmt.Errorf("AddHit: %w", err)
	}

	return nil
}
