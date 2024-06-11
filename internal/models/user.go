// Package models ...
package models

import (
	"context"
	"errors"
	"fmt"
)

var (
	// ErrUserAlreadyExists is not unique login error
	ErrUserAlreadyExists = errors.New("login is already exists")

	// ErrUserNotExists user not found error
	ErrUserNotExists = errors.New("user is not exists")
)

// User is a user model
type User struct {
	ID           string
	Login        string
	PasswordHash string
}

// UserDTO is a data transfer object for user
type UserDTO struct {
	Login    string
	Password string
}

// UserStorage is a interface for working with users in database
type UserStorage interface {
	// CreateUser is a method for adding new user
	CreateUser(ctx context.Context, u *UserDTO) (*User, error)
	// GetUser is a method for getting user
	GetUser(ctx context.Context, u *UserDTO) (*User, error)
}

// GetUser is a method for getting user
func (u *UserDTO) GetUser(ctx context.Context, db UserStorage) (*User, error) {
	user, err := db.GetUser(ctx, u)
	if err != nil {
		if !errors.Is(err, ErrUserNotExists) {
			return nil, fmt.Errorf("getting user error: %w", err)
		}
		return nil, fmt.Errorf("%w: %v", ErrUserNotExists, err)
	}
	return user, nil
}

// CreateUser is a method for registering new user
func (u *UserDTO) CreateUser(ctx context.Context, db UserStorage) (*User, error) {
	if u.Login == "" {
		return nil, ErrUserNotExists
	}
	user, err := db.CreateUser(ctx, u)
	if err != nil {
		if !errors.Is(err, ErrUserAlreadyExists) {
			return nil, fmt.Errorf("adding user error: %w", err)
		}
		return nil, fmt.Errorf("%w: %v", ErrUserAlreadyExists, err)
	}
	return user, nil
}
