package pet

import (
	"context"
	"errors"

	"github.com/my-pets/api/internal/models"
)

// ErrNotFound is returned when a pet does not exist or does not belong to the requesting user.
var ErrNotFound = errors.New("pet not found")

// Repository defines the persistence contract for the pet domain.
// All write and read operations are scoped to ownerID to ensure users
// can only access their own pets.
type Repository interface {
	GetPaginated(ctx context.Context, ownerID string, page, perPage int) ([]models.Pet, int64, error)
	GetByID(ctx context.Context, id, ownerID string) (models.Pet, error)
	Create(ctx context.Context, ownerID string, payload CreatePetPayload) (models.Pet, error)
	Update(ctx context.Context, id, ownerID string, payload UpdatePetPayload) (models.Pet, error)
	Delete(ctx context.Context, id, ownerID string) error
	CountByOwner(ctx context.Context, ownerID string) (int64, error)
}
