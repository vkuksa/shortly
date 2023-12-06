package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/vkuksa/shortly/internal/domain"
	"github.com/vkuksa/shortly/internal/interface/repository"
)

type Options struct {
	Address  string
	Password string
	DB       int
}

type Handler struct {
	client *redis.Client
}

func NewClient(o Options) (*Handler, error) {
	c := redis.NewClient(&redis.Options{
		Addr:     o.Address,
		Password: o.Password,
		DB:       o.DB,
	})

	if err := c.Ping().Err(); err != nil {
		return nil, err
	}

	return &Handler{client: c}, nil
}

func (r *Handler) GetLink(_ context.Context, uuid string) (*domain.Link, error) {
	if err := repository.ValidateKey(uuid); err != nil {
		return nil, err
	}

	dataString, err := r.client.Get(uuid).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, repository.ErrValueNotFound
		}

		return nil, err
	}

	var result *domain.Link
	if err = json.Unmarshal([]byte(dataString), result); err != nil {
		return nil, fmt.Errorf("GetLink: %w", err)
	}

	return result, nil
}

func (r *Handler) StoreLink(_ context.Context, link *domain.Link) error {
	if err := repository.ValidateKey(link.UUID); err != nil {
		return err
	}

	data, err := json.Marshal(link)
	if err != nil {
		return err
	}

	err = r.client.Set(link.UUID, string(data), 0).Err()
	if err != nil {
		return err
	}

	return nil
}

// func (c *Handler) Delete(k string) error {
// 	_, err := c.client.Del(k).Result()
// 	return err
// }

// func (c *Handler) Close() error {
// 	return c.client.Close()
// }
