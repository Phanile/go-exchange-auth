package suite

import (
	"context"
	"github.com/Phanile/go-exchange-auth/internal/app"
	"github.com/Phanile/go-exchange-auth/internal/config"
	authv1 "github.com/Phanile/go-exchange-protos/generated/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"net"
	"strconv"
	"testing"
)

const grpcHost = "localhost"

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient authv1.AuthClient
}

func NewSuite(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadByPath("../config/config.yaml")

	//t.Setenv("JWT_PRIVATE_KEY", "%your_key%")
	//t.Setenv("PGSQL_CONNECTION_STRING", "%your_conn_string%")

	application := app.NewApp(slog.Default(), cfg) // ZAP?
	go application.GRPCServer.MustRun()

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	conn, err := grpc.NewClient(grpcAddress(cfg), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		t.Fatal("grpc connection failed: ", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: authv1.NewAuthClient(conn),
	}
}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(int(cfg.GRPC.Port)))
}
