package link

import (
	"context"
	"time"
)

type Repository interface {
	Get(ctx context.Context, shortened string) (*ShortenedLink, error)
	Store(ctx context.Context, url *ShortenedLink) error
	AddHit(ctx context.Context, id string) error
	UpdateExpiration(ctx context.Context, id string, expiresAt time.Time) error
}

type Encoder interface {
	Encode(original string) (string, error)
}
