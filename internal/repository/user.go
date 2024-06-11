package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/SerjRamone/vaultme/internal/models"
)

const uniqueConstraintLogin = "user_login_key"

// CreateUser is a method for registering new user
func (db *DB) CreateUser(ctx context.Context, u *models.UserDTO) (*models.User, error) {
	row := db.pool.QueryRow(
		ctx,
		`INSERT INTO user(login, password) VALUES ($1, $2) RETURNING id, login, password;`,
		u.Login,
		u.Password,
	)

	user := models.User{}
	if err := row.Scan(&user.ID, &user.Login, &user.PasswordHash); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) && pgErr.ConstraintName == uniqueConstraintLogin {
				return nil, models.ErrUserAlreadyExists
			}
			return nil, fmt.Errorf("user with same login is already exists: %w", err)
		}
		return nil, fmt.Errorf("row scan error: %w", err)
	}

	return &user, nil
}

// GetUser is a method for getting user
func (db *DB) GetUser(ctx context.Context, u *models.UserDTO) (*models.User, error) {
	row := db.pool.QueryRow(
		ctx,
		`SELECT id, login, password FROM users WHERE login = $1;`,
		u.Login,
	)

	user := models.User{}
	if err := row.Scan(&user.ID, &user.Login, &user.PasswordHash); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrUserNotExists
		}
		return nil, fmt.Errorf("row scan error: %w", err)
	}

	return &user, nil
}
