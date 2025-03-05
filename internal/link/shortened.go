package link

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ShortenedLink struct {
	Id string `json:"id" bson:"_id"`

	Shortened string `json:"shortened" bson:"shortened"`
	Original  string `json:"original" bson:"original"`

	Hits int `json:"hits" bson:"hits"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt" bson:"expiresAt"`

	Repository Repository `json:"-" bson:"-"`
}

func (u *ShortenedLink) ResetExpiration(ctx context.Context) error {
	expiresAt := time.Now().Add(24 * time.Hour)
	err := u.Repository.UpdateExpiration(ctx, u.Id, expiresAt)
	if err != nil {
		return err
	}
	u.ExpiresAt = expiresAt
	return nil
}

func (u *ShortenedLink) Hit(ctx context.Context) error {
	if err := u.Repository.AddHit(ctx, u.Id); err != nil {
		return err
	}
	u.Hits++
	return nil
}

func NewShortenedLink(ctx context.Context, repo Repository, original, shortened string) *ShortenedLink {
	now := time.Now()
	return &ShortenedLink{
		Id:         uuid.New().String(),
		Shortened:  shortened,
		Original:   original,
		CreatedAt:  now,
		ExpiresAt:  now.Add(24 * time.Hour),
		Repository: repo,
	}
}
