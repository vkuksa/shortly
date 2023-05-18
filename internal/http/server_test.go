package http_test

import (
	"context"
	"io"
	"net/http"
	"testing"

	shortlyhttp "github.com/vkuksa/shortly/internal/http"
	"github.com/vkuksa/shortly/mock"
)

// Server represents a test wrapper for internal/http.Server
// It attaches mocks to the server & initializes on a random port.
type Server struct {
	*shortlyhttp.Server

	LinkService *mock.LinkService
}

// MustOpenServer is a test helper function for starting a new test HTTP server.
// Fail on error.
func MustOpenServer(tb testing.TB) *Server {
	tb.Helper()

	ls := &mock.LinkService{}

	// Initialize wrapper and set test configuration settings.
	s := &Server{Server: shortlyhttp.NewServer(), LinkService: ls}
	s.Server.LinkService = ls

	// Begin running test server.
	if err := s.Open(); err != nil {
		tb.Fatal(err)
	}
	return s
}

// MustCloseServer is a test helper function for shutting down the server.
// Fail on error.
func MustCloseServer(tb testing.TB, s *Server) {
	tb.Helper()
	if err := s.Close(); err != nil {
		tb.Fatal(err)
	}
}

// MustNewRequest creates a new HTTP request using the server's base URL and
// attaching a user session based on the context.
func (s *Server) MustNewRequest(tb testing.TB, ctx context.Context, method, url string, body io.Reader) *http.Request {
	tb.Helper()

	// Create new net/http request with server's base URL.
	r, err := http.NewRequest(method, s.URL()+url, body)
	if err != nil {
		tb.Fatal(err)
	}
	return r
}
