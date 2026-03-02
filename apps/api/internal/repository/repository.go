package repository

import "github.com/my-pets/api/internal/models"

// PetRepository defines the contract for pet persistence.
// Swap implementations (postgres, in-memory, mock) without touching handlers.
type PetRepository interface {
	GetAll() ([]models.Pet, error)
	GetByID(id int) (models.Pet, error)
	Create(pet models.Pet) (models.Pet, error)
	Update(id int, pet models.Pet) (models.Pet, error)
	Delete(id int) error
}
