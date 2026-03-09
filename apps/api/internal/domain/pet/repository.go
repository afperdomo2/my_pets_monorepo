package pet

import (
	"context"
	"errors"

	"github.com/my-pets/api/internal/models"
)

// ErrNotFound is returned when a pet does not exist in the database.
var ErrNotFound = errors.New("pet not found")

// Repository defines the persistence contract for the pet domain.
// All methods accept a context to support cancellation and timeouts.
type Repository interface {
	GetAll(ctx context.Context) ([]models.Pet, error)
	// GetPaginated returns a page of pets and the total count.
	GetPaginated(ctx context.Context, page, perPage int) ([]models.Pet, int64, error)
	GetByID(ctx context.Context, id string) (models.Pet, error)
	Create(ctx context.Context, payload PetPayload) (models.Pet, error)
	Update(ctx context.Context, id string, payload PetPayload) (models.Pet, error)
	Delete(ctx context.Context, id string) error
}
