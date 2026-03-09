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
}

// UpdatePetPayload is the shape accepted from HTTP clients on update.
// WeightGrams and LifeStage are intentionally excluded — weight is managed
// via a dedicated weighing log, and life stage is derived from weight.
type UpdatePetPayload struct {
	Name           string `json:"name"             binding:"required,min=1,max=100"`
	Species        string `json:"species"          binding:"required,oneof=dog cat bird rabbit fish other"`
	Breed          string `json:"breed"            binding:"omitempty,max=100"`
	BirthDate      string `json:"birth_date"       binding:"required"`
	BirthDateExact bool   `json:"birth_date_exact"`
}
