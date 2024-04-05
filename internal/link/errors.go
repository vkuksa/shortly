package link

import (
	"errors"
)

var (
	ErrInternal = errors.New("internal")
	ErrNotFound = errors.New("not_found")
	ErrBadInput = errors.New("bad_input")
	ErrConflict = errors.New("conflict")
)
