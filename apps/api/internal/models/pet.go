package models

import "time"

// Pet is the full entity stored in the database and returned by the API.
// GORM uses the struct tags to manage the schema via AutoMigrate.
type Pet struct {
	ID        string    `json:"id"         gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name      string    `json:"name"       gorm:"type:varchar(100);not null"`
	Species   string    `json:"species"    gorm:"type:varchar(50);not null"`
	Breed     string    `json:"breed"      gorm:"type:varchar(100)"`
	Age       int       `json:"age"        gorm:"default:0"`
	Owner     string    `json:"owner"      gorm:"type:varchar(100)"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
