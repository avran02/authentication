package service

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"

	"github.com/avran02/authentication/internal/models"
	"github.com/avran02/authentication/internal/pkg/jwt"
	"github.com/avran02/authentication/internal/repo"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(
		ctx context.Context,
		username, password string,
		email *string,
	) (id, accessToken, refreshToken string, err error)
	Login(ctx context.Context, username, password string) (id, accessToken, refreshToken string, err error)
	RefreshTokens(ctx context.Context, token string) (accessToken, refreshToken string, err error)
	ValidateToken(ctx context.Context, token string) (string, error)
	Logout(ctx context.Context, token string) (bool, error)
}

type service struct {
	repo repo.Repo
	jwt  jwt.Generator
}

func (s *service) Register(
	ctx context.Context,
	username, password string,
	email *string,
) (id, accessToken, refreshToken string, err error) {
	slog.Info("Registering user: " + username)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := s.repo.FindUserByUsername(ctx, username)
	if err != nil && !errors.Is(err, repo.ErrUserNotFound) {
		return "", "", "", fmt.Errorf("failed to find user: %w", err)
	}
	if user != nil {
		return "", "", "", ErrUserAlreadyExists
	}

	id = uuid.NewString()
	if err = s.repo.CreateUser(ctx, models.User{
		ID:       id,
		Email:    email,
		Username: username,
		Password: string(hashedPassword),
	}); err != nil {
		return "", "", "", fmt.Errorf("failed ti create user: %w", err)
	}

	// todo: refactor duplicates part of login
	accessToken, accessTokenID, refreshToken, err := s.jwt.Generate(id)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate tokens: %w", err)
	}

	if err = s.saveRefreshToken(ctx, id, accessTokenID, []byte(refreshToken)); err != nil {
		return "", "", "", fmt.Errorf("failed to save refresh token: %w", err)
	}

	return id, accessToken, refreshToken, nil
}

func (s *service) Login(ctx context.Context, username, password string) (id, accessToken, refreshToken string, err error) {
	slog.Info("Logging in user: " + username)
	user, err := s.repo.FindUserByUsername(ctx, username)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to find user: %w", err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", "", ErrWrongCredentials
	}

	if err = s.repo.DeleteAllUserTokens(ctx, user.ID); err != nil {
		return "", "", "", fmt.Errorf("failed to delete all user tokens: %w", err)
	}

	accessToken, accessTokenID, refreshToken, err := s.jwt.Generate(user.ID)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate tokens: %w", err)
	}

	if err = s.saveRefreshToken(ctx, user.ID, accessTokenID, []byte(refreshToken)); err != nil {
		return "", "", "", fmt.Errorf("failed to save refresh token: %w", err)
	}

	return user.ID, accessToken, refreshToken, nil
}

func (s *service) RefreshTokens(ctx context.Context, refreshTokenStr string) (newAccessToken, newRefreshToken string, err error) {
	slog.Info("authenticationService.RefreshTokens")
	refreshToken, err := s.jwt.ParseRefreshToken(refreshTokenStr)
	if err != nil {
		return "", "", fmt.Errorf("authenticationService.RefreshTokens: can't validate refresh token: %w", err)
	}

	writtenRefreshTokenHash, writtenAccessTokenID, err := s.repo.GetRefreshTokenInfo(ctx, refreshToken.Subject)
	slog.Debug("authenticationService.RefreshTokens", "writtenRefreshTokenHash", writtenRefreshTokenHash, "writtenAccessTokenID", writtenAccessTokenID)
	if err != nil {
		return "", "", fmt.Errorf("authenticationService.RefreshTokens: can't get refresh token info: %w", err)
	}

	hashedRefreshToken := sha256.New().Sum([]byte(refreshTokenStr))
	encodedRefreshToken := base64.RawStdEncoding.EncodeToString(hashedRefreshToken)
	if writtenRefreshTokenHash != encodedRefreshToken {
		return "", "", ErrTokenDoesntExist
	}

	if refreshToken.AccessTokenID != writtenAccessTokenID {
		slog.Error("wrong access token id", "writtenAccessTokenID", writtenAccessTokenID, "refreshToken.AccessTokenID", refreshToken.AccessTokenID)
		return "", "", fmt.Errorf("wrong access token id: %w", ErrWrongTokensPair)
	}

	newAccessToken, newAccessTokenID, newRefreshToken, err := s.jwt.Generate(refreshToken.Subject)
	if err != nil {
		return "", "", fmt.Errorf("authenticationService.RefreshTokens: can't generate new tokens: %w", err)
	}

	if err = s.repo.DeleteAllUserTokens(ctx, refreshToken.Subject); err != nil {
		return "", "", fmt.Errorf("authenticationService.RefreshTokens: can't delete all user tokens: %w", err)
	}

	if err = s.saveRefreshToken(ctx, refreshToken.Subject, newAccessTokenID, []byte(newRefreshToken)); err != nil {
		return "", "", fmt.Errorf("authenticationService.RefreshTokens: can't save new refresh token: %w", err)
	}

	return newAccessToken, newRefreshToken, nil
}

func (s *service) ValidateToken(ctx context.Context, token string) (string, error) {
	claims, err := s.jwt.ParseAccessToken(token)
	if err != nil {
		return "", fmt.Errorf("failed to parse access token: %w", err)
	}

	_, writtenAccessTokenID, err := s.repo.GetRefreshTokenInfo(ctx, claims.Subject)
	if err != nil {
		return "", fmt.Errorf("authenticationService.RefreshTokens: can't get refresh token info: %w", err)
	}
	if writtenAccessTokenID != claims.ID {
		return "", fmt.Errorf("wrong access token id: %w", ErrWrongTokensPair)
	}

	return claims.Subject, nil
}

func (s *service) Logout(ctx context.Context, token string) (bool, error) {
	claims, err := s.jwt.ParseAccessToken(token)
	if err != nil {
		return false, fmt.Errorf("failed to parse access token: %w", err)
	}

	_, writtenAccessTokenID, err := s.repo.GetRefreshTokenInfo(ctx, claims.Subject)
	if err != nil {
		return false, fmt.Errorf("authenticationService.RefreshTokens: can't get refresh token info: %w", err)
	}
	if writtenAccessTokenID != claims.ID {
		return false, fmt.Errorf("wrong access token id: %w", ErrWrongTokensPair)
	}

	err = s.repo.DeleteAllUserTokens(ctx, claims.Subject)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *service) saveRefreshToken(ctx context.Context, userID, accessTokenID string, refreshToken []byte) error {
	hasedRefreshToken := sha256.New().Sum(refreshToken)
	encodedRefreshToken := base64.RawStdEncoding.EncodeToString(hasedRefreshToken)
	return s.repo.WriteRefreshToken(ctx, userID, accessTokenID, encodedRefreshToken)
}

func New(repo repo.Repo, jwt jwt.Generator) Service {
	return &service{
		repo: repo,
		jwt:  jwt,
	}
}
