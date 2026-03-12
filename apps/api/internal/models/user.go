package models

import "time"

// AuthProvider indicates how the user authenticates.
type AuthProvider string

const (
	AuthProviderLocal  AuthProvider = "local"
	AuthProviderGoogle AuthProvider = "google"
)

// User is the full entity stored in the database and returned by the API.
// Password and GoogleID are excluded from JSON responses.
// IsSystemUser marks the initial bootstrap account created during first-time setup.
type User struct {
	ID           string       `json:"id"             gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name         string       `json:"name"           gorm:"type:text;not null"`
	Email        string       `json:"email"          gorm:"type:text;uniqueIndex:users_email_key;not null"`
	Password     *string      `json:"-"              gorm:"type:varchar(255)"`
	IsSystemUser bool         `json:"is_system_user" gorm:"not null;default:false"`
	AuthProvider AuthProvider `json:"auth_provider"  gorm:"type:varchar(20);not null;default:'local'"`
	GoogleID     *string      `json:"-"              gorm:"type:varchar(255);uniqueIndex:users_google_id_key"`
	PetLimit     int          `json:"pet_limit"      gorm:"not null;default:5"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}
