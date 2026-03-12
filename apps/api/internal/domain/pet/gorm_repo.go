package pet

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/my-pets/api/internal/models"
)

// gormRepo implements Repository using GORM.
type gormRepo struct {
	db *gorm.DB
}

// NewGormRepo creates a new GORM-backed pet repository.
func NewGormRepo(db *gorm.DB) Repository {
	return &gormRepo{db: db}
}

func (r *gormRepo) GetPaginated(ctx context.Context, ownerID string, page, perPage int) ([]models.Pet, int64, error) {
	var pets []models.Pet
	var total int64

	base := r.db.WithContext(ctx).Model(&models.Pet{}).Where("owner_id = ?", ownerID)

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("pet.GetPaginated count: %w", err)
	}

	offset := (page - 1) * perPage
	if err := base.Order("created_at DESC").Limit(perPage).Offset(offset).Find(&pets).Error; err != nil {
		return nil, 0, fmt.Errorf("pet.GetPaginated: %w", err)
	}

	return pets, total, nil
}

func (r *gormRepo) GetByID(ctx context.Context, id, ownerID string) (models.Pet, error) {
	var p models.Pet
	result := r.db.WithContext(ctx).Where("id = ? AND owner_id = ?", id, ownerID).First(&p)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return models.Pet{}, ErrNotFound
	}
	if result.Error != nil {
		return models.Pet{}, fmt.Errorf("pet.GetByID: %w", result.Error)
	}
	return p, nil
}

func (r *gormRepo) Create(ctx context.Context, ownerID string, payload CreatePetPayload) (models.Pet, error) {
	birthDate, err := time.Parse("2006-01-02", payload.BirthDate)
	if err != nil {
		return models.Pet{}, fmt.Errorf("pet.Create: invalid birth_date format: %w", err)
	}

	p := models.Pet{
		OwnerID:        ownerID,
		Name:           payload.Name,
		Species:        payload.Species,
		Breed:          payload.Breed,
		BirthDate:      birthDate,
		BirthDateExact: payload.BirthDateExact,
		WeightGrams:    payload.WeightGrams,
		LifeStage:      payload.LifeStage,
	}

	if result := r.db.WithContext(ctx).Create(&p); result.Error != nil {
		return models.Pet{}, fmt.Errorf("pet.Create: %w", result.Error)
	}
	return p, nil
}

func (r *gormRepo) Update(ctx context.Context, id, ownerID string, payload UpdatePetPayload) (models.Pet, error) {
	var p models.Pet
	result := r.db.WithContext(ctx).Where("id = ? AND owner_id = ?", id, ownerID).First(&p)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return models.Pet{}, ErrNotFound
	}
	if result.Error != nil {
		return models.Pet{}, fmt.Errorf("pet.Update: %w", result.Error)
	}

	birthDate, err := time.Parse("2006-01-02", payload.BirthDate)
	if err != nil {
		return models.Pet{}, fmt.Errorf("pet.Update: invalid birth_date format: %w", err)
	}

	p.Name = payload.Name
	p.Species = payload.Species
	p.Breed = payload.Breed
	p.BirthDate = birthDate
	p.BirthDateExact = payload.BirthDateExact
	// WeightGrams and LifeStage are intentionally not updated here.

	if result := r.db.WithContext(ctx).Save(&p); result.Error != nil {
		return models.Pet{}, fmt.Errorf("pet.Update: %w", result.Error)
	}
	return p, nil
}

func (r *gormRepo) Delete(ctx context.Context, id, ownerID string) error {
	result := r.db.WithContext(ctx).Where("id = ? AND owner_id = ?", id, ownerID).Delete(&models.Pet{})
	if result.Error != nil {
		return fmt.Errorf("pet.Delete: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *gormRepo) CountByOwner(ctx context.Context, ownerID string) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Pet{}).Where("owner_id = ?", ownerID).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("pet.CountByOwner: %w", err)
	}
	return count, nil
}
