package grpc

import (
	"fmt"
	authGRPC "github.com/Phanile/go-exchange-auth/internal/grpc/auth"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       uint
}

func NewGRPCApp(log *slog.Logger, port uint) *App {
	grpcServer := grpc.NewServer()
	authGRPC.Register(grpcServer)

	return &App{
		log:        log,
		gRPCServer: grpcServer,
		port:       port,
	}
}

func (app *App) MustRun() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func (app *App) Run() error {
	const op = "grpcApp.Run"

	app.log.With(
		slog.String("op", op),
		slog.Int("port", int(app.port)),
	)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", app.port))

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	app.log.Info("grpc server is running on", slog.String("addr", listener.Addr().String()))

	listenErr := app.gRPCServer.Serve(listener)

	if listenErr != nil {
		return fmt.Errorf("%s: %w", op, listenErr)
	}

	return nil
}

func (app *App) Stop() {
	const op = "grpcApp.Stop"

	app.log.With(
		slog.String("op", op),
	)

	app.log.Info("grpc server is shutting down")

	app.gRPCServer.GracefulStop()
}
