package health_record

import (
	"context"
	"errors"

	"github.com/my-pets/api/internal/models"
)

// ErrNotFound se retorna cuando el registro de salud no existe o no pertenece al usuario.
var ErrNotFound = errors.New("health record not found")

// Repository define el contrato de persistencia para el dominio health_record.
// Todas las operaciones están delimitadas por ownerID (el usuario autenticado),
// usando el campo user_id de la tabla para filtrado directo sin joins.
type Repository interface {
	// GetAllByOwner retorna todos los registros de salud del usuario, de todas sus mascotas.
	GetAllByOwner(ctx context.Context, ownerID string, page, perPage int) ([]models.HealthRecord, int64, error)

	// GetByPetID retorna los registros de salud de una mascota específica del usuario.
	GetByPetID(ctx context.Context, petID, ownerID string, page, perPage int) ([]models.HealthRecord, int64, error)

	// GetByPetIDAndCategory retorna los registros de salud de una mascota filtrados por categoría.
	GetByPetIDAndCategory(ctx context.Context, petID, category, ownerID string, page, perPage int) ([]models.HealthRecord, int64, error)

	// GetByID retorna un registro de salud por su ID, validando que pertenezca al usuario.
	GetByID(ctx context.Context, id, ownerID string) (models.HealthRecord, error)

	// Create crea un nuevo registro de salud. El pet_id viene en el payload.
	// Verifica que la mascota pertenezca al usuario antes de crear.
	Create(ctx context.Context, ownerID string, payload CreateHealthRecordPayload) (models.HealthRecord, error)

	// Update actualiza todos los campos editables de un registro. No permite cambiar pet_id.
	Update(ctx context.Context, id, ownerID string, payload UpdateHealthRecordPayload) (models.HealthRecord, error)

	// UpdateStatus actualiza únicamente el campo status de un registro.
	UpdateStatus(ctx context.Context, id, ownerID string, payload UpdateStatusPayload) (models.HealthRecord, error)

	// Delete elimina un registro de salud validando que pertenezca al usuario.
	Delete(ctx context.Context, id, ownerID string) error

	// GetUpcomingRecords retorna los próximos registros pendientes de aplicación del usuario.
	// Si category no está vacío, filtra por esa categoría.
	// Retorna máximo 'limit' registros, ordenados por due_date ASC.
	// No tiene paginación.
	GetUpcomingRecords(ctx context.Context, ownerID, category string, limit int) ([]models.HealthRecord, error)
}
