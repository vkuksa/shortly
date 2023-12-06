package usecase

import (
	"testing"
)

func TestLinkService(t *testing.T) {
	// t.Run("GenerateShortenedLink", func(t *testing.T) {
	// 	if _, err := lserv.GenerateShortenedLink(context.Background(), ""); err == nil {
	// 		t.Fatal("expected error")
	// 	} else if ErrorCode(err) != ErrInvalid {
	// 		t.Fatalf("unexpected error: %#v", err)
	// 	}

	// 	url := "example.com"
	// 	// First execution should store link object and return valid data
	// 	link, err := lserv.GenerateShortenedLink(context.Background(), url)
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	} else if link.UUID == "" || link.CreatedAt.IsZero() || link.ExpiresAt.IsZero() {
	// 		t.Fatal("link object not initialised")
	// 	}

	// 	prevExpiresAt := link.ExpiresAt
	// 	// Second execution should update existing object's ExpiresAt field
	// 	if link, err = lserv.GenerateShortenedLink(context.Background(), url); err != nil {
	// 		t.Fatal(err)
	// 	} else if link.ExpiresAt == prevExpiresAt {
	// 		t.Fatal("link object not updated")
	// 	}
	// })

	// t.Run("GetOriginalLink", func(t *testing.T) {
	// 	if _, err := lserv.GetOriginalLink(context.Background(), ""); err == nil {
	// 		t.Fatal("expected error")
	// 	}

	// 	uuid := "test"
	// 	url := "example.com"
	// 	if _, err := lserv.GetOriginalLink(context.Background(), uuid); err == nil {
	// 		t.Fatal("expected error")
	// 	} else if ErrorCode(err) != ErrNotFound {
	// 		t.Fatalf("unexpected error: %#v", err)
	// 	}

	// 	lstor.Stor[uuid] = domain.Link{UUID: uuid, URL: url}

	// 	if link, err := lserv.GetOriginalLink(context.Background(), uuid); err != nil {
	// 		t.Fatal(err)
	// 	} else if link == nil || link.URL != url || link.UUID != uuid {
	// 		t.Fatal("invalid link retrieved")
	// 	}
	// })

	// t.Run("Add Hit", func(t *testing.T) {
	// 	uuid := "test"
	// 	if err := lserv.AddHit(context.Background(), uuid); err == nil {
	// 		t.Fatal("expected error")
	// 	} else if ErrorCode(err) != ErrNotFound {
	// 		t.Fatalf("unexpected error: %#v", err)
	// 	}

	// 	lstor.Stor["test"] = domain.Link{UUID: "test"}

	// 	if err := lserv.AddHit(context.Background(), uuid); err != nil {
	// 		t.Fatal(err)
	// 	} else if lstor.Stor["test"].Count != 1 {
	// 		t.Fatal("count was not incremented")
	// 	}
	// })
}
