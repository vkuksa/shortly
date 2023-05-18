package inmem

import (
	"context"
	"testing"
	"time"

	"github.com/vkuksa/shortly"
)

func TestSaveLink(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		storage := NewLinkStorage()

		l1 := &shortly.Link{UUID: "1", URL: "test-1.com"}
		if err := storage.SaveLink(context.Background(), l1); err != nil {
			t.Fatal(err)
		}

		l2 := &shortly.Link{UUID: "2", URL: "test-2.com"}
		if err := storage.SaveLink(context.Background(), l2); err != nil {
			t.Fatal(err)
		}

		if len(storage.stor) != 2 {
			t.Fatal("expected storage size 2")
		}
	})

	t.Run("EINVALID", func(t *testing.T) {
		storage := NewLinkStorage()

		if err := storage.SaveLink(context.Background(), nil); err == nil {
			t.Fatal("expected error")
		} else if shortly.ErrorCode(err) != shortly.ERRINVALID {
			t.Fatalf("unexpected error: %#v", err)
		}

	})

	t.Run("EEXISTS", func(t *testing.T) {
		storage := NewLinkStorage()

		l := &shortly.Link{UUID: "1", URL: "test.com"}
		_ = storage.SaveLink(context.Background(), l)
		if err := storage.SaveLink(context.Background(), l); err == nil {
			t.Fatal("expected error")
		} else if shortly.ErrorCode(err) != shortly.ERRCONFLICT {
			t.Fatalf("unexpected error: %#v", err)
		}
	})
}

func TestFindLink(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		storage := NewLinkStorage()

		if err := storage.SaveLink(context.Background(), &shortly.Link{UUID: "1", URL: "test-1.com"}); err != nil {
			t.Fatal(err)
		}

		if err := storage.SaveLink(context.Background(), &shortly.Link{UUID: "2", URL: "test-2.com"}); err != nil {
			t.Fatal(err)
		}

		if len(storage.stor) != 2 {
			t.Fatal("expected storage size 2")
		}

		link, err := storage.FindLink(context.Background(), "2")
		if err != nil {
			t.Fatal(err)
		} else if link.URL != "test-2.com" {
			t.Fatal("Wrong origin of retrieved link")
		}
	})

	t.Run("ENOTFOUND", func(t *testing.T) {
		storage := NewLinkStorage()

		if _, err := storage.FindLink(context.Background(), "0"); err == nil {
			t.Fatal("expected error")
		} else if shortly.ErrorCode(err) != shortly.ERRNOTFOUND {
			t.Fatalf("unexpected error: %#v", err)
		}
	})
}

func TestDeleteLink(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		storage := NewLinkStorage()

		if err := storage.SaveLink(context.Background(), &shortly.Link{UUID: "1", URL: "test-1.com"}); err != nil {
			t.Fatal(err)
		}

		if err := storage.SaveLink(context.Background(), &shortly.Link{UUID: "2", URL: "test-2.com"}); err != nil {
			t.Fatal(err)
		}

		if len(storage.stor) != 2 {
			t.Fatal("expected storage size 2")
		}

		if err := storage.DeleteLink(context.Background(), "1"); err != nil {
			t.Fatal(err)
		} else if _, exists := storage.stor["1"]; exists {
			t.Fatal("link was not removed from underlying data structure")
		} else if len(storage.stor) != 1 {
			t.Fatal("expected storage size 1")
		}

		if err := storage.DeleteLink(context.Background(), "2"); err != nil {
			t.Fatal(err)
		} else if len(storage.stor) != 0 {
			t.Fatal("expected storage size 0")
		}
	})

	t.Run("ENOTFOUND", func(t *testing.T) {
		storage := NewLinkStorage()

		if err := storage.DeleteLink(context.Background(), "0"); err == nil {
			t.Fatal("expected error")
		} else if shortly.ErrorCode(err) != shortly.ERRNOTFOUND {
			t.Fatalf("unexpected error: %#v", err)
		}
	})
}

func TestUpdateLink(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		storage := NewLinkStorage()

		uuid := "1"
		if err := storage.SaveLink(context.Background(), &shortly.Link{UUID: uuid, URL: "test-1.com", Count: 3}); err != nil {
			t.Fatal(err)
		}

		count := 4
		now := time.Now()
		tomorrow := now.AddDate(0, 0, 1)
		upd := shortly.LinkUpdate{
			Count:     &count,
			CreatedAt: &now,
			ExpiresAt: &tomorrow,
		}
		if link, err := storage.UpdateLink(context.Background(), uuid, upd); err != nil {
			t.Fatal(err)
		} else if storage.stor[uuid].Count != 4 || storage.stor[uuid].CreatedAt != now || storage.stor[uuid].ExpiresAt != tomorrow {
			t.Fatal("Inmem value was not updated")
		} else if link.Count != 4 || link.CreatedAt != now || link.ExpiresAt != tomorrow {
			t.Fatal("Returned value is not updated")
		}
	})

	t.Run("ENOTFOUND", func(t *testing.T) {
		storage := NewLinkStorage()

		if _, err := storage.UpdateLink(context.Background(), "0", shortly.LinkUpdate{}); err == nil {
			t.Fatal("expected error")
		} else if shortly.ErrorCode(err) != shortly.ERRNOTFOUND {
			t.Fatalf("unexpected error: %#v", err)
		}
	})
}
