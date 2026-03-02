package pet

import "time"

// PetPayload is the shape accepted from HTTP clients on create and update.
// It deliberately excludes server-managed fields (ID, timestamps).
type PetPayload struct {
	Name    string `json:"name"    binding:"required,min=1,max=100"`
	Species string `json:"species" binding:"required,oneof=dog cat bird rabbit fish other"`
	Breed   string `json:"breed"   binding:"omitempty,max=100"`
	Age     int    `json:"age"     binding:"min=0,max=100"`
	Owner   string `json:"owner"   binding:"omitempty,max=100"`
}

// Pet is the full entity returned by the API and stored in the database.
// It embeds PetPayload so fields are not repeated.
type Pet struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	PetPayload
}
