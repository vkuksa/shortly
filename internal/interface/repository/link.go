package repository

import (
	"context"

	"github.com/vkuksa/shortly/internal/domain"
)

type DataStore interface {
	GetLink(ctx context.Context, uuid string) (*domain.Link, error)
	StoreLink(ctx context.Context, link *domain.Link) error
	IncHit(ctx context.Context, uuid string) error
}

type LinkRepository struct {
	db DataStore
}

func New(db DataStore) *LinkRepository {
	return &LinkRepository{db: db}
}

func (r *LinkRepository) GetLink(ctx context.Context, uuid string) (*domain.Link, error) {
	if err := ValidateKey(uuid); err != nil {
		return nil, err
	}

	return r.db.GetLink(ctx, uuid)
}

func (r *LinkRepository) StoreLink(ctx context.Context, link *domain.Link) error {
	if err := ValidateKey(link.UUID); err != nil {
		return err
	}

	return r.db.StoreLink(ctx, link)
}

func (r *LinkRepository) IncHit(ctx context.Context, uuid string) error {
	if err := ValidateKey(uuid); err != nil {
		return err
	}

	return r.db.IncHit(ctx, uuid)
}
