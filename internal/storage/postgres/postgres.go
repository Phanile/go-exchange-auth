package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Phanile/go-exchange-auth/internal/domain/models"
	"github.com/Phanile/go-exchange-auth/internal/storage"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(postgresConfig string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("postgres", postgresConfig)

	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) Connection() *sql.DB {
	return s.db
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "storage.postgres.SaveUser"

	var id int64
	err := s.db.QueryRowContext(ctx, "insert into Users(email, pass_hash) values ($1, $2) returning id", email, passHash).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("%s : %w", op, err)
	}

	return id, nil
}

func (s *Storage) UserById(ctx context.Context, id int64) (*models.User, error) {
	const op = "storage.postgres.UserById"

	var user models.User

	err := s.db.QueryRowContext(ctx, "select * from Users where id = $1", id).Scan(&user.Id, &user.Email, &user.PassHash)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s : %w", op, storage.ErrUserNotFound)
		}

		return nil, fmt.Errorf("%s : %w", op, err)
	}

	return &user, nil
}

func (s *Storage) UserByEmail(ctx context.Context, email string) (*models.User, error) {
	const op = "storage.postgres.UserByEmail"

	var user models.User

	err := s.db.QueryRowContext(ctx, "select * from Users where email = $1", email).Scan(&user.Id, &user.Email, &user.PassHash)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s : %w", op, storage.ErrUserNotFound)
		}

		return nil, fmt.Errorf("%s : %w", op, err)
	}

	return &user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userId int64) (bool, error) {
	const op = "storage.postgres.IsAdmin"

	var exists bool
	err := s.db.QueryRowContext(ctx, "select exists(select 1 from Admins where user_id = $1)", userId).Scan(&exists)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return exists, nil
		}

		return exists, fmt.Errorf("%s : %w", op, err)
	}

	return exists, nil
}
