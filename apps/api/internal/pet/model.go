package pet

import "time"

// PetPayload is the shape accepted from HTTP clients on create and update.
// It deliberately excludes server-managed fields (ID, timestamps).
type PetPayload struct {
	Name    string `json:"name"    binding:"required"`
	Species string `json:"species" binding:"required"`
	Breed   string `json:"breed"`
	Age     int    `json:"age"`
	Owner   string `json:"owner"`
}

// Pet is the full entity returned by the API and stored in the database.
// It embeds PetPayload so fields are not repeated.
type Pet struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	PetPayload
}
