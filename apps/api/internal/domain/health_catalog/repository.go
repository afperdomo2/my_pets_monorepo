package health_catalog

import (
	"context"
	"errors"

	"github.com/my-pets/api/internal/models"
)

// ErrNotFound se retorna cuando un registro de la guía de salud no existe.
var ErrNotFound = errors.New("health catalog not found")

// Repository define las operaciones de persistencia para la guía de salud.
type Repository interface {
	GetPaginatedByCategory(ctx context.Context, category string, page, perPage int, speciesFilter *string) ([]models.HealthCatalog, int64, error)
	GetByID(ctx context.Context, id string) (models.HealthCatalog, error)
	GetBySpecies(ctx context.Context, species string) ([]models.HealthCatalog, error)
	Create(ctx context.Context, payload CreateHealthCatalogPayload) (models.HealthCatalog, error)
	Update(ctx context.Context, id string, payload UpdateHealthCatalogPayload) (models.HealthCatalog, error)
	Delete(ctx context.Context, id string) error
}
