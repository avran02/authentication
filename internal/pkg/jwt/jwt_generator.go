package jwt

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/avran02/authentication/internal/config"
	"github.com/avran02/authentication/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Generator interface {
	Generate(userID string) (accessToken, accessTokenID, refreshToken string, err error)
	ParseAccessToken(token string) (models.AccessTokenClaims, error)
	ParseRefreshToken(token string) (models.RefreshTokenClaims, error)
}

type jwtGenerator struct {
	config config.JWT
}

func (j *jwtGenerator) Generate(userID string) (accessToken, accessTokenID, refreshToken string, err error) {
	slog.Info("pkg.jwt.Generate")
	accessClaims := j.newAccessClaims(userID)
	unsignedAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, accessClaims)
	accessToken, err = unsignedAccessToken.SignedString([]byte(j.config.Secret))
	if err != nil {
		return "", "", "", fmt.Errorf("pkg.jwt.Generate: failed to sign token: %w", err)
	}

	refreshClaims := j.newRefreshClaims(userID, accessClaims.ID)
	unsignedRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS512, refreshClaims)
	refreshToken, err = unsignedRefreshToken.SignedString([]byte(j.config.Secret))
	if err != nil {
		return "", "", "", fmt.Errorf("pkg.jwt.Generate: failed to sign token: %w", err)
	}

	return accessToken, accessClaims.ID, refreshToken, nil
}

func (j *jwtGenerator) ParseAccessToken(token string) (models.AccessTokenClaims, error) {
	slog.Info("pkg.jwt.ParseAccessToken")
	if token == "" {
		return models.AccessTokenClaims{}, ErrEmptyToken
	}

	parsedToken, err := jwt.ParseWithClaims(token, &models.AccessTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(j.config.Secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS512.Alg()}))
	if err != nil {
		return models.AccessTokenClaims{}, fmt.Errorf("pkg.jwt.ValidateAccessToken: failed to parse token: %w", err)
	}

	if !parsedToken.Valid {
		return models.AccessTokenClaims{}, ErrInvalidToken
	}

	claims, ok := parsedToken.Claims.(*models.AccessTokenClaims)
	if !ok {
		return models.AccessTokenClaims{}, ErrInvalidToken
	}

	if time.Now().After(claims.ExpiresAt.Time) {
		return models.AccessTokenClaims{}, ErrExpiredToken
	}

	return *claims, nil
}

func (j *jwtGenerator) ParseRefreshToken(token string) (models.RefreshTokenClaims, error) {
	slog.Info("pkg.jwt.ParseRefreshToken")

	if token == "" {
		return models.RefreshTokenClaims{}, ErrEmptyToken
	}

	parsedToken, err := jwt.ParseWithClaims(token, &models.RefreshTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(j.config.Secret), nil
	})
	if err != nil {
		return models.RefreshTokenClaims{}, fmt.Errorf("pkg.jwt.ParseRefreshToken: failed to parse token: %w", err)
	}

	if !parsedToken.Valid {
		return models.RefreshTokenClaims{}, ErrInvalidToken
	}

	claims, ok := parsedToken.Claims.(*models.RefreshTokenClaims)
	if !ok {
		return models.RefreshTokenClaims{}, ErrInvalidToken
	}

	slog.Debug("pkg.jwt.ParseRefreshToken", "claims", claims)

	if time.Now().After(claims.ExpiresAt.Time) {
		return models.RefreshTokenClaims{}, ErrExpiredToken
	}

	return *claims, nil
}

func (j *jwtGenerator) newAccessClaims(userID string) models.AccessTokenClaims {
	tokenLifetime := time.Duration(j.config.AccessExp) * time.Second
	accessTokenExpiresAt := jwt.NewNumericDate(time.Now().Add(tokenLifetime))

	return models.AccessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: accessTokenExpiresAt,
			Subject:   userID,
			ID:        uuid.New().String(),
		},
	}
}

func (j *jwtGenerator) newRefreshClaims(userID, accessTokenID string) models.RefreshTokenClaims {
	tokenLifetime := time.Duration(j.config.RefreshExp) * time.Second
	refreshTokenExpiresAt := jwt.NewNumericDate(time.Now().Add(tokenLifetime))

	return models.RefreshTokenClaims{
		AccessTokenID: accessTokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: refreshTokenExpiresAt,
			Subject:   userID,
			ID:        uuid.New().String(),
		},
	}
}

func NewJwtGenerator(config config.JWT) Generator {
	return &jwtGenerator{
		config: config,
	}
}
