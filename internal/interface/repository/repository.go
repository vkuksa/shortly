package repository

import (
	"errors"
)

var (
	ErrKeyIsEmpty    = errors.New("key is empty")
	ErrValueNotFound = errors.New("value not found")
)

func ValidateKey(k string) error {
	if k == "" {
		return ErrKeyIsEmpty
	}

	return nil
}
