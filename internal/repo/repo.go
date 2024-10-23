package repo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"

	"github.com/avran02/authentication/internal/config"
	"github.com/avran02/authentication/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repo interface {
	CreateUser(ctx context.Context, user models.User) error
	FindUserByUsername(ctx context.Context, username string) (*models.User, error)
	DeleteAllUserTokens(ctx context.Context, userID string) error
	WriteRefreshToken(ctx context.Context, userID, accessTokenID, refreshToken string) error
	GetRefreshTokenInfo(ctx context.Context, userID string) (writtenRefreshTokenHash, writtenAccessTokenID string, err error)
}

type repo struct {
	client           *mongo.Client
	userCollection   *mongo.Collection
	tokensCollection *mongo.Collection
}

func (r *repo) CreateUser(ctx context.Context, user models.User) error {
	_, err := r.userCollection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}

func (r *repo) FindUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user *models.User
	err := r.userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *repo) DeleteAllUserTokens(ctx context.Context, userID string) error {
	_, err := r.tokensCollection.DeleteMany(ctx, bson.M{"userID": userID})
	if err != nil {
		return fmt.Errorf("failed to delete tokens: %w", err)
	}
	return nil
}

func (r *repo) WriteRefreshToken(ctx context.Context, userID, accessTokenID, refreshToken string) error {
	token := bson.M{
		"userID":        userID,
		"accessTokenID": accessTokenID,
		"refreshToken":  refreshToken,
	}

	_, err := r.tokensCollection.InsertOne(ctx, token)
	if err != nil {
		return fmt.Errorf("failed to insert token: %w", err)
	}
	return nil
}

func (r *repo) GetRefreshTokenInfo(ctx context.Context, userID string) (string, string, error) {
	var token bson.M
	if err := r.tokensCollection.FindOne(ctx, bson.M{"userID": userID}).Decode(&token); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", "", ErrTokenNotFound
		}
		return "", "", fmt.Errorf("failed to find token: %w", err)
	}

	return token["refreshToken"].(string), token["accessTokenID"].(string), nil
}

func New(conf *config.DB) Repo {
	client := mustConnectDB(conf)
	usersCollection := client.Database("auth").Collection("users")
	tokensCollection := client.Database("auth").Collection("tokens")

	return &repo{
		client:           client,
		userCollection:   usersCollection,
		tokensCollection: tokensCollection,
	}
}

func mustConnectDB(config *config.DB) *mongo.Client {
	dsn := getDsn(config)
	opt := options.Client().ApplyURI(dsn)
	client, err := mongo.Connect(context.Background(), opt)
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %s", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("failed to ping MongoDB: %s", err)
	}

	slog.Info("MongoDB connected")
	return client
}

func getDsn(config *config.DB) string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s", config.User, config.Password, config.Host, config.Port)
}
