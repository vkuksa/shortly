//nolint

package test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vkuksa/shortly/pkg/storage"
	"github.com/vkuksa/shortly/pkg/storage/bbolt"
	"github.com/vkuksa/shortly/pkg/storage/inmem"
	"github.com/vkuksa/shortly/pkg/storage/redis"

	goredis "github.com/go-redis/redis"
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

func MustCreateStorage[V any](tb testing.TB, kind string, o interface{}) (storage storage.Storage[V], closer func()) {
	var err error
	switch kind {
	case "inmem":
		return inmem.NewStorage[V](), func() {}
	case "bbolt":
		stor, err := bbolt.NewStorage[V](o.(bbolt.Options))
		if err != nil {
			break
		}

		return stor, func() {
			_ = stor.Close()
			_ = os.Remove("test.db")
		}
	case "redis":
		o := o.(redis.Options)

		// Create Redis client
		client := goredis.NewClient(&goredis.Options{
			Addr:     o.Address,
			Password: o.Password,
			DB:       o.DB,
		})
		err := client.Ping().Err()
		if err != nil {
			tb.Fatalf("An error occurred during testing the connection to the server: %v\n", err)
		}

		stor, err := redis.NewClient[V](o)
		if err != nil {
			break
		}

		return stor, func() {
			err := client.FlushAll().Err()
			if err != nil {
				panic(err)
			}

			_ = client.Close()
			_ = stor.Close()
		}
	default:
		tb.Fatalf("Storage %s is not supported", kind)
		return nil, func() {}
	}

	tb.Fatalf("Failed to create %s storage with %s", kind, err.Error())
	return nil, func() {}
}
