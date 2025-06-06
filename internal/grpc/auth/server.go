package auth

import (
	"context"
	"errors"
	"github.com/Phanile/go-exchange-auth/internal/lib/prometheus"
	"github.com/Phanile/go-exchange-auth/internal/services/auth"
	"github.com/Phanile/go-exchange-auth/internal/storage"
	"github.com/Phanile/go-exchange-protos/generated/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, email string, password string) (token string, err error)
	Register(ctx context.Context, email string, password string) (userId int64, err error)
	IsAdmin(ctx context.Context, userId int64) (isAdmin bool, err error)
}

type ServerAPI struct {
	authv1.UnimplementedAuthServer
	auth Auth
}

func Register(grpc *grpc.Server, auth Auth) {
	authv1.RegisterAuthServer(grpc, &ServerAPI{
		auth: auth,
	})
}

func (s *ServerAPI) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	prometheus.LoginAttempts.WithLabelValues("attempt_login").Inc()

	if req.GetEmail() == "" {
		prometheus.LoginErrors.WithLabelValues("empty_email").Inc()
		return nil, status.Error(codes.InvalidArgument, "Email is required")
	}

	if req.GetPassword() == "" {
		prometheus.LoginErrors.WithLabelValues("empty_password").Inc()
		return nil, status.Error(codes.InvalidArgument, "Password is required")
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword())

	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			prometheus.LoginErrors.WithLabelValues("invalid_credentials").Inc()
			return nil, status.Error(codes.InvalidArgument, "Invalid Credentials")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &authv1.LoginResponse{
		Token: token,
	}, nil
}

func (s *ServerAPI) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	prometheus.RegisterAttempts.WithLabelValues("attempt").Inc()

	if req.GetEmail() == "" {
		prometheus.RegisterErrors.WithLabelValues("empty_email").Inc()
		return nil, status.Error(codes.InvalidArgument, "Email is required")
	}

	if req.GetPassword() == "" {
		prometheus.RegisterErrors.WithLabelValues("empty_password").Inc()
		return nil, status.Error(codes.InvalidArgument, "Password is required")
	}

	userId, err := s.auth.Register(ctx, req.GetEmail(), req.GetPassword())

	if err != nil {
		if errors.Is(err, auth.ErrUserExist) {
			prometheus.RegisterErrors.WithLabelValues("user_exists").Inc()
			return nil, status.Error(codes.AlreadyExists, "User already exists")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &authv1.RegisterResponse{
		UserId: userId,
	}, nil
}

func (s *ServerAPI) IsAdmin(ctx context.Context, req *authv1.AdminRequest) (*authv1.AdminResponse, error) {
	if req.GetUserId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "UserId is less or equal to zero")
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())

	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "User not found")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &authv1.AdminResponse{
		IsAdmin: isAdmin,
	}, nil
}
