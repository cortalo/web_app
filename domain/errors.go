package domain

import "errors"

var (
	ErrUsernameExist = errors.New("username exist")
	ErrWrongPassword = errors.New("wrong password")
)
