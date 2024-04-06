package bbolt

import (
	"context"
	"encoding/json"

	"github.com/vkuksa/shortly/internal/domain"
	"github.com/vkuksa/shortly/internal/link"
	bolt "go.etcd.io/bbolt"
)

type Options struct {
	File   string
	Bucket string
}

type Handler struct {
	db         *bolt.DB
	bucketName string
}

func NewHandler(o Options) (*Handler, error) {
	db, err := bolt.Open(o.File, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(o.Bucket))
		return err
	})
	if err != nil {
		return nil, err
	}

	return &Handler{db: db, bucketName: o.Bucket}, nil
}

func (s *Handler) StoreLink(_ context.Context, link *domain.Link) error {
	data, err := json.Marshal(link)
	if err != nil {
		return err
	}

	err = s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.bucketName))
		return b.Put([]byte(link.UUID), data)
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Handler) GetLink(_ context.Context, uuid domain.UUID) (*domain.Link, error) {
	var data []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.bucketName))
		data = b.Get([]byte(uuid))
		return nil
	})
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, link.ErrNotFound
	}

	var result *domain.Link
	if err = json.Unmarshal(data, result); err != nil {
		return nil, err
	}

	return result, nil
}

// func (s *Handler) Delete(k string) error {
// 	if err := storage.ValidateKey(k); err != nil {
// 		return fmt.Errorf("delete: %w", err)
// 	}

// 	return s.db.Update(func(tx *bolt.Tx) error {
// 		b := tx.Bucket([]byte(s.bucketName))
// 		return b.Delete([]byte(k))
// 	})
// }

// func (s *Handler) Close() error {
// 	return s.db.Close()
// }
