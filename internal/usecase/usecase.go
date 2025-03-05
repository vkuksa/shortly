package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/vkuksa/shortly/internal/link"
)

type UseCase struct {
	factory link.Factory
}

func New(f link.Factory) UseCase {
	return UseCase{factory: f}
}

type ShortenUrlInput struct {
	url string
}

func NewShortenUrlInput(url string) ShortenUrlInput {
	return ShortenUrlInput{url: url}
}

func (i ShortenUrlInput) Validate() error {
	if i.url == "" {
		return errors.New("url is empty")
	}
	return nil
}

func (uc *UseCase) ShortenUrl(ctx context.Context, input ShortenUrlInput) (*link.ShortenedLink, error) {
	if err := input.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", link.ErrBadInput, err)
	}

	link, err := uc.factory.NewLinkFromOriginal(ctx, input.url)
	if err != nil {
		return nil, err
	}

	err = link.ResetExpiration(ctx)
	if err != nil {
		slog.Error("failed to reset expiration", slog.Any("error", err))
	}
	return link, nil
}

type GetOriginalInput struct {
	shortened string
}

func NewGetOriginalInput(s string) GetOriginalInput {
	return GetOriginalInput{shortened: s}
}

func (i GetOriginalInput) Validate() error {
	if i.shortened == "" {
		return errors.New("shortened is empty")
	}
	return nil
}

func (uc *UseCase) GetOriginal(ctx context.Context, input GetOriginalInput) (*link.ShortenedLink, error) {
	link, err := uc.factory.NewLinkFromShortened(ctx, input.shortened)
	if err != nil {
		return nil, err
	}

	err = link.Hit(ctx)
	if err != nil {
		slog.Error("failed to reset expiration", slog.Any("error", err))
	}
	return link, nil
}

func (uc *UseCase) GetOriginalWithoutHit(ctx context.Context, input GetOriginalInput) (*link.ShortenedLink, error) {
	return uc.factory.NewLinkFromShortened(ctx, input.shortened)
}
