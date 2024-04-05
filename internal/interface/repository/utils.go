package repository

import "github.com/vkuksa/shortly/internal/link"

func ValidateKey(k string) error {
	if k == "" {
		return link.ErrBadInput
	}

	return nil
}
