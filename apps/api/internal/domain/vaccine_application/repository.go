package vaccine_application

import (
	"context"
	"errors"

	"github.com/my-pets/api/internal/models"
)

// ErrNotFound se retorna cuando la aplicación de vacuna no existe.
var ErrNotFound = errors.New("vaccine application not found")

// Repository define el contrato de persistencia para vaccine_application.
type Repository interface {
	// GetByHealthRecordID retorna todas las aplicaciones de un health_record.
	GetByHealthRecordID(ctx context.Context, healthRecordID string) ([]models.VaccineApplication, error)

	// GetByID retorna una aplicación por su ID.
	GetByID(ctx context.Context, id string) (models.VaccineApplication, error)

	// Create crea una nueva aplicación de vacuna.
	Create(ctx context.Context, app models.VaccineApplication) (models.VaccineApplication, error)

	// Update actualiza una aplicación existente.
	Update(ctx context.Context, app models.VaccineApplication) (models.VaccineApplication, error)

	// Delete elimina una aplicación.
	Delete(ctx context.Context, id string) error

	// UpdateHealthRecordAfterApplication actualiza last_dose_date y applied_doses_count en health_record.
	UpdateHealthRecordAfterApplication(ctx context.Context, healthRecordID string, applicationDate string) error
}
