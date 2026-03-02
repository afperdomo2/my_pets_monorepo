package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type postgresRepo struct {
	db *sql.DB
}

// NewPostgresRepo returns a Repository backed by PostgreSQL.
func NewPostgresRepo(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

// Migrate creates the users table if it does not exist and applies column additions.
func Migrate(ctx context.Context, db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id             SERIAL PRIMARY KEY,
			name           VARCHAR(100)        NOT NULL,
			email          VARCHAR(255) UNIQUE NOT NULL,
			password       VARCHAR(255),
			is_system_user BOOLEAN             NOT NULL DEFAULT FALSE,
			auth_provider  VARCHAR(20)         NOT NULL DEFAULT 'local',
			google_id      VARCHAR(255) UNIQUE,
			created_at     TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
			updated_at     TIMESTAMPTZ         NOT NULL DEFAULT NOW()
		)`,
		// Idempotent additions for existing tables created before this migration.
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS auth_provider VARCHAR(20) NOT NULL DEFAULT 'local'`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS google_id VARCHAR(255)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS users_google_id_idx ON users (google_id) WHERE google_id IS NOT NULL`,
		// password can now be NULL for Google-only accounts
		`ALTER TABLE users ALTER COLUMN password DROP NOT NULL`,
	}
	for _, q := range queries {
		if _, err := db.ExecContext(ctx, q); err != nil {
			return err
		}
	}
	return nil
}

// userColumns are the safe columns returned in API responses (no password).
const userColumns = `id, name, email, is_system_user, auth_provider, created_at, updated_at`

func scanUser(row interface {
	Scan(dest ...any) error
}) (User, error) {
	var u User
	err := row.Scan(
		&u.ID, &u.Name, &u.Email,
		&u.IsSystemUser, &u.AuthProvider,
		&u.CreatedAt, &u.UpdatedAt,
	)
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

// GetPasswordByEmail returns the bcrypt hash stored for the given email.
// Used only during login — the hash is never sent to clients.
func (r *postgresRepo) GetPasswordByEmail(ctx context.Context, email string) (string, error) {
	var hash sql.NullString
	err := r.db.QueryRowContext(ctx,
		`SELECT password FROM users WHERE email = $1`, email,
	).Scan(&hash)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", err
	}
	if !hash.Valid || hash.String == "" {
		// Account exists but has no password (Google-only account).
		return "", ErrWrongProvider
	}
	return hash.String, nil
}

func (r *postgresRepo) Create(ctx context.Context, payload CreateUserPayload, isSystemUser bool) (User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	row := r.db.QueryRowContext(ctx,
		`INSERT INTO users (name, email, password, is_system_user, auth_provider)
		 VALUES ($1, $2, $3, $4, 'local')
		 RETURNING `+userColumns,
		payload.Name, payload.Email, string(hash), isSystemUser,
	)
	u, err := scanUser(row)
	if err != nil {
		if isUniqueViolation(err) {
			return u, ErrEmailTaken
		}
		return u, err
	}
	return u, nil
}

func (r *postgresRepo) Update(ctx context.Context, id int, payload UpdateUserPayload) (User, error) {
	var row *sql.Row
	if payload.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
		if err != nil {
			return User{}, err
		}
		row = r.db.QueryRowContext(ctx,
			`UPDATE users
			 SET name=$1, email=$2, password=$3, updated_at=$4
			 WHERE id=$5
			 RETURNING `+userColumns,
			payload.Name, payload.Email, string(hash), time.Now(), id,
		)
	} else {
		row = r.db.QueryRowContext(ctx,
			`UPDATE users
			 SET name=$1, email=$2, updated_at=$3
			 WHERE id=$4
			 RETURNING `+userColumns,
			payload.Name, payload.Email, time.Now(), id,
		)
	}
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

func (r *postgresRepo) HasUsers(ctx context.Context) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM users LIMIT 1)`,
	).Scan(&exists)
	return exists, err
}

func (r *postgresRepo) SystemUserExists(ctx context.Context) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM users WHERE is_system_user = TRUE)`,
	).Scan(&exists)
	return exists, err
}

// UpsertByGoogleID creates a new user or returns the existing one for a given Google account.
// If the email already exists with a local provider, ErrWrongProvider is returned.
func (r *postgresRepo) UpsertByGoogleID(ctx context.Context, info GoogleUserInfo, isSystemUser bool) (User, error) {
	row := r.db.QueryRowContext(ctx,
		`INSERT INTO users (name, email, google_id, is_system_user, auth_provider)
		 VALUES ($1, $2, $3, $4, 'google')
		 ON CONFLICT (google_id) DO UPDATE
		   SET name = EXCLUDED.name,
		       updated_at = NOW()
		 RETURNING `+userColumns,
		info.Name, info.Email, info.GoogleID, isSystemUser,
	)
	u, err := scanUser(row)
	if err != nil {
		if isUniqueViolation(err) {
			// email unique constraint — account exists with local provider
			return u, ErrWrongProvider
		}
		return u, err
	}
	return u, nil
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
