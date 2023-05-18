package business

import (
	"context"
	"testing"

	"github.com/vkuksa/shortly"
	"github.com/vkuksa/shortly/mock"
)

func TestLinkService(t *testing.T) {
	lstor := &mock.LinkStorage{}
	lserv := NewLinkService(lstor)
	cleanStor := func() {
		lstor.Stor = make(map[string]*shortly.Link)
	}

	t.Run("GenerateShortenedLink", func(t *testing.T) {
		cleanStor()

		if _, err := lserv.GenerateShortenedLink(context.Background(), ""); err == nil {
			t.Fatal("expected error")
		} else if shortly.ErrorCode(err) != shortly.ERRINVALID {
			t.Fatalf("unexpected error: %#v", err)
		}

		url := "example.com"
		// First execution should store link object and return valid data
		link, err := lserv.GenerateShortenedLink(context.Background(), url)
		if err != nil {
			t.Fatal(err)
		} else if link.UUID == "" || link.CreatedAt.IsZero() || link.ExpiresAt.IsZero() {
			t.Fatal("link object not initialised")
		}

		prevExpiresAt := link.ExpiresAt
		// Second execution should update existing object's ExpiresAt field
		if link, err = lserv.GenerateShortenedLink(context.Background(), url); err != nil {
			t.Fatal(err)
		} else if link.ExpiresAt == prevExpiresAt {
			t.Fatal("link object not updated")
		}
	})

	t.Run("GetOriginalLink", func(t *testing.T) {
		cleanStor()

		uuid := "test"
		url := "example.com"
		if _, err := lserv.GetOriginalLink(context.Background(), uuid); err == nil {
			t.Fatal("expected error")
		} else if shortly.ErrorCode(err) != shortly.ERRNOTFOUND {
			t.Fatalf("unexpected error: %#v", err)
		}

		lstor.Stor[uuid] = &shortly.Link{UUID: uuid, URL: url}

		if link, err := lserv.GetOriginalLink(context.Background(), uuid); err != nil {
			t.Fatal(err)
		} else if link == nil || link.URL != url || link.UUID != uuid {
			t.Fatal("invalid link retrieved")
		}
	})

	t.Run("Add Hit", func(t *testing.T) {
		cleanStor()

		uuid := "test"
		if err := lserv.AddHit(context.Background(), uuid); err == nil {
			t.Fatal("expected error")
		} else if shortly.ErrorCode(err) != shortly.ERRNOTFOUND {
			t.Fatalf("unexpected error: %#v", err)
		}

		lstor.Stor["test"] = &shortly.Link{UUID: "test"}

		if err := lserv.AddHit(context.Background(), uuid); err != nil {
			t.Fatal(err)
		} else if lstor.Stor["test"].Count != 1 {
			t.Fatal("count was not incremented")
		}
	})
}
