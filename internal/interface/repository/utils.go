package repository

import (
	"github.com/vkuksa/shortly/internal/domain"
	"github.com/vkuksa/shortly/internal/link"
)

func ValidateKey(k domain.UUID) error {
	if k == "" {
		return link.ErrBadInput
	}

	return nil
}
