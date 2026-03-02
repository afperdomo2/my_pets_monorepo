package user

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type postgresRepo struct {
	db *sql.DB
}

// NewPostgresRepo returns a Repository backed by PostgreSQL.
func NewPostgresRepo(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

// Migrate creates the users table if it does not exist.
func Migrate(ctx context.Context, db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id             SERIAL PRIMARY KEY,
		name           VARCHAR(100)        NOT NULL,
		email          VARCHAR(255) UNIQUE NOT NULL,
		password       VARCHAR(255)        NOT NULL,
		is_system_user BOOLEAN             NOT NULL DEFAULT FALSE,
		created_at     TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
		updated_at     TIMESTAMPTZ         NOT NULL DEFAULT NOW()
	);`
	_, err := db.ExecContext(ctx, query)
	return err
}

const userColumns = `id, name, email, is_system_user, created_at, updated_at`

func scanUser(row interface {
	Scan(dest ...any) error
}) (User, error) {
	var u User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.IsSystemUser, &u.CreatedAt, &u.UpdatedAt)
	return u, err
}

func (r *postgresRepo) GetAll(ctx context.Context) ([]User, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT `+userColumns+` FROM users ORDER BY id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *postgresRepo) GetByID(ctx context.Context, id int) (User, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT `+userColumns+` FROM users WHERE id = $1`, id,
	)
	u, err := scanUser(row)
	if errors.Is(err, sql.ErrNoRows) {
		return u, ErrNotFound
	}
	return u, err
}

func (r *postgresRepo) GetByEmail(ctx context.Context, email string) (User, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT `+userColumns+` FROM users WHERE email = $1`, email,
	)
	u, err := scanUser(row)
	if errors.Is(err, sql.ErrNoRows) {
		return u, ErrNotFound
	}
	return u, err
}

func (r *postgresRepo) Create(ctx context.Context, payload UserPayload, isSystemUser bool) (User, error) {
	row := r.db.QueryRowContext(ctx,
		`INSERT INTO users (name, email, password, is_system_user)
		 VALUES ($1, $2, $3, $4)
		 RETURNING `+userColumns,
		payload.Name, payload.Email, payload.Password, isSystemUser,
	)
	u, err := scanUser(row)
	if err != nil {
		// Unique constraint violation on email (pgx error code 23505)
		if isUniqueViolation(err) {
			return u, ErrEmailTaken
		}
		return u, err
	}
	return u, nil
}

func (r *postgresRepo) Update(ctx context.Context, id int, payload UserPayload) (User, error) {
	row := r.db.QueryRowContext(ctx,
		`UPDATE users
		 SET name=$1, email=$2, password=$3, updated_at=$4
		 WHERE id=$5
		 RETURNING `+userColumns,
		payload.Name, payload.Email, payload.Password, time.Now(), id,
	)
	u, err := scanUser(row)
	if errors.Is(err, sql.ErrNoRows) {
		return u, ErrNotFound
	}
	if err != nil {
		if isUniqueViolation(err) {
			return u, ErrEmailTaken
		}
		return u, err
	}
	return u, nil
}

func (r *postgresRepo) Delete(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
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

func (r *postgresRepo) SystemUserExists(ctx context.Context) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM users WHERE is_system_user = TRUE)`,
	).Scan(&exists)
	return exists, err
}

// sqlstateErr is satisfied by pgx error types that expose the SQLSTATE code.
type sqlstateErr interface {
	SQLState() string
}

// isUniqueViolation reports whether err is a PostgreSQL unique-constraint
// violation (SQLSTATE 23505). pgx/v5 errors implement SQLState().
func isUniqueViolation(err error) bool {
	var se sqlstateErr
	return errors.As(err, &se) && se.SQLState() == "23505"
}
