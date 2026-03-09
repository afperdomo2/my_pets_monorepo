package user

import (
	"context"
	"errors"

	"github.com/my-pets/api/internal/models"
)

// ErrNotFound is returned when a user does not exist in the database.
var ErrNotFound = errors.New("user not found")

// ErrEmailTaken is returned when attempting to register a duplicate email.
var ErrEmailTaken = errors.New("email already in use")

// ErrWrongProvider is returned when a user tries to log in with the wrong provider.
var ErrWrongProvider = errors.New("account registered with a different login method")

// Repository defines the persistence contract for the user domain.
// All methods accept a context to support cancellation and timeouts.
type Repository interface {
	GetAll(ctx context.Context) ([]models.User, error)
	// GetPaginated returns a page of users and the total count.
	GetPaginated(ctx context.Context, page, perPage int) ([]models.User, int64, error)
	GetByID(ctx context.Context, id string) (models.User, error)
	GetByEmail(ctx context.Context, email string) (models.User, error)
	Create(ctx context.Context, payload CreateUserPayload, isSystemUser bool) (models.User, error)
	Update(ctx context.Context, id string, payload UpdateUserPayload) (models.User, error)
	Delete(ctx context.Context, id string) error
	// HasUsers reports whether any user exists in the database (used for setup detection).
	HasUsers(ctx context.Context) (bool, error)
	// SystemUserExists reports whether at least one system user has been created.
	SystemUserExists(ctx context.Context) (bool, error)
	// GetPasswordByEmail retrieves the stored bcrypt hash for login verification.
	// The hash is never included in User responses.
	GetPasswordByEmail(ctx context.Context, email string) (hash string, err error)
	// UpsertByGoogleID creates or updates a user authenticated via Google OAuth.
	UpsertByGoogleID(ctx context.Context, info GoogleUserInfo, isSystemUser bool) (models.User, error)
}
