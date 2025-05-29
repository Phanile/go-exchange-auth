package app

import (
	"database/sql"
	grpcApp "github.com/Phanile/go-exchange-auth/internal/app/grpc"
	"github.com/Phanile/go-exchange-auth/internal/config"
	"github.com/Phanile/go-exchange-auth/internal/services/auth"
	"github.com/Phanile/go-exchange-auth/internal/storage/postgres"
	"github.com/pressly/goose/v3"
	"log/slog"
)

type App struct {
	GRPCServer *grpcApp.App
}

func NewApp(log *slog.Logger, config *config.Config) *App {
	storage, err := postgres.NewStorage(config.Dsn)

	if err != nil {
		panic(err)
	}

	runMigrations(storage.Connection())

	authService := auth.NewAuth(log, storage, storage, config.TokenTTL)
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
