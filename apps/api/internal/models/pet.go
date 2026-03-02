package models

import "time"

type Pet struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name" binding:"required"`
	Species   string    `json:"species" binding:"required"` // dog, cat, bird, etc.
	Breed     string    `json:"breed"`
	Age       int       `json:"age"`
	Owner     string    `json:"owner"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
