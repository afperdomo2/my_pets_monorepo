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

// parseDateString convierte una cadena "YYYY-MM-DD" en time.Time.
func parseDateString(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

// petBelongsToOwner verifica que la mascota existe y pertenece al usuario autenticado.
// Retorna el error ErrNotFound si la mascota no existe o no pertenece al usuario.
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

func (r *gormRepo) GetAllByOwner(ctx context.Context, ownerID string, page, perPage int) ([]models.HealthRecord, int64, error) {
	var records []models.HealthRecord
	var total int64

	base := r.db.WithContext(ctx).Model(&models.HealthRecord{}).Where("user_id = ?", ownerID)

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("health_record.GetAllByOwner count: %w", err)
	}

	offset := (page - 1) * perPage
	if err := base.Preload("Pet").Order("created_at DESC").Limit(perPage).Offset(offset).Find(&records).Error; err != nil {
		return nil, 0, fmt.Errorf("health_record.GetAllByOwner: %w", err)
	}

	return records, total, nil
}

func (r *gormRepo) GetByPetID(ctx context.Context, petID, ownerID string, page, perPage int) ([]models.HealthRecord, int64, error) {
	var records []models.HealthRecord
	var total int64

	base := r.db.WithContext(ctx).Model(&models.HealthRecord{}).
		Where("pet_id = ? AND user_id = ?", petID, ownerID)

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("health_record.GetByPetID count: %w", err)
	}

	offset := (page - 1) * perPage
	if err := base.Preload("Pet").Order("created_at DESC").Limit(perPage).Offset(offset).Find(&records).Error; err != nil {
		return nil, 0, fmt.Errorf("health_record.GetByPetID: %w", err)
	}

	return records, total, nil
}

func (r *gormRepo) GetByPetIDAndCategory(ctx context.Context, petID, category, ownerID string, page, perPage int) ([]models.HealthRecord, int64, error) {
	var records []models.HealthRecord
	var total int64

	base := r.db.WithContext(ctx).Model(&models.HealthRecord{}).
		Where("pet_id = ? AND category = ? AND user_id = ?", petID, category, ownerID)

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("health_record.GetByPetIDAndCategory count: %w", err)
	}

	offset := (page - 1) * perPage
	if err := base.Preload("Pet").Order("created_at DESC").Limit(perPage).Offset(offset).Find(&records).Error; err != nil {
		return nil, 0, fmt.Errorf("health_record.GetByPetIDAndCategory: %w", err)
	}

	return records, total, nil
}

func (r *gormRepo) GetByID(ctx context.Context, id, ownerID string) (models.HealthRecord, error) {
	var rec models.HealthRecord
	result := r.db.WithContext(ctx).
		Preload("VaccineApplications").
		Where("id = ? AND user_id = ?", id, ownerID).
		First(&rec)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return models.HealthRecord{}, ErrNotFound
	}
	if result.Error != nil {
		return models.HealthRecord{}, fmt.Errorf("health_record.GetByID: %w", result.Error)
	}
	return rec, nil
}

func (r *gormRepo) Create(ctx context.Context, ownerID string, payload CreateHealthRecordPayload) (models.HealthRecord, error) {
	// Verificar que la mascota pertenece al usuario antes de crear.
	if err := r.petBelongsToOwner(ctx, payload.PetID, ownerID); err != nil {
		return models.HealthRecord{}, err
	}

	// application_date es obligatorio
	t, err := parseDateString(payload.ApplicationDate)
	if err != nil {
		return models.HealthRecord{}, fmt.Errorf("health_record.Create: formato de application_date inválido: %w", err)
	}

	var nextDoseDate *time.Time
	if payload.NextDoseDate != nil {
		t, err := parseDateString(*payload.NextDoseDate)
		if err != nil {
			return models.HealthRecord{}, fmt.Errorf("health_record.Create: formato de next_dose_date inválido: %w", err)
		}
		nextDoseDate = &t
	}

	// Al crear el registro, la primera dosis ya está aplicada
	// Por lo tanto, last_dose_date = application_date y applied_doses_count = 1
	rec := models.HealthRecord{
		PetID:             payload.PetID,
		UserID:            ownerID,
		HealthCatalogID:   payload.HealthCatalogID,
		Category:          payload.Category,
		Name:              payload.Name,
		ApplicationDate:   t,
		NextDoseDate:      nextDoseDate,
		Notes:             payload.Notes,
		TotalDoses:        payload.TotalDoses,
		LastDoseDate:      &t,
		AppliedDosesCount: 1,
	}

	if result := r.db.WithContext(ctx).Create(&rec); result.Error != nil {
		return models.HealthRecord{}, fmt.Errorf("health_record.Create: %w", result.Error)
	}
	return rec, nil
}

func (r *gormRepo) Update(ctx context.Context, id, ownerID string, payload UpdateHealthRecordPayload) (models.HealthRecord, error) {
	var rec models.HealthRecord
	result := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, ownerID).First(&rec)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return models.HealthRecord{}, ErrNotFound
	}
	if result.Error != nil {
		return models.HealthRecord{}, fmt.Errorf("health_record.Update: %w", result.Error)
	}

	rec.Category = payload.Category
	rec.Name = payload.Name
	rec.Notes = payload.Notes

	// Actualizar total_doses (puede ser null para limpiar el valor)
	rec.TotalDoses = payload.TotalDoses

	// Actualizar application_date si se proporciona
	if payload.ApplicationDate != nil {
		t, err := parseDateString(*payload.ApplicationDate)
		if err != nil {
			return models.HealthRecord{}, fmt.Errorf("health_record.Update: formato de application_date inválido: %w", err)
		}
		rec.ApplicationDate = t
	}

	// Actualizar next_dose_date si se proporciona
	if payload.NextDoseDate != nil {
		t, err := parseDateString(*payload.NextDoseDate)
		if err != nil {
			return models.HealthRecord{}, fmt.Errorf("health_record.Update: formato de next_dose_date inválido: %w", err)
		}
		rec.NextDoseDate = &t
	}

	if result := r.db.WithContext(ctx).Save(&rec); result.Error != nil {
		return models.HealthRecord{}, fmt.Errorf("health_record.Update: %w", result.Error)
	}
	return rec, nil
}

func (r *gormRepo) Delete(ctx context.Context, id, ownerID string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("health_record_id = ?", id).
			Delete(&models.VaccineApplication{}).Error; err != nil {
			return fmt.Errorf("health_record.Delete vaccine_applications: %w", err)
		}
		result := tx.Where("id = ? AND user_id = ?", id, ownerID).
			Delete(&models.HealthRecord{})
		if result.Error != nil {
			return fmt.Errorf("health_record.Delete: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			return ErrNotFound
		}
		return nil
	})
}

func (r *gormRepo) GetUpcomingRecords(ctx context.Context, ownerID, category string, page, perPage int) ([]models.HealthRecord, int64, error) {
	var records []models.HealthRecord
	var total int64

	// Construir query base: registros con próxima dosis programada y sin aplicación
	query := r.db.WithContext(ctx).
		Model(&models.HealthRecord{}).
		Where("user_id = ? AND next_dose_date IS NOT NULL AND application_date IS NULL", ownerID)

	// Filtrar por categoría si se proporciona
	if category != "" {
		query = query.Where("category = ?", category)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("health_record.GetUpcomingRecords count: %w", err)
	}

	offset := (page - 1) * perPage
	if err := query.
		Preload("Pet").
		Order("next_dose_date ASC").
		Limit(perPage).
		Offset(offset).
		Find(&records).Error; err != nil {
		return nil, 0, fmt.Errorf("health_record.GetUpcomingRecords: %w", err)
	}

	return records, total, nil
}

func (r *gormRepo) UpdateLastDoseDate(ctx context.Context, healthRecordID string, lastDoseDate string) error {
	t, err := parseDateString(lastDoseDate)
	if err != nil {
		return fmt.Errorf("health_record.UpdateLastDoseDate: formato de fecha inválido: %w", err)
	}

	result := r.db.WithContext(ctx).
		Model(&models.HealthRecord{}).
		Where("id = ?", healthRecordID).
		Update("last_dose_date", t)

	if result.Error != nil {
		return fmt.Errorf("health_record.UpdateLastDoseDate: %w", result.Error)
	}

	return nil
}

func (r *gormRepo) IncrementAppliedDosesCount(ctx context.Context, healthRecordID string) error {
	result := r.db.WithContext(ctx).
		Model(&models.HealthRecord{}).
		Where("id = ?", healthRecordID).
		Update("applied_doses_count", gorm.Expr("applied_doses_count + 1"))

	if result.Error != nil {
		return fmt.Errorf("health_record.IncrementAppliedDosesCount: %w", result.Error)
	}

	return nil
}
