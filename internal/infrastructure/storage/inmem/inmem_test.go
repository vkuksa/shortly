package inmem_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vkuksa/shortly/internal/domain"
	"github.com/vkuksa/shortly/internal/infrastructure/storage/inmem"
)

func TestStorage(t *testing.T) {
	tests := []struct {
		name           string
		prepareStorage func(*inmem.Storage)
		uuid           string
		link           *domain.Link
		expectError    bool
	}{
		{
			name:           "Store and Get valid link",
			prepareStorage: func(_ *inmem.Storage) {},
			uuid:           "valid-uuid",
			link:           &domain.Link{UUID: "valid-uuid", URL: "http://example.com"},
			expectError:    false,
		},
		{
			name:           "Get non-existing link",
			prepareStorage: func(_ *inmem.Storage) {},
			uuid:           "non-existent-uuid",
			link:           nil,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := inmem.NewStorage()
			tt.prepareStorage(storage)

			if tt.link != nil {
				err := storage.StoreLink(context.Background(), tt.link)
				assert.Equal(t, tt.expectError, err != nil, "StoreLink() error")
			}

			got, err := storage.GetLink(context.Background(), tt.uuid)
			if tt.expectError {
				assert.Error(t, err, "GetLink() expected error")
			} else {
				assert.NoError(t, err, "GetLink() unexpected error")
				assert.NotNil(t, got, "GetLink() expected a non-nil result")
				gotBytes, _ := json.Marshal(got)
				wantBytes, _ := json.Marshal(tt.link)
				assert.Equal(t, string(wantBytes), string(gotBytes), "GetLink() result mismatch")
			}
		})
	}
}
