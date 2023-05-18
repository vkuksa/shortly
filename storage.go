package shortly

import (
	"context"
	"time"
)

// Represents an interface that serves storage-related operation upon links
type LinkStorage interface {
	SaveLink(ctx context.Context, link *Link) error

	FindLink(ctx context.Context, uuid string) (*Link, error)

	DeleteLink(ctx context.Context, uuid string) error

	UpdateLink(ctx context.Context, uuid string, upd LinkUpdate) (*Link, error)

	// Opens up storage for data processing
	Open(ctx context.Context) error

	// Gracefully closes storage
	Close() error
}

// Fields as pointers to emulate optionality
type LinkUpdate struct {
	Count     *int
	CreatedAt *time.Time
	ExpiresAt *time.Time
}
