package shortly

import (
	"context"
	"time"
)

// Represents a link in a system
type Link struct {
	// UUID is also a shortened link
	UUID string `json:"uuid"`
	// Original url to redirrect to
	URL string `json:"url"`

	// Amount of times that link was reddirrected
	Count int `json:"count"`

	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// Represents an interface that serves domain-related operations
type LinkService interface {
	// Generates a shortened link for the provided URL and stores the Link data
	GenerateShortenedLink(ctx context.Context, url string) (*Link, error)

	// Retrieves the original link for the provided shortened version
	GetOriginalLink(ctx context.Context, uuid string) (*Link, error)

	// Increments counter of a link
	AddHit(ctx context.Context, uuid string) error
}
