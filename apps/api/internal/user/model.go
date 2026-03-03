package user

import "time"

// AuthProvider indicates how the user authenticates.
type AuthProvider string

const (
	AuthProviderLocal  AuthProvider = "local"
	AuthProviderGoogle AuthProvider = "google"
)

// CreateUserPayload is the shape accepted from HTTP clients on user creation.
// Password is required because local accounts always need one.
type CreateUserPayload struct {
	Name     string `json:"name"     binding:"required,min=2,max=100"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

// UpdateUserPayload is the shape accepted from HTTP clients on user update.
// Password is optional: if omitted the existing hash is kept unchanged.
type UpdateUserPayload struct {
	Name  string `json:"name"  binding:"required,min=2,max=100"`
	Email string `json:"email" binding:"required,email"`
}

// GoogleUserInfo contains the profile data returned by Google's userinfo endpoint.
type GoogleUserInfo struct {
	GoogleID string
	Email    string
	Name     string
}

// User is the full entity returned by the API and stored in the database.
// Password is excluded from JSON responses to avoid leaking it.
// is_system_user marks the initial bootstrap account created during first-time setup.
type User struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Email        string       `json:"email"`
	IsSystemUser bool         `json:"is_system_user"`
	AuthProvider AuthProvider `json:"auth_provider"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}
