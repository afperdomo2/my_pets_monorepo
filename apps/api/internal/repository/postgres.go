package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/my-pets/api/internal/models"
)

type postgresRepo struct {
	db *sql.DB
}

// NewPostgresRepo returns a PetRepository backed by PostgreSQL.
func NewPostgresRepo(db *sql.DB) PetRepository {
	return &postgresRepo{db: db}
}

// Migrate creates the pets table if it does not exist.
func Migrate(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS pets (
		id         SERIAL PRIMARY KEY,
		name       VARCHAR(100) NOT NULL,
		species    VARCHAR(50)  NOT NULL,
		breed      VARCHAR(100),
		age        INT          DEFAULT 0,
		owner      VARCHAR(100),
		created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
	);`
	_, err := db.Exec(query)
	return err
}

func (r *postgresRepo) GetAll() ([]models.Pet, error) {
	rows, err := r.db.Query(
		`SELECT id, name, species, breed, age, owner, created_at, updated_at FROM pets ORDER BY id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pets []models.Pet
	for rows.Next() {
		var p models.Pet
		if err := rows.Scan(&p.ID, &p.Name, &p.Species, &p.Breed, &p.Age, &p.Owner, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		pets = append(pets, p)
	}
	return pets, rows.Err()
}

func (r *postgresRepo) GetByID(id int) (models.Pet, error) {
	var p models.Pet
	err := r.db.QueryRow(
		`SELECT id, name, species, breed, age, owner, created_at, updated_at FROM pets WHERE id = $1`, id,
	).Scan(&p.ID, &p.Name, &p.Species, &p.Breed, &p.Age, &p.Owner, &p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return p, ErrNotFound
	}
	return p, err
}

func (r *postgresRepo) Create(pet models.Pet) (models.Pet, error) {
	err := r.db.QueryRow(
		`INSERT INTO pets (name, species, breed, age, owner)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, name, species, breed, age, owner, created_at, updated_at`,
		pet.Name, pet.Species, pet.Breed, pet.Age, pet.Owner,
	).Scan(&pet.ID, &pet.Name, &pet.Species, &pet.Breed, &pet.Age, &pet.Owner, &pet.CreatedAt, &pet.UpdatedAt)
	return pet, err
}

func (r *postgresRepo) Update(id int, pet models.Pet) (models.Pet, error) {
	pet.UpdatedAt = time.Now()
	err := r.db.QueryRow(
		`UPDATE pets
		 SET name=$1, species=$2, breed=$3, age=$4, owner=$5, updated_at=$6
		 WHERE id=$7
		 RETURNING id, name, species, breed, age, owner, created_at, updated_at`,
		pet.Name, pet.Species, pet.Breed, pet.Age, pet.Owner, pet.UpdatedAt, id,
	).Scan(&pet.ID, &pet.Name, &pet.Species, &pet.Breed, &pet.Age, &pet.Owner, &pet.CreatedAt, &pet.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return pet, ErrNotFound
	}
	return pet, err
}

func (r *postgresRepo) Delete(id int) error {
	res, err := r.db.Exec(`DELETE FROM pets WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

// ErrNotFound is returned when a pet does not exist in the database.
var ErrNotFound = errors.New("pet not found")
