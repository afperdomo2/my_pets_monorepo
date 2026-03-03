package pet

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type postgresRepo struct {
	db *sql.DB
}

// NewPostgresRepo returns a Repository backed by PostgreSQL.
func NewPostgresRepo(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

// Migrate creates the pets table if it does not exist.
func Migrate(ctx context.Context, db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS pets (
		id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name       VARCHAR(100) NOT NULL,
		species    VARCHAR(50)  NOT NULL,
		breed      VARCHAR(100),
		age        INT          DEFAULT 0,
		owner      VARCHAR(100),
		created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
	);`
	_, err := db.ExecContext(ctx, query)
	return err
}

const petColumns = `id, name, species, breed, age, owner, created_at, updated_at`

func scanPet(row interface {
	Scan(dest ...any) error
}) (Pet, error) {
	var p Pet
	err := row.Scan(
		&p.ID, &p.Name, &p.Species, &p.Breed,
		&p.Age, &p.Owner, &p.CreatedAt, &p.UpdatedAt,
	)
	return p, err
}

func (r *postgresRepo) GetAll(ctx context.Context) ([]Pet, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT `+petColumns+` FROM pets ORDER BY id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pets := []Pet{}
	for rows.Next() {
		p, err := scanPet(rows)
		if err != nil {
			return nil, err
		}
		pets = append(pets, p)
	}
	return pets, rows.Err()
}

func (r *postgresRepo) GetByID(ctx context.Context, id string) (Pet, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT `+petColumns+` FROM pets WHERE id = $1`, id,
	)
	p, err := scanPet(row)
	if errors.Is(err, sql.ErrNoRows) {
		return p, ErrNotFound
	}
	return p, err
}

func (r *postgresRepo) Create(ctx context.Context, payload PetPayload) (Pet, error) {
	newID := uuid.New().String()
	row := r.db.QueryRowContext(ctx,
		`INSERT INTO pets (id, name, species, breed, age, owner)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING `+petColumns,
		newID, payload.Name, payload.Species, payload.Breed, payload.Age, payload.Owner,
	)
	return scanPet(row)
}

func (r *postgresRepo) Update(ctx context.Context, id string, payload PetPayload) (Pet, error) {
	row := r.db.QueryRowContext(ctx,
		`UPDATE pets
		 SET name=$1, species=$2, breed=$3, age=$4, owner=$5, updated_at=$6
		 WHERE id=$7
		 RETURNING `+petColumns,
		payload.Name, payload.Species, payload.Breed, payload.Age, payload.Owner,
		time.Now(), id,
	)
	p, err := scanPet(row)
	if errors.Is(err, sql.ErrNoRows) {
		return p, ErrNotFound
	}
	return p, err
}

func (r *postgresRepo) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM pets WHERE id = $1`, id)
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
