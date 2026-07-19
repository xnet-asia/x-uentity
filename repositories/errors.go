package repositories

import "errors"

var (
	ErrNotFound     = errors.New("entity not found")
	ErrInvalidKey   = errors.New("invalid key provided")
	ErrCreateFailed = errors.New("failed to create entity")
)
