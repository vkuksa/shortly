//nolint

package inmem

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vkuksa/shortly/pkg/storage"
)

// Ensure service implements interface.
var _ storage.Storage[int] = (*Storage[int])(nil)

func TestSetAndGet(t *testing.T) {
	storage := NewStorage[int]()

	var value int
	found, err := storage.Get("", &value)
	assert.Error(t, err, "Trying to get by invalid key")
	assert.False(t, found, "Value should not be found")

	found, err = storage.Get("key0", &value)
	assert.NoError(t, err, "Failed to get value")
	assert.False(t, found, "A value was found, but no value was expected")
	assert.Equal(t, 0, value, "Value expected to not be initialised")

	err = storage.Set("", 10)
	assert.Error(t, err, "Trying to set invalid key")

	err = storage.Set("key1", 10)
	assert.NoError(t, err, "Failed to set value")

	found, err = storage.Get("key1", &value)
	assert.NoError(t, err, "Failed to get value")
	assert.True(t, found, "Expected value to be found")
	assert.Equal(t, 10, value, "Expected value to be 10")
}

func TestDelete(t *testing.T) {
	storage := NewStorage[int]()

	err := storage.Delete("")
	assert.Error(t, err, "Delete on invalid key")

	err = storage.Set("key1", 10)
	assert.NoError(t, err, "Failed to set value")

	err = storage.Delete("key1")
	assert.NoError(t, err, "Failed to delete value")

	var value int
	found, err := storage.Get("key1", &value)
	assert.NoError(t, err, "Failed to get value")
	assert.False(t, found, "Expected value to be deleted")
}

//nolint:errcheck
func TestClose(t *testing.T) {
	storage := NewStorage[int]()
	key := "test"
	value := 1

	err := storage.Close()
	assert.NoError(t, err, "Failed to close storage")

	// Attempt to set a value after closing
	assert.Panics(t, func() { storage.Set(key, value) })
	assert.Panics(t, func() { storage.Get(key, &value) })
	assert.Panics(t, func() { storage.Delete(key) })
}
