package jwt_test

import (
	"testing"
	"time"

	"github.com/avran02/authentication/internal/config"
	"github.com/avran02/authentication/internal/models"
	jwtGenerator "github.com/avran02/authentication/internal/pkg/jwt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

var (
	cfg = config.JWT{
		Secret:     "supersecretkey",
		AccessExp:  3600,
		RefreshExp: 86400,
	}
	gen    = jwtGenerator.NewJwtGenerator(cfg)
	userID = "user123"
)

func TestJwtGenerator_Generate(t *testing.T) {
	accessToken, accessTokenID, refreshToken, _, err := gen.Generate(userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)
	assert.NotEmpty(t, accessTokenID)
}

func TestJwtGenerator_ParseAccessToken(t *testing.T) {
	accessToken, _, _, _, err := gen.Generate(userID)
	assert.NoError(t, err)

	claims, err := gen.ParseAccessToken(accessToken)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.Subject)
}

func TestJwtGenerator_ParseAccessToken_ExpiredToken(t *testing.T) {
	claims := models.AccessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Minute)),
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signedToken, err := token.SignedString([]byte(cfg.Secret))
	assert.NoError(t, err)

	_, err = gen.ParseAccessToken(signedToken)
	assert.Error(t, err)
}

func TestJwtGenerator_ParseRefreshToken(t *testing.T) {
	_, accessTokenID, refreshToken, _, err := gen.Generate(userID)
	assert.NoError(t, err)

	claims, err := gen.ParseRefreshToken(refreshToken)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.Subject)
	assert.Equal(t, accessTokenID, claims.AccessTokenID)
}

func TestJwtGenerator_ParseRefreshToken_ExpiredToken(t *testing.T) {
	claims := models.RefreshTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Minute)),
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signedToken, err := token.SignedString([]byte(cfg.Secret))
	assert.NoError(t, err)

	_, err = gen.ParseRefreshToken(signedToken)
	assert.Error(t, err)
}
