package link

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
)

type Factory struct {
	repo    Repository
	encoder Encoder
}

func NewFactory(r Repository, e Encoder) Factory {
	return Factory{repo: r, encoder: e}
}

func (f Factory) NewLinkFromOriginal(ctx context.Context, original string) (*ShortenedLink, error) {
	shortened, err := f.encoder.Encode(original)
	if err != nil {
		return nil, err
	}
	slog.Debug("Shortened link", slog.String("shortened", shortened), slog.String("original", original))

	link, err := f.getLink(ctx, shortened)
	if err != nil {
		if errors.Is(err, ErrNotFound) || link == nil {
			link = NewShortenedLink(f.repo, original, shortened)
			if err := f.repo.Store(ctx, link); err != nil {
				return nil, fmt.Errorf("store link: %w", err)
			}
		} else {
			return nil, err
		}
	}
	return link, nil
}

func (f Factory) NewLinkFromShortened(ctx context.Context, shortened string) (*ShortenedLink, error) {
	return f.getLink(ctx, shortened)
}

func (f Factory) getLink(ctx context.Context, shortened string) (*ShortenedLink, error) {
	link, err := f.repo.Get(ctx, shortened)
	if err != nil {
		return nil, fmt.Errorf("get link: %w", err)
	}
	link.Repository = f.repo
	slog.Debug("Retrieved link", slog.Any("link", link))
	return link, nil
}
