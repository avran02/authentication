package service

import "errors"

var (
	ErrUserNotFound      = errors.New("user already exists")
	ErrTokenDoesntExist  = errors.New("token doesn't exist")
	ErrWrongCredentials  = errors.New("wrong credentials")
	ErrWrongTokensPair   = errors.New("wrong tokens pair")
	ErrUserAlreadyExists = errors.New("user already exists")
)
