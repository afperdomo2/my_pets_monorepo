package health_record

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/my-pets/api/internal/models"
)

// gormRepo implementa Repository usando GORM.
type gormRepo struct {
	db *gorm.DB
}

// NewGormRepo crea un nuevo repositorio de health_record respaldado por GORM.
func NewGormRepo(db *gorm.DB) Repository {
	return &gormRepo{db: db}
}

// petBelongsToOwner verifica que la mascota existe y pertenece al usuario autenticado.
// Retorna error si la mascota no existe o no pertenece al usuario.
func (r *gormRepo) petBelongsToOwner(ctx context.Context, petID, ownerID string) error {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.Pet{}).
		Where("id = ? AND user_id = ?", petID, ownerID).
		Count(&count).Error; err != nil {
		return fmt.Errorf("health_record: verificar mascota: %w", err)
	}
	if count == 0 {
		return ErrNotFound
	}
	return nil
}

// parseDateString convierte una cadena "YYYY-MM-DD" en time.Time.
func parseDateString(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

func (r *gormRepo) GetPaginated(ctx context.Context, petID, ownerID string, page, perPage int) ([]models.HealthRecord, int64, error) {
	// Verificar que la mascota pertenece al usuario antes de listar.
	if err := r.petBelongsToOwner(ctx, petID, ownerID); err != nil {
		return nil, 0, err
	}

	var records []models.HealthRecord
	var total int64

	base := r.db.WithContext(ctx).Model(&models.HealthRecord{}).Where("pet_id = ?", petID)

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("health_record.GetPaginated count: %w", err)
	}

	offset := (page - 1) * perPage
	if err := base.Order("due_date ASC").Limit(perPage).Offset(offset).Find(&records).Error; err != nil {
		return nil, 0, fmt.Errorf("health_record.GetPaginated: %w", err)
	}

	return records, total, nil
}

func (r *gormRepo) GetByID(ctx context.Context, id, petID, ownerID string) (models.HealthRecord, error) {
	// Verificar que la mascota pertenece al usuario antes de leer el registro.
	if err := r.petBelongsToOwner(ctx, petID, ownerID); err != nil {
		return models.HealthRecord{}, err
	}

	var rec models.HealthRecord
	result := r.db.WithContext(ctx).
		Where("id = ? AND pet_id = ?", id, petID).
		First(&rec)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return models.HealthRecord{}, ErrNotFound
	}
	if result.Error != nil {
		return models.HealthRecord{}, fmt.Errorf("health_record.GetByID: %w", result.Error)
	}
	return rec, nil
}

func (r *gormRepo) Create(ctx context.Context, petID, ownerID string, payload CreateHealthRecordPayload) (models.HealthRecord, error) {
	// Verificar que la mascota pertenece al usuario antes de crear.
	if err := r.petBelongsToOwner(ctx, petID, ownerID); err != nil {
		return models.HealthRecord{}, err
	}

	dueDate, err := parseDateString(payload.DueDate)
	if err != nil {
		return models.HealthRecord{}, fmt.Errorf("health_record.Create: formato de due_date inválido: %w", err)
	}

	var applicationDate *time.Time
	if payload.ApplicationDate != nil {
		t, err := parseDateString(*payload.ApplicationDate)
		if err != nil {
			return models.HealthRecord{}, fmt.Errorf("health_record.Create: formato de application_date inválido: %w", err)
		}
		applicationDate = &t
	}

	status := "pending"
	if payload.Status != "" {
		status = payload.Status
	}

	rec := models.HealthRecord{
		PetID:           petID,
		HealthCatalogID: payload.HealthCatalogID,
		Category:        payload.Category,
		Name:            payload.Name,
		Status:          status,
		ApplicationDate: applicationDate,
		DueDate:         dueDate,
		Notes:           payload.Notes,
	}

	if result := r.db.WithContext(ctx).Create(&rec); result.Error != nil {
		return models.HealthRecord{}, fmt.Errorf("health_record.Create: %w", result.Error)
	}
	return rec, nil
}

func (r *gormRepo) Update(ctx context.Context, id, petID, ownerID string, payload UpdateHealthRecordPayload) (models.HealthRecord, error) {
	// Verificar que la mascota pertenece al usuario antes de actualizar.
	if err := r.petBelongsToOwner(ctx, petID, ownerID); err != nil {
		return models.HealthRecord{}, err
	}

	var rec models.HealthRecord
	result := r.db.WithContext(ctx).Where("id = ? AND pet_id = ?", id, petID).First(&rec)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return models.HealthRecord{}, ErrNotFound
	}
	if result.Error != nil {
		return models.HealthRecord{}, fmt.Errorf("health_record.Update: %w", result.Error)
	}

	dueDate, err := parseDateString(payload.DueDate)
	if err != nil {
		return models.HealthRecord{}, fmt.Errorf("health_record.Update: formato de due_date inválido: %w", err)
	}

	var applicationDate *time.Time
	if payload.ApplicationDate != nil {
		t, err := parseDateString(*payload.ApplicationDate)
		if err != nil {
			return models.HealthRecord{}, fmt.Errorf("health_record.Update: formato de application_date inválido: %w", err)
		}
		applicationDate = &t
	}

	rec.Category = payload.Category
	rec.Name = payload.Name
	rec.DueDate = dueDate
	rec.ApplicationDate = applicationDate
	rec.Notes = payload.Notes

	if payload.Status != "" {
		rec.Status = payload.Status
	}

	if result := r.db.WithContext(ctx).Save(&rec); result.Error != nil {
		return models.HealthRecord{}, fmt.Errorf("health_record.Update: %w", result.Error)
	}
	return rec, nil
}

func (r *gormRepo) UpdateStatus(ctx context.Context, id, petID, ownerID string, payload UpdateStatusPayload) (models.HealthRecord, error) {
	// Verificar que la mascota pertenece al usuario antes de actualizar el status.
	if err := r.petBelongsToOwner(ctx, petID, ownerID); err != nil {
		return models.HealthRecord{}, err
	}

	var rec models.HealthRecord
	result := r.db.WithContext(ctx).Where("id = ? AND pet_id = ?", id, petID).First(&rec)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return models.HealthRecord{}, ErrNotFound
	}
	if result.Error != nil {
		return models.HealthRecord{}, fmt.Errorf("health_record.UpdateStatus: %w", result.Error)
	}

	rec.Status = payload.Status

	// Si se marca como 'applied' y se provee application_date, actualizarla.
	if payload.ApplicationDate != nil {
		t, err := parseDateString(*payload.ApplicationDate)
		if err != nil {
			return models.HealthRecord{}, fmt.Errorf("health_record.UpdateStatus: formato de application_date inválido: %w", err)
		}
		rec.ApplicationDate = &t
	}

	if result := r.db.WithContext(ctx).Save(&rec); result.Error != nil {
		return models.HealthRecord{}, fmt.Errorf("health_record.UpdateStatus: %w", result.Error)
	}
	return rec, nil
}

func (r *gormRepo) Delete(ctx context.Context, id, petID, ownerID string) error {
	// Verificar que la mascota pertenece al usuario antes de eliminar.
	if err := r.petBelongsToOwner(ctx, petID, ownerID); err != nil {
		return err
	}

	result := r.db.WithContext(ctx).
		Where("id = ? AND pet_id = ?", id, petID).
		Delete(&models.HealthRecord{})
	if result.Error != nil {
		return fmt.Errorf("health_record.Delete: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
