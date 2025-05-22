package auth

import (
	"context"
	"github.com/Phanile/go-exchange-protos/generated/go/auth"
	"google.golang.org/grpc"
)

type ServerAPI struct {
	authv1.UnimplementedAuthServer
}

func Register(grpc *grpc.Server) {
	authv1.RegisterAuthServer(grpc, &ServerAPI{})
}

func (s *ServerAPI) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	return &authv1.LoginResponse{
		Token: "Mem",
	}, nil
}

func (s *ServerAPI) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	return &authv1.RegisterResponse{
		UserId: 1,
	}, nil
}

func (s *ServerAPI) IsAdmin(ctx context.Context, req *authv1.AdminRequest) (*authv1.AdminResponse, error) {
	return &authv1.AdminResponse{
		IsAdmin: false,
	}, nil
}
