package usecase

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/vkuksa/shortly/internal/domain"
)

type MockLinkRepository struct {
	mock.Mock
}

func (m *MockLinkRepository) GetLink(ctx context.Context, uuid string) (*domain.Link, error) {
	args := m.Called(ctx, uuid)
	return args.Get(0).(*domain.Link), args.Error(1)
}

func (m *MockLinkRepository) StoreLink(ctx context.Context, link *domain.Link) error {
	args := m.Called(ctx, link)
	return args.Error(0)
}

// func TestLinkUseCase(t *testing.T) {
// 	mockRepo := &MockLinkRepository{Stor: make(map[string]domain.Link)}
// 	uc := NewLinkUseCase(mockRepo)
// 	ctx := context.Background()

// 	// Test GenerateShortenedLink
// 	t.Run("GenerateShortenedLink", func(t *testing.T) {
// 		// Define test cases
// 		tests := []struct {
// 			name        string
// 			url         string
// 			expectError bool
// 			// Add other necessary fields
// 		}{
// 			{
// 				name:        "Empty URL",
// 				url:         "",
// 				expectError: true,
// 			},
// 			// Add more test cases as needed
// 		}

// 		for _, tt := range tests {
// 			t.Run(tt.name, func(t *testing.T) {
// 				// Mock setup if necessary
// 				// ...

// 				link, err := uc.GenerateShortenedLink(ctx, tt.url)
// 				if tt.expectError {
// 					assert.Error(t, err)
// 				} else {
// 					assert.NoError(t, err)
// 					assert.NotNil(t, link)
// 					// Additional assertions based on expected outcome
// 				}
// 			})
// 		}
// 	})

// 	// Test GetOriginalLink
// 	t.Run("GetOriginalLink", func(t *testing.T) {
// 		// Define test cases
// 		tests := []struct {
// 			name        string
// 			uuid        string
// 			expectError bool
// 			// Add other necessary fields
// 		}{
// 			// Add test cases here
// 		}

// 		for _, tt := range tests {
// 			t.Run(tt.name, func(t *testing.T) {
// 				// Mock setup if necessary
// 				// ...

// 				link, err := uc.GetOriginalLink(ctx, tt.uuid)
// 				if tt.expectError {
// 					assert.Error(t, err)
// 				} else {
// 					assert.NoError(t, err)
// 					assert.NotNil(t, link)
// 					// Additional assertions based on expected outcome
// 				}
// 			})
// 		}
// 	})
// }
