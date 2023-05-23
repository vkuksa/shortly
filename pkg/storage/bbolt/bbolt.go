package bbolt

import (
	"encoding/json"
	"fmt"

	bolt "go.etcd.io/bbolt"

	"github.com/vkuksa/shortly/pkg/storage"
)

type Storage[V any] struct {
	db         *bolt.DB
	bucketName string
}

type Options struct {
	File   string `toml:"file"`
	Bucket string `toml:"bucket"`
}

func NewStorage[V any](o Options) (*Storage[V], error) {
	// Open DB
	db, err := bolt.Open(o.File, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("NewStorage: %w", err)
	}

	// Create bucket
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(o.Bucket))
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("NewStorage: %w", err)
	}

	return &Storage[V]{db: db, bucketName: o.Bucket}, nil
}

// Saves data into bbolt storage
// Returns an error, if given key is not valid
func (s *Storage[V]) Set(k string, v V) error {
	if err := storage.ValidateKey(k); err != nil {
		return fmt.Errorf("Set: %w", err)
	}

	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("Set: %w", err)
	}

	err = s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.bucketName))
		return b.Put([]byte(k), data)
	})
	if err != nil {
		return err
	}

	return nil
}

// Retrieves a stored value by a given key
// Returns a false, nil if no value have been found for a given key
// Returns an error if it occured during retrieving of value
// Expects keys that are not ""
func (s *Storage[V]) Get(k string, v *V) (bool, error) {
	if err := storage.ValidateKey(k); err != nil {
		return false, fmt.Errorf("Get: %w", err)
	}

	var data []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.bucketName))
		data = b.Get([]byte(k))
		return nil
	})
	if err != nil {
		return false, fmt.Errorf("Get: %w", err)
	}

	// If no value was found return false
	if data == nil {
		return false, nil
	}

	if err = json.Unmarshal(data, v); err != nil {
		return true, fmt.Errorf("Get: %w", err)
	}

	return true, nil
}

// Deletes a key-value pair from a storage
// Returns an error if given key is not valid or update operation failed
func (s *Storage[V]) Delete(k string) error {
	if err := storage.ValidateKey(k); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.bucketName))
		return b.Delete([]byte(k))
	})
}

// Gracefull close of database
func (s *Storage[V]) Close() error {
	return s.db.Close()
}
