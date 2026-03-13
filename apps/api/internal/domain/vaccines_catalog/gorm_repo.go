package vaccines_catalog

import (
	"context"
	"fmt"

	"github.com/lib/pq"
	"gorm.io/gorm"

	"github.com/my-pets/api/internal/models"
)

type gormRepo struct {
	db *gorm.DB
}

func NewGormRepo(db *gorm.DB) Repository {
	return &gormRepo{db: db}
}

func (r *gormRepo) GetPaginated(ctx context.Context, page, perPage int) ([]models.VaccineCatalog, int64, error) {
	var vaccines []models.VaccineCatalog
	var total int64

	base := r.db.WithContext(ctx).Model(&models.VaccineCatalog{})

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("vaccines_catalog.GetPaginated count: %w", err)
	}

	offset := (page - 1) * perPage
	if err := base.Order("name ASC").Limit(perPage).Offset(offset).Find(&vaccines).Error; err != nil {
		return nil, 0, fmt.Errorf("vaccines_catalog.GetPaginated: %w", err)
	}

	return vaccines, total, nil
}

func (r *gormRepo) GetByID(ctx context.Context, id string) (models.VaccineCatalog, error) {
	var v models.VaccineCatalog
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&v)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return models.VaccineCatalog{}, ErrNotFound
		}
		return models.VaccineCatalog{}, fmt.Errorf("vaccines_catalog.GetByID: %w", result.Error)
	}
	return v, nil
}

func (r *gormRepo) GetBySpecies(ctx context.Context, species string) ([]models.VaccineCatalog, error) {
	var vaccines []models.VaccineCatalog
	err := r.db.WithContext(ctx).
		Where("? = ANY(species)", species).
		Order("name ASC").
		Find(&vaccines).Error
	if err != nil {
		return nil, fmt.Errorf("vaccines_catalog.GetBySpecies: %w", err)
	}
	return vaccines, nil
}

func (r *gormRepo) Create(ctx context.Context, payload CreateVaccineCatalogPayload) (models.VaccineCatalog, error) {
	v := models.VaccineCatalog{
		Name:            payload.Name,
		Species:         pq.StringArray(payload.Species),
		FrequencyMonths: payload.FrequencyMonths,
		IsMandatory:     payload.IsMandatory,
	}

	if result := r.db.WithContext(ctx).Create(&v); result.Error != nil {
		return models.VaccineCatalog{}, fmt.Errorf("vaccines_catalog.Create: %w", result.Error)
	}
	return v, nil
}

func (r *gormRepo) Update(ctx context.Context, id string, payload UpdateVaccineCatalogPayload) (models.VaccineCatalog, error) {
	var v models.VaccineCatalog
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&v)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return models.VaccineCatalog{}, ErrNotFound
		}
		return models.VaccineCatalog{}, fmt.Errorf("vaccines_catalog.Update: %w", result.Error)
	}

	v.Name = payload.Name
	v.Species = pq.StringArray(payload.Species)
	v.FrequencyMonths = payload.FrequencyMonths
	v.IsMandatory = payload.IsMandatory

	if result := r.db.WithContext(ctx).Save(&v); result.Error != nil {
		return models.VaccineCatalog{}, fmt.Errorf("vaccines_catalog.Update: %w", result.Error)
	}
	return v, nil
}

func (r *gormRepo) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.VaccineCatalog{})
	if result.Error != nil {
		return fmt.Errorf("vaccines_catalog.Delete: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
