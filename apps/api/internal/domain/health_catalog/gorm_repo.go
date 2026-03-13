package health_catalog

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/lib/pq"
	"github.com/my-pets/api/internal/models"
)

type gormRepo struct {
	db *gorm.DB
}

// NewGormRepo construye un Repository respaldado por GORM.
func NewGormRepo(db *gorm.DB) Repository {
	return &gormRepo{db: db}
}

// GetPaginatedByCategory retorna registros de la guía de salud paginados con filtro opcional de especie, filtrados por categoría.
func (r *gormRepo) GetPaginatedByCategory(ctx context.Context, category string, page, perPage int, speciesFilter *string) ([]models.HealthCatalog, int64, error) {
	var items []models.HealthCatalog
	var total int64

	base := r.db.WithContext(ctx).Model(&models.HealthCatalog{}).Where("category = ?", category)

	// Aplicar filtro de especie si se proporciona
	if speciesFilter != nil && *speciesFilter != "" {
		base = base.Where("? = ANY(species)", *speciesFilter)
	}

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("health_catalog.GetPaginatedByCategory count: %w", err)
	}

	offset := (page - 1) * perPage
	if err := base.Order("name ASC").Limit(perPage).Offset(offset).Find(&items).Error; err != nil {
		return nil, 0, fmt.Errorf("health_catalog.GetPaginatedByCategory: %w", err)
	}

	return items, total, nil
}

// GetByID retorna un registro de la guía de salud por su ID.
func (r *gormRepo) GetByID(ctx context.Context, id string) (models.HealthCatalog, error) {
	var item models.HealthCatalog
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&item)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return models.HealthCatalog{}, ErrNotFound
		}
		return models.HealthCatalog{}, fmt.Errorf("health_catalog.GetByID: %w", result.Error)
	}
	return item, nil
}

// GetBySpecies retorna todos los registros de la guía de salud aplicables a una especie.
func (r *gormRepo) GetBySpecies(ctx context.Context, species string) ([]models.HealthCatalog, error) {
	var items []models.HealthCatalog
	err := r.db.WithContext(ctx).
		Where("? = ANY(species)", species).
		Order("name ASC").
		Find(&items).Error
	if err != nil {
		return nil, fmt.Errorf("health_catalog.GetBySpecies: %w", err)
	}
	return items, nil
}

// Create persiste un nuevo registro en la guía de salud.
func (r *gormRepo) Create(ctx context.Context, payload CreateHealthCatalogPayload) (models.HealthCatalog, error) {
	item := models.HealthCatalog{
		Name:            payload.Name,
		Category:        payload.Category,
		Description:     payload.Description,
		Species:         pq.StringArray(payload.Species),
		FrequencyMonths: payload.FrequencyMonths,
		IsMandatory:     payload.IsMandatory,
	}

	if result := r.db.WithContext(ctx).Create(&item); result.Error != nil {
		return models.HealthCatalog{}, fmt.Errorf("health_catalog.Create: %w", result.Error)
	}
	return item, nil
}

// Update modifica un registro existente de la guía de salud.
func (r *gormRepo) Update(ctx context.Context, id string, payload UpdateHealthCatalogPayload) (models.HealthCatalog, error) {
	var item models.HealthCatalog
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&item)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return models.HealthCatalog{}, ErrNotFound
		}
		return models.HealthCatalog{}, fmt.Errorf("health_catalog.Update: %w", result.Error)
	}

	item.Name = payload.Name
	item.Category = payload.Category
	item.Description = payload.Description
	item.Species = pq.StringArray(payload.Species)
	item.FrequencyMonths = payload.FrequencyMonths
	item.IsMandatory = payload.IsMandatory

	if result := r.db.WithContext(ctx).Save(&item); result.Error != nil {
		return models.HealthCatalog{}, fmt.Errorf("health_catalog.Update: %w", result.Error)
	}
	return item, nil
}

// Delete elimina un registro de la guía de salud por su ID.
func (r *gormRepo) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.HealthCatalog{})
	if result.Error != nil {
		return fmt.Errorf("health_catalog.Delete: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
