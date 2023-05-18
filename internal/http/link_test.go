package http_test

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/vkuksa/shortly"
)

func TestLinkEndpoints(t *testing.T) {
	// Start the mocked HTTP test server.
	s := MustOpenServer(t)
	defer MustCloseServer(t, s)

	testUrl := "http://example.com"
	testUuid := "test"
	exmpl := &shortly.Link{
		URL:  testUrl,
		UUID: testUuid,
	}
	s.LinkService.GetOriginalLinkFn = func(ctx context.Context, uuid string) (*shortly.Link, error) {
		if uuid != testUuid {
			return nil, shortly.Errorf(shortly.ERRNOTFOUND, "Expected")
		}

		return exmpl, nil
	}
	s.LinkService.AddHitFn = func(ctx context.Context, uuid string) error {
		return nil
	}
	s.LinkService.GenerateShortenedLinkFn = func(ctx context.Context, url string) (*shortly.Link, error) {
		if url != testUrl {
			return nil, shortly.Errorf(shortly.ERRINVALID, "Expected")
		}

		return exmpl, nil
	}

	t.Run("Handle Root", func(t *testing.T) {
		// Issue request for see if root handler is present
		resp, err := http.DefaultClient.Do(s.MustNewRequest(t, context.Background(), "GET", "/", nil))
		if err != nil {
			t.Fatal(err)
			// We expect server error, as template will not be loaded properly during testing
			// The main point is to see whether endpoint is present and can be accessed
		} else if got, want := resp.StatusCode, http.StatusOK; got != want {
			t.Fatalf("StatusCode=%v, want %v", got, want)
		}
	})

	t.Run("Handle Link Redirrection", func(t *testing.T) {
		// Issue request for performing redirrection via shortened url
		resp, err := http.DefaultClient.Do(s.MustNewRequest(t, context.Background(), "GET", "/test", nil))
		if err != nil {
			t.Fatal(err)
		} else if got, want := resp.StatusCode, http.StatusOK; got != want {
			t.Fatalf("StatusCode=%v, want %v", got, want)
		}

		// Issue request for querying expected error
		resp, err = http.DefaultClient.Do(s.MustNewRequest(t, context.Background(), "GET", "/test_notfound", nil))
		if err != nil {
			t.Fatal(err)
		} else if got, want := resp.StatusCode, http.StatusNotFound; got != want {
			t.Fatalf("StatusCode=%v, want %v", got, want)
		}
	})

	t.Run("Handle Link Storage", func(t *testing.T) {
		// Issue request for storing and generating of shortened link
		formData := url.Values{}
		formData.Set("url", testUrl)
		req := s.MustNewRequest(t, context.Background(), "POST", "/", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		} else if got, want := resp.StatusCode, http.StatusOK; got != want {
			t.Fatalf("StatusCode=%v, want %v", got, want)
		}

		req = s.MustNewRequest(t, context.Background(), "POST", "/", nil)
		// Issue request for expecting an error, when invalid url provided
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		} else if got, want := resp.StatusCode, http.StatusBadRequest; got != want {
			t.Fatalf("StatusCode=%v, want %v", got, want)
		}
	})
}
