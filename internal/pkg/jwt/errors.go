package jwt

import "errors"

var (
	ErrEmptyToken   = errors.New("token is empty")
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token is expired")
)
