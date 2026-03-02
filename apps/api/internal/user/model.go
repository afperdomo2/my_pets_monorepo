package user

import "time"

// UserPayload is the shape accepted from HTTP clients on create and update.
// It deliberately excludes server-managed fields (ID, timestamps, is_system_user).
type UserPayload struct {
	Name     string `json:"name"     binding:"required"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// User is the full entity returned by the API and stored in the database.
// Password is excluded from JSON responses to avoid leaking it.
// is_system_user marks the initial bootstrap account created during first-time setup.
type User struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	IsSystemUser bool      `json:"is_system_user"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
