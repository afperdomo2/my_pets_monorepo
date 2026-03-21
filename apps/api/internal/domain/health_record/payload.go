package health_record

// CreateHealthRecordPayload es el cuerpo aceptado al crear un registro de salud.
// pet_id es obligatorio en el body — ya no se recibe por URL.
// Si se provee health_catalog_id, el handler copiará name y category desde el catálogo.
// Si no se provee health_catalog_id, name y category son obligatorios en el payload.
//
// HealthRecord solo maneja vacunas (vaccine) y desparasitaciones (deworming).
// Los exámenes se manejan en la tabla exams.
type CreateHealthRecordPayload struct {
	PetID           string  `json:"pet_id"            binding:"required,uuid"`
	HealthCatalogID *string `json:"health_catalog_id" binding:"omitempty,uuid"`
	Category        string  `json:"category"          binding:"omitempty,oneof=vaccine deworming"`
	Name            string  `json:"name"              binding:"omitempty,max=100"`
	ApplicationDate string  `json:"application_date"  binding:"required"`
	NextDoseDate    *string `json:"next_dose_date"    binding:"omitempty"`
	Notes           *string `json:"notes"             binding:"omitempty,max=1000"`
}

// UpdateHealthRecordPayload es el cuerpo aceptado al actualizar un registro de salud completo.
// pet_id y health_catalog_id no son actualizables para preservar la integridad del historial.
type UpdateHealthRecordPayload struct {
	Category     string  `json:"category"          binding:"required,oneof=vaccine deworming"`
	Name         string  `json:"name"              binding:"required,min=1,max=100"`
	ApplicationDate *string `json:"application_date" binding:"omitempty"`
	NextDoseDate *string `json:"next_dose_date"    binding:"omitempty"`
	Notes        *string `json:"notes"             binding:"omitempty,max=1000"`
}
