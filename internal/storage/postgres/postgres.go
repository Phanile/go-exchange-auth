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

	stmt, err := s.db.Prepare("insert into Users(email, pass_hash) values (?, ?)")

	if err != nil {
		return 0, fmt.Errorf("%s : %w", op, err)
	}

	defer stmt.Close()

	res, errExec := stmt.ExecContext(ctx, email, passHash)

	if errExec != nil {
		return 0, fmt.Errorf("%s : %w", op, errExec)
	}

	id, errId := res.LastInsertId()

	if errId != nil {
		return 0, fmt.Errorf("%s : %w", op, errId)
	}

	return id, nil
}

func (s *Storage) UserById(ctx context.Context, id int64) (*models.User, error) {
	const op = "storage.postgres.UserById"

	stmt, err := s.db.Prepare("select * from Users where id = ?")

	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}

	defer stmt.Close()

	var user models.User

	row := stmt.QueryRowContext(ctx, id)
	errScan := row.Scan(user.Id, user.Email, user.PassHash)

	if errScan != nil {
		if errors.Is(errScan, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s : %w", op, storage.ErrUserNotFound)
		}

		return nil, fmt.Errorf("%s : %w", op, errScan)
	}

	return &user, nil
}

func (s *Storage) UserByEmail(ctx context.Context, email string) (*models.User, error) {
	const op = "storage.postgres.UserByEmail"

	stmt, err := s.db.Prepare("select * from Users where email = ?")

	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}

	defer stmt.Close()

	var user models.User

	row := stmt.QueryRowContext(ctx, email)
	errScan := row.Scan(user.Id, user.Email, user.PassHash)

	if errScan != nil {
		if errors.Is(errScan, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s : %w", op, storage.ErrUserNotFound)
		}

		return nil, fmt.Errorf("%s : %w", op, errScan)
	}

	return &user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userId int64) (bool, error) {
	const op = "storage.postgres.IsAdmin"

	stmt, err := s.db.Prepare("select * from Admins where user_id = ?")

	if err != nil {
		return false, fmt.Errorf("%s : %w", op, err)
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, userId)
	errScan := row.Scan()

	if errScan != nil {
		if errors.Is(errScan, sql.ErrNoRows) {
			return false, nil
		}

		return false, fmt.Errorf("%s : %w", op, errScan)
	}

	return true, nil
}
