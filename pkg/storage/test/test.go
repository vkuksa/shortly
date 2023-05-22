//nolint

package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vkuksa/shortly/pkg/storage"
	"github.com/vkuksa/shortly/pkg/storage/bbolt"
	"github.com/vkuksa/shortly/pkg/storage/inmem"
)

func SetGetDelete(tb testing.TB, storage storage.Storage[int]) {
	var value int
	found, err := storage.Get("", &value)
	assert.Error(tb, err, "Trying to get by invalid key")
	assert.False(tb, found, "Value should not be found")

	err = storage.Delete("")
	assert.Error(tb, err, "Delete on invalid key")

	found, err = storage.Get("key0", &value)
	assert.NoError(tb, err, "Failed to get value")
	assert.False(tb, found, "A value was found, but no value was expected")
	assert.Equal(tb, 0, value, "Value expected to not be initialised")

	err = storage.Set("", 10)
	assert.Error(tb, err, "Trying to set invalid key")

	err = storage.Set("key1", 10)
	assert.NoError(tb, err, "Failed to set value")

	found, err = storage.Get("key1", &value)
	assert.NoError(tb, err, "Failed to get value")
	assert.True(tb, found, "Expected value to be found")
	assert.Equal(tb, 10, value, "Expected value to be 10")

	err = storage.Delete("key1")
	assert.NoError(tb, err, "Failed to delete value")

	found, err = storage.Get("key1", &value)
	assert.NoError(tb, err, "Failed to get value")
	assert.False(tb, found, "Expected value to be deleted")
}

//nolint:errcheck
func Close(tb testing.TB, storage storage.Storage[int]) {
	key := "test"
	value := 1

	err := storage.Close()
	assert.NoError(tb, err, "Failed to close storage")

	// Attempt to set, get, delete a value after closing
	err = storage.Set(key, value)
	assert.Error(tb, err)

	_, err = storage.Get(key, &value)
	assert.Error(tb, err)

	err = storage.Delete(key)
	assert.Error(tb, err)
}

func MustCreateStorage[V any](tb testing.TB, kind string) storage.Storage[V] {
	var err error
	switch kind {
	case "inmem":
		return inmem.NewStorage[V]()
	case "bbolt":
		stor, err := bbolt.NewStorage[V]("test.db", "test")
		if err != nil {
			break
		}

		return stor
	default:
		tb.Fatalf("Storage %s is not supported", kind)
		return nil
	}

	tb.Fatalf("Failed to create %s storage with %s", kind, err.Error())
	return nil
}

//nolint:errcheck
func MustCleanupStorage[V any](tb testing.TB, s storage.Storage[V]) {
	if err := s.Close(); err != nil {
		tb.Fatal(err)
	}

	switch s := s.(type) {
	case *inmem.Storage[V]:
		return
	case *bbolt.Storage[V]:
		s.Cleanup()
		return
	default:
		tb.Fatalf("Trying to cleanup unknown storage type")
	}

}
