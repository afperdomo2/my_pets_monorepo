package health_catalog

// CreateHealthCatalogPayload es la estructura aceptada al crear un nuevo registro en la guía de salud.
type CreateHealthCatalogPayload struct {
	Name            string   `json:"name"             binding:"required,min=1,max=100"`
	Category        string   `json:"category"         binding:"required,oneof=vaccine deworming exam"`
	Description     string   `json:"description"      binding:"max=1000"`
	Species         []string `json:"species"          binding:"required,min=1,dive,oneof=dog cat bird rabbit fish other"`
	FrequencyMonths *int     `json:"frequency_months" binding:"omitempty,min=1,max=120"`
	IsMandatory     bool     `json:"is_mandatory"`
}

// UpdateHealthCatalogPayload es la estructura aceptada al actualizar un registro de la guía de salud.
type UpdateHealthCatalogPayload struct {
	Name            string   `json:"name"             binding:"required,min=1,max=100"`
	Category        string   `json:"category"         binding:"required,oneof=vaccine deworming exam"`
	Description     string   `json:"description"      binding:"max=1000"`
	Species         []string `json:"species"          binding:"required,min=1,dive,oneof=dog cat bird rabbit fish other"`
	FrequencyMonths *int     `json:"frequency_months" binding:"omitempty,min=1,max=120"`
	IsMandatory     bool     `json:"is_mandatory"`
}
