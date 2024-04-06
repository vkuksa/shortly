package test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vkuksa/shortly/internal/domain"
)

func Test_Functional(t *testing.T) {
	const ExpectedLocation = "https://www.google.com/"
	var uuid domain.UUID

	tests := []struct {
		name         string
		method       string
		path         func() string // using func to dynamically calculate it with uuid
		body         io.Reader
		expectedCode int
		validate     func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:         "Link Shortening",
			method:       http.MethodPost,
			path:         func() string { return "/links" },
			body:         strings.NewReader(fmt.Sprintf(`{"url": "%s"}`, ExpectedLocation)),
			expectedCode: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var link domain.Link
				assert.NoError(t, json.NewDecoder(w.Body).Decode(&link), "error decoding response body")
				assert.NotEmpty(t, link.UUID, "UUID should not be empty")
				uuid = link.UUID
			},
		},
		{
			name:         "Redirrection",
			method:       http.MethodGet,
			path:         func() string { return fmt.Sprintf("/%s", uuid) },
			body:         nil,
			expectedCode: http.StatusFound,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.NotEmpty(t, w.Result().Header["Location"], "location header empty")
				location := w.Result().Header["Location"]
				assert.Equal(t, ExpectedLocation, location[0], "location header empty")
			},
		},
		{
			name:         "Retrieval",
			method:       http.MethodGet,
			path:         func() string { return fmt.Sprintf("/links/%s", uuid) },
			body:         nil,
			expectedCode: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var link domain.Link
				if err := json.NewDecoder(w.Body).Decode(&link); err != nil {
					t.Errorf("error decoding response body: %v", err)
				}

				assert.Equal(t, 1, link.Count, "invalid usage counter")
				assert.Equal(t, ExpectedLocation, link.URL, "url")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, tc.path(), tc.body)
			assert.NoError(t, err, "new request")

			w := httptest.NewRecorder()
			httpServer.Serve(w, req)

			assert.Equal(t, tc.expectedCode, w.Code, "unexpected status code")

			if tc.validate != nil {
				tc.validate(t, w)
			}
		})
	}
}
