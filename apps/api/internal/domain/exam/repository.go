package exam

import (
	"context"
	"errors"

	"github.com/my-pets/api/internal/models"
)

// ErrNotFound se retorna cuando el examen no existe.
var ErrNotFound = errors.New("exam not found")

// Repository define el contrato de persistencia para exam.
type Repository interface {
	// GetAllByOwner retorna todos los exámenes del usuario con paginación.
	GetAllByOwner(ctx context.Context, ownerID string, page, perPage int) ([]models.Exam, int64, error)

	// GetByPetID retorna los exámenes de una mascota específica con paginación.
	GetByPetID(ctx context.Context, petID, ownerID string, page, perPage int) ([]models.Exam, int64, error)

	// GetByID retorna un examen por su ID, validando que pertenezca al usuario.
	GetByID(ctx context.Context, id, ownerID string) (models.Exam, error)

	// GetByIDWithResults retorna un examen con sus resultados.
	GetByIDWithResults(ctx context.Context, id, ownerID string) (models.Exam, []models.ExamResult, error)

	// Create crea un nuevo examen.
	Create(ctx context.Context, exam models.Exam) (models.Exam, error)

	// Update actualiza un examen existente.
	Update(ctx context.Context, exam models.Exam) (models.Exam, error)

	// Delete elimina un examen.
	Delete(ctx context.Context, id, ownerID string) error

	// CreateResults crea múltiples resultados para un examen.
	CreateResults(ctx context.Context, results []models.ExamResult) error

	// DeleteResultsByExamID elimina todos los resultados de un examen.
	DeleteResultsByExamID(ctx context.Context, examID string) error

	// GetResultsByExamID retorna los resultados de un examen.
	GetResultsByExamID(ctx context.Context, examID string) ([]models.ExamResult, error)
}
