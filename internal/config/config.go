package config

import (
	"log"
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Server struct {
	Host     string
	GRPCPort string
	HTTPPort string
	LogLevel string
}

type DB struct {
	Host     string
	Port     string
	User     string
	Password string
}

type JWT struct {
	Secret     string
	AccessExp  int
	RefreshExp int
}

type Config struct {
	Server Server
	DB     DB
	JWT    JWT
	CORS   CORSConfig
	Cookie CookieConfig
}

func New() *Config {
	if os.Getenv("LOAD_DOT_ENV") != "false" {
		if err := godotenv.Load(); err != nil {
			log.Fatal("can't load .env file")
		}
		slog.Info("Loaded .env file")
	}

	accessExpStr := os.Getenv("JWT_ACCESS_EXP")
	refreshExpStr := os.Getenv("JWT_REFRESH_EXP")
	accessExp, err := strconv.Atoi(accessExpStr)
	if err != nil {
		slog.Warn("JWT_ACCESS_EXP is not an int, using default value: 3600")
		accessExp = 3600
	}

	refreshExp, err := strconv.Atoi(refreshExpStr)
	if err != nil {
		slog.Warn("JWT_REFRESH_EXP is not an int, using default value: 86400")
		refreshExp = 86400
	}

	slog.Info("env config loaded")

	corsConfig, cookieConfig := getYmlConfig()
	slog.Info("config.yml loaded")

	return &Config{
		Server: Server{
			Host:     os.Getenv("SERVER_HOST"),
			GRPCPort: os.Getenv("SERVER_GRPC_PORT"),
			HTTPPort: os.Getenv("SERVER_HTTP_PORT"),
			LogLevel: os.Getenv("SERVER_LOG_LEVEL"),
		},
		DB: DB{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
		},
		JWT: JWT{
			Secret:     os.Getenv("JWT_SECRET"),
			AccessExp:  accessExp,
			RefreshExp: refreshExp,
		},
		CORS:   corsConfig,
		Cookie: cookieConfig,
	}
}
