package app

import (
	grpcApp "github.com/Phanile/go-exchange-auth/internal/app/grpc"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcApp.App
}

func NewApp(log *slog.Logger, port uint, tokenTTL time.Duration) *App {
	gRPCApp := grpcApp.NewGRPCApp(log, port)
	return &App{
		GRPCServer: gRPCApp,
	}
}
