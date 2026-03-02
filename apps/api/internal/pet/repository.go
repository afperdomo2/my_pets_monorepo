package pet

import (
	"context"
	"errors"
)

// ErrNotFound is returned when a pet does not exist in the database.
var ErrNotFound = errors.New("pet not found")

// Repository defines the persistence contract for the pet domain.
// All methods accept a context to support cancellation and timeouts.
type Repository interface {
	GetAll(ctx context.Context) ([]Pet, error)
	GetByID(ctx context.Context, id int) (Pet, error)
	Create(ctx context.Context, payload PetPayload) (Pet, error)
	Update(ctx context.Context, id int, payload PetPayload) (Pet, error)
	Delete(ctx context.Context, id int) error
}
