package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/Phanile/go-exchange-auth/internal/domain/models"
	"github.com/Phanile/go-exchange-auth/internal/lib/jwt"
	"github.com/Phanile/go-exchange-auth/internal/storage"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExist          = errors.New("user already exists")
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	tokenTTL     time.Duration
}

func NewAuth(log *slog.Logger, userSaver UserSaver, userProvider UserProvider, tokenTTL time.Duration) *Auth {
	return &Auth{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		tokenTTL:     tokenTTL,
	}
}

func (a *Auth) Login(ctx context.Context, email string, password string) (token string, err error) {
	const op = "auth.Login"

	a.log.With(
		slog.String("op", op),
	)

	user, errUser := a.userProvider.UserByEmail(ctx, email)

	if errUser != nil {
		if errors.Is(errUser, storage.ErrUserNotFound) {
			a.log.Error("user not found")

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		a.log.Error("failed to get user")

		return "", fmt.Errorf("%s: %w", op, errUser)
	}

	errPass := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password))

	if errPass != nil {
		a.log.Error("invalid credentials")

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	jwtToken := jwt.NewToken(user, a.tokenTTL, "")

	return jwtToken, nil
}

func (a *Auth) Register(ctx context.Context, email string, password string) (userId int64, err error) {
	const op = "auth.Register"

	a.log.With(
		slog.String("op", op),
	)

	hashedPass, e := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	if e != nil {
		a.log.Error("failed to hash password")
		return 0, fmt.Errorf("%s: %w", op, e)
	}

	id, saveErr := a.userSaver.SaveUser(ctx, email, hashedPass)

	if saveErr != nil {
		if errors.Is(saveErr, ErrUserExist) {
			return 0, fmt.Errorf("%s: %w", op, ErrUserExist)
		}

		a.log.Error("failed to save user")
		return 0, fmt.Errorf("%s: %w", op, saveErr)
	}

	return id, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userId int64) (isAdmin bool, err error) {
	const op = "auth.IsAdmin"

	a.log.With(
		slog.String("op", op),
	)

	admin, errAdmin := a.userProvider.IsAdmin(ctx, userId)

	if errAdmin != nil {
		return false, fmt.Errorf("%s: %w", op, errAdmin)
	}

	return admin, nil
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (userId int64, err error)
}

type UserProvider interface {
	UserByEmail(ctx context.Context, email string) (*models.User, error)
	UserById(ctx context.Context, id int64) (*models.User, error)
	IsAdmin(ctx context.Context, userId int64) (bool, error)
}
