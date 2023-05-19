package storage

import (
	"errors"
)

// Represents an interface that serves storage-related operation
type Storage[V any] interface {
	Set(key string, value V) error

	Get(key string, value *V) (bool, error)

	Delete(key string) error

	// Gracefully closes storage
	Close() error
}

func ValidateKey(key string) error {
	if key == "" {
		return errors.New("Empty key provided")
	}

	return nil
}
