package vaccines_catalog

// CreateVaccineCatalogPayload es la estructura aceptada al crear una nueva vaccine.
type CreateVaccineCatalogPayload struct {
	Name            string   `json:"name"             binding:"required,min=1,max=100"`
	Species         []string `json:"species"          binding:"required,min=1,dive,oneof=dog cat bird rabbit fish other"`
	FrequencyMonths int      `json:"frequency_months" binding:"required,min=1,max=120"`
	IsMandatory     bool     `json:"is_mandatory"`
}

// UpdateVaccineCatalogPayload es la estructura aceptada al actualizar una vaccine.
type UpdateVaccineCatalogPayload struct {
	Name            string   `json:"name"             binding:"required,min=1,max=100"`
	Species         []string `json:"species"          binding:"required,min=1,dive,oneof=dog cat bird rabbit fish other"`
	FrequencyMonths int      `json:"frequency_months" binding:"required,min=1,max=120"`
	IsMandatory     bool     `json:"is_mandatory"`
}
