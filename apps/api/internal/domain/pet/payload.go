package pet

// CreatePetPayload is the shape accepted from HTTP clients on create.
// It includes fields that can only be set at creation time (WeightGrams, LifeStage).
type CreatePetPayload struct {
	Name           string  `json:"name"             binding:"required,min=1,max=100"`
	Species        string  `json:"species"          binding:"required,oneof=dog cat bird rabbit fish other"`
	Breed          string  `json:"breed"            binding:"omitempty,max=100"`
	BirthDate      string  `json:"birth_date"       binding:"required"`
	BirthDateExact bool    `json:"birth_date_exact"`
	WeightGrams    *int    `json:"weight_grams"     binding:"omitempty,min=1"`
	LifeStage      *string `json:"life_stage"       binding:"omitempty"`
	Size           *string `json:"size"             binding:"omitempty"`
}

// UpdatePetPayload is the shape accepted from HTTP clients on update.
// WeightGrams is intentionally excluded — weight is managed via a dedicated weighing log.
// LifeStage is NOT sent by the client; it is computed by the handler and stored here
// before being passed to the repository.
// Species is intentionally excluded — it cannot be changed after creation.
type UpdatePetPayload struct {
	Name           string  `json:"name"             binding:"required,min=1,max=100"`
	Breed          string  `json:"breed"            binding:"omitempty,max=100"`
	BirthDate      string  `json:"birth_date"       binding:"required"`
	BirthDateExact bool    `json:"birth_date_exact"`
	Size           *string `json:"size"             binding:"omitempty"`
	LifeStage      *string `json:"-"` // set by handler, not from client
}
