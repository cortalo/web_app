package domain

import "errors"

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrInvalidName    = errors.New("invalid name")
)
