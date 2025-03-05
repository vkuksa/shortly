package inmem_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/vkuksa/shortly/internal/infrastructure/storage/inmem"
	"github.com/vkuksa/shortly/internal/link"
)

func TestStorage(t *testing.T) {
	ctx := context.Background()
	storage := inmem.NewStorage()

	l := link.ShortenedLink{
		Id:        uuid.New().String(),
		Shortened: "short1",
		Original:  "original1",
		Hits:      0,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	t.Run("Store and Get", func(t *testing.T) {
		testLink := l
		if err := storage.Store(ctx, &testLink); err != nil {
			t.Fatalf("Store() error = %v", err)
		}

		got, err := storage.Get(ctx, testLink.Shortened)
		if err != nil {
			t.Fatalf("Get() error = %v", err)
		}

		if got.Id != testLink.Id {
			t.Errorf("Get() got link ID = %v, want %v", got.Id, testLink.Id)
		}
		if got.Original != testLink.Original {
			t.Errorf("Get() got link Original = %v, want %v", got.Original, testLink.Original)
		}
		if got.Hits != testLink.Hits {
			t.Errorf("Get() got link Hits = %v, want %v", got.Hits, testLink.Hits)
		}
	})

	t.Run("AddHit increments link hits", func(t *testing.T) {
		testLink := l
		err := storage.AddHit(ctx, testLink.Id)
		if err != nil {
			t.Fatalf("AddHit() error = %v", err)
		}

		got, err := storage.Get(ctx, testLink.Shortened)
		if err != nil {
			t.Fatalf("Get() error = %v", err)
		}

		// We started at 0, so we expect 1 now
		if got.Hits != (testLink.Hits + 1) {
			t.Errorf("Hits after AddHit() = %v, want %v", got.Hits, testLink.Hits+1)
		}
	})

	t.Run("UpdateExpiration changes link expiration", func(t *testing.T) {
		testLink := l
		newExpiration := time.Now().Add(48 * time.Hour)
		err := storage.UpdateExpiration(ctx, testLink.Id, newExpiration)
		if err != nil {
			t.Fatalf("UpdateExpiration() error = %v", err)
		}

		got, err := storage.Get(ctx, testLink.Shortened)
		if err != nil {
			t.Fatalf("Get() error = %v", err)
		}

		// Compare times with a tolerance, or use .Equal if you want an exact match
		if !got.ExpiresAt.Equal(newExpiration) {
			t.Errorf("ExpiresAt after UpdateExpiration() = %v, want %v", got.ExpiresAt, newExpiration)
		}
	})

	t.Run("Get non-existent link returns ErrNotFound", func(t *testing.T) {
		_, err := storage.Get(ctx, "does-not-exist")
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
		if err != link.ErrNotFound {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})
}
