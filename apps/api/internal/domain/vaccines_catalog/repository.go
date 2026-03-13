package vaccines_catalog

import (
	"context"
	"errors"

	"github.com/my-pets/api/internal/models"
)

var ErrNotFound = errors.New("vaccine catalog not found")

type Repository interface {
	GetPaginated(ctx context.Context, page, perPage int, speciesFilter *string) ([]models.VaccineCatalog, int64, error)
	GetByID(ctx context.Context, id string) (models.VaccineCatalog, error)
	GetBySpecies(ctx context.Context, species string) ([]models.VaccineCatalog, error)
	Create(ctx context.Context, payload CreateVaccineCatalogPayload) (models.VaccineCatalog, error)
	Update(ctx context.Context, id string, payload UpdateVaccineCatalogPayload) (models.VaccineCatalog, error)
	Delete(ctx context.Context, id string) error
}
