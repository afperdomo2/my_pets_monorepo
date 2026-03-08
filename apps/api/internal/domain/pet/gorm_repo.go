package pet

import (
	"context"
	"errors"
	"fmt"

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

func (r *gormRepo) GetAll(ctx context.Context) ([]models.Pet, error) {
	var pets []models.Pet
	if result := r.db.WithContext(ctx).Order("id").Find(&pets); result.Error != nil {
		return nil, fmt.Errorf("pet.GetAll: %w", result.Error)
	}
	return pets, nil
}

func (r *gormRepo) GetByID(ctx context.Context, id string) (models.Pet, error) {
	var p models.Pet
	result := r.db.WithContext(ctx).First(&p, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return models.Pet{}, ErrNotFound
	}
	if result.Error != nil {
		return models.Pet{}, fmt.Errorf("pet.GetByID: %w", result.Error)
	}
	return p, nil
}

func (r *gormRepo) Create(ctx context.Context, payload PetPayload) (models.Pet, error) {
	p := models.Pet{
		Name:    payload.Name,
		Species: payload.Species,
		Breed:   payload.Breed,
		Age:     payload.Age,
		Owner:   payload.Owner,
	}
	if result := r.db.WithContext(ctx).Create(&p); result.Error != nil {
		return models.Pet{}, fmt.Errorf("pet.Create: %w", result.Error)
	}
	return p, nil
}

func (r *gormRepo) Update(ctx context.Context, id string, payload PetPayload) (models.Pet, error) {
	var p models.Pet
	result := r.db.WithContext(ctx).First(&p, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return models.Pet{}, ErrNotFound
	}
	if result.Error != nil {
		return models.Pet{}, fmt.Errorf("pet.Update: %w", result.Error)
	}

	p.Name = payload.Name
	p.Species = payload.Species
	p.Breed = payload.Breed
	p.Age = payload.Age
	p.Owner = payload.Owner

	if result := r.db.WithContext(ctx).Save(&p); result.Error != nil {
		return models.Pet{}, fmt.Errorf("pet.Update: %w", result.Error)
	}
	return p, nil
}

func (r *gormRepo) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&models.Pet{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("pet.Delete: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
