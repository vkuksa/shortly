package test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vkuksa/shortly/internal/link"
)

func serviceURL() string {
	if url := os.Getenv("SERVICE_URL"); url != "" {
		return url
	}
	// Default to service container address within Docker network.
	return "http://shortly-svc:8081"
}

func Test_Functional(t *testing.T) {
	const ExpectedLocation = "https://www.google.com/"
	var shortened string

	tests := []struct {
		name         string
		method       string
		path         func() string // using func to dynamically calculate it with uuid
		body         io.Reader
		expectedCode int
		validate     func(*testing.T, *http.Response)
	}{
		{
			name:         "Link Shortening",
			method:       http.MethodPost,
			path:         func() string { return serviceURL() + "/links" },
			body:         strings.NewReader(fmt.Sprintf(`{"url": "%s"}`, ExpectedLocation)),
			expectedCode: http.StatusCreated,
			validate: func(t *testing.T, res *http.Response) {
				var link link.ShortenedLink
				assert.NoError(t, json.NewDecoder(res.Body).Decode(&link), "error decoding response body")
				assert.NotEmpty(t, link.Shortened, "UUID should not be empty")
				shortened = link.Shortened
			},
		},
		{
			name:         "Redirrection",
			method:       http.MethodGet,
			path:         func() string { return fmt.Sprintf("%s/%s", serviceURL(), shortened) },
			body:         nil,
			expectedCode: http.StatusFound,
			validate: func(t *testing.T, res *http.Response) {
				assert.NotEmpty(t, res.Header["Location"], "location header empty")
				location := res.Header["Location"]
				assert.Equal(t, ExpectedLocation, location[0], "location header empty")
			},
		},
		{
			name:         "Retrieval",
			method:       http.MethodGet,
			path:         func() string { return fmt.Sprintf("%s/links/%s", serviceURL(), shortened) },
			body:         nil,
			expectedCode: http.StatusOK,
			validate: func(t *testing.T, res *http.Response) {
				var link link.ShortenedLink
				if err := json.NewDecoder(res.Body).Decode(&link); err != nil {
					t.Errorf("error decoding response body: %v", err)
				}

				assert.Equal(t, 1, link.Hits, "invalid hit counter")
				assert.Equal(t, ExpectedLocation, link.Original, "url")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, tc.path(), tc.body)
			assert.NoError(t, err, "new request")

			client := &http.Client{
				// Prevent client from following redirects.
				CheckRedirect: func(*http.Request, []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}
			res, err := client.Do(req)
			if err != nil {
				t.Fatalf("error making request: %v", err)
			}

			assert.Equal(t, tc.expectedCode, res.StatusCode, "unexpected status code")

			if tc.validate != nil {
				tc.validate(t, res)
			}
		})
	}
}
