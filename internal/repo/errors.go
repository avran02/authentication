package repo

import "errors"

var (
	ErrUserNotFound  = errors.New("user does not exists")
	ErrTokenNotFound = errors.New("token doesn't exist")
)
