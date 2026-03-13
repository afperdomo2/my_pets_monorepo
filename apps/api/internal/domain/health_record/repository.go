package health_record

import (
	"context"
	"errors"

	"github.com/my-pets/api/internal/models"
)

// ErrNotFound se retorna cuando el registro de salud no existe o no pertenece a la mascota del usuario.
var ErrNotFound = errors.New("health record not found")

// Repository define el contrato de persistencia para el dominio health_record.
// Todas las operaciones de escritura y lectura están delimitadas por petID y ownerID
// para garantizar que los usuarios solo accedan a registros de sus propias mascotas.
type Repository interface {
	// GetPaginated retorna una página de registros de salud de una mascota y el total.
	GetPaginated(ctx context.Context, petID, ownerID string, page, perPage int) ([]models.HealthRecord, int64, error)

	// GetByID retorna un registro de salud por su ID, validando que pertenezca a la mascota del usuario.
	GetByID(ctx context.Context, id, petID, ownerID string) (models.HealthRecord, error)

	// Create crea un nuevo registro de salud para la mascota indicada.
	Create(ctx context.Context, petID, ownerID string, payload CreateHealthRecordPayload) (models.HealthRecord, error)

	// Update actualiza todos los campos editables de un registro de salud existente.
	Update(ctx context.Context, id, petID, ownerID string, payload UpdateHealthRecordPayload) (models.HealthRecord, error)

	// UpdateStatus actualiza únicamente el campo status de un registro de salud.
	UpdateStatus(ctx context.Context, id, petID, ownerID string, payload UpdateStatusPayload) (models.HealthRecord, error)

	// Delete elimina un registro de salud de la base de datos.
	Delete(ctx context.Context, id, petID, ownerID string) error
}
