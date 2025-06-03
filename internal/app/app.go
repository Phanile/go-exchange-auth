package app

import (
	"database/sql"
	grpcApp "github.com/Phanile/go-exchange-auth/internal/app/grpc"
	"github.com/Phanile/go-exchange-auth/internal/config"
	"github.com/Phanile/go-exchange-auth/internal/services/auth"
	"github.com/Phanile/go-exchange-auth/internal/storage/postgres"
	"github.com/golang-jwt/jwt"
	"github.com/pressly/goose/v3"
	"log/slog"
	"os"
)

type App struct {
	GRPCServer *grpcApp.App
}

func NewApp(log *slog.Logger, config *config.Config) *App {
	dsn := os.Getenv("PGSQL_CONNECTION_STRING")

	if dsn == "" {
		panic("PGSQL_CONNECTION_STRING environment variable not set")
	}

	storage, err := postgres.NewStorage(dsn)

	if err != nil {
		panic(err)
	}

	privatePem := os.Getenv("JWT_PRIVATE_KEY")

	if privatePem == "" {
		panic("JWT_PRIVATE_KEY environment variable not set")
	}

	secretKey, errParse := jwt.ParseRSAPrivateKeyFromPEM([]byte(privatePem))

	if errParse != nil {
		panic(errParse)
	}

	runMigrations(storage.Connection())

	authService := auth.NewAuth(log, storage, storage, config.TokenTTL, secretKey)
	gRPCApp := grpcApp.NewGRPCApp(log, authService, config.GRPC.Port)

	return &App{
		GRPCServer: gRPCApp,
	}
}

func runMigrations(db *sql.DB) {
	goose.SetBaseFS(nil)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		panic(err)
	}
}
