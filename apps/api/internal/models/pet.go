package models

import "time"

// Pet is the full entity stored in the database and returned by the API.
// GORM uses the struct tags to manage the schema via AutoMigrate.
type Pet struct {
	ID             string    `json:"id"               gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID         string    `json:"user_id"         gorm:"type:uuid;not null;index"`
	Name           string    `json:"name"             gorm:"type:varchar(100);not null"`
	Species        string    `json:"species"          gorm:"type:varchar(50);not null"`
	Breed          string    `json:"breed"            gorm:"type:varchar(100)"`
	BirthDate      time.Time `json:"birth_date"       gorm:"not null"`
	BirthDateExact bool      `json:"birth_date_exact" gorm:"not null;default:false"`
	WeightGrams    *int      `json:"weight_grams"     gorm:"default:null"`
	LifeStage      *string   `json:"life_stage"       gorm:"type:varchar(20);default:null"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
