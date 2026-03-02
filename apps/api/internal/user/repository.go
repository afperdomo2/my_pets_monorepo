package user

import (
	"context"
	"errors"
)

// ErrNotFound is returned when a user does not exist in the database.
var ErrNotFound = errors.New("user not found")

// ErrEmailTaken is returned when attempting to register a duplicate email.
var ErrEmailTaken = errors.New("email already in use")

// Repository defines the persistence contract for the user domain.
// All methods accept a context to support cancellation and timeouts.
type Repository interface {
	GetAll(ctx context.Context) ([]User, error)
	GetByID(ctx context.Context, id int) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	Create(ctx context.Context, payload UserPayload, isSystemUser bool) (User, error)
	Update(ctx context.Context, id int, payload UserPayload) (User, error)
	Delete(ctx context.Context, id int) error
	// SystemUserExists reports whether at least one system user has been created.
	SystemUserExists(ctx context.Context) (bool, error)
}
