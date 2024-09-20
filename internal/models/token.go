package models

import "github.com/golang-jwt/jwt/v5"

type AccessTokenClaims struct {
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	AccessTokenID string `json:"access_token_id"`
	jwt.RegisteredClaims
}
