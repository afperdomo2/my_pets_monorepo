package vaccine_application

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

// NewGormRepo crea un nuevo repositorio de vaccine_application respaldado por GORM.
func NewGormRepo(db *gorm.DB) Repository {
	return &gormRepo{db: db}
}

// parseDateString convierte una cadena "YYYY-MM-DD" en time.Time.
func parseDateString(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

func (r *gormRepo) GetByHealthRecordID(ctx context.Context, healthRecordID string) ([]models.VaccineApplication, error) {
	var apps []models.VaccineApplication
	result := r.db.WithContext(ctx).
		Where("health_record_id = ?", healthRecordID).
		Order("application_date ASC").
		Find(&apps)
	if result.Error != nil {
		return nil, fmt.Errorf("vaccine_application.GetByHealthRecordID: %w", result.Error)
	}
	return apps, nil
}

func (r *gormRepo) GetByID(ctx context.Context, id string) (models.VaccineApplication, error) {
	var app models.VaccineApplication
	result := r.db.WithContext(ctx).First(&app, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return models.VaccineApplication{}, ErrNotFound
	}
	if result.Error != nil {
		return models.VaccineApplication{}, fmt.Errorf("vaccine_application.GetByID: %w", result.Error)
	}
	return app, nil
}

func (r *gormRepo) Create(ctx context.Context, app models.VaccineApplication) (models.VaccineApplication, error) {
	result := r.db.WithContext(ctx).Create(&app)
	if result.Error != nil {
		return models.VaccineApplication{}, fmt.Errorf("vaccine_application.Create: %w", result.Error)
	}
	return app, nil
}

func (r *gormRepo) Update(ctx context.Context, app models.VaccineApplication) (models.VaccineApplication, error) {
	result := r.db.WithContext(ctx).Save(&app)
	if result.Error != nil {
		return models.VaccineApplication{}, fmt.Errorf("vaccine_application.Update: %w", result.Error)
	}
	return app, nil
}

func (r *gormRepo) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&models.VaccineApplication{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("vaccine_application.Delete: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
