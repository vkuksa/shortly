package redis

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/vkuksa/shortly/pkg/storage"
)

const (
	DefaultAddress = "localhost:6379"
)

type Client[V any] struct {
	c *redis.Client
}

type Options struct {
	Address  string `toml:"address"`
	Password string `toml:"password"`
	DB       int    `toml:"db"`
}

// NewClient creates a new Redis client
func NewClient[V any](options Options) (*Client[V], error) {
	client := redis.NewClient(&redis.Options{
		Addr:     options.Address,
		Password: options.Password,
		DB:       options.DB,
	})

	if err := client.Ping().Err(); err != nil {
		return nil, err
	}

	return &Client[V]{c: client}, nil
}

// Set stores the given value for the given key.
// Values are automatically marshalled to JSON
// The key must not be ""
func (c Client[V]) Set(k string, v V) error {
	if err := storage.ValidateKey(k); err != nil {
		return err
	}

	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	err = c.c.Set(k, string(data), 0).Err()
	if err != nil {
		return err
	}
	return nil
}

// Retrieves a stored value by a given key
// Returns a false, nil if no value have been found for a given key
// Returns an error if it occured during retrieving of value
// Expects keys that are not ""
func (c Client[V]) Get(k string, v *V) (found bool, err error) {
	if err := storage.ValidateKey(k); err != nil {
		return false, err
	}

	dataString, err := c.c.Get(k).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}

	if err = json.Unmarshal([]byte(dataString), v); err != nil {
		return true, fmt.Errorf("Get: %w", err)
	}

	return true, nil
}

// Deletes a key-value pair from a storage
// Returns an error if given key is not valid or update operation failed
func (c Client[V]) Delete(k string) error {
	if err := storage.ValidateKey(k); err != nil {
		return err
	}

	_, err := c.c.Del(k).Result()
	return err
}

// Close closes the client.
// It must be called to release any open resources.
func (c Client[V]) Close() error {
	return c.c.Close()
}
