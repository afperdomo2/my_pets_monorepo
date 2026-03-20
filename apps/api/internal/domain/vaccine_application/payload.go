package vaccine_application

// CreatePayload es el cuerpo aceptado al crear una aplicación de vacuna.
// health_record_id es obligatorio — vincula la aplicación al registro de salud.
// application_date es la fecha en que se aplicó la dosis.
// notes es opcional — permite almacenar observaciones (lote, veterinario, etc.).
type CreatePayload struct {
	HealthRecordID string  `json:"health_record_id" binding:"required,uuid"`
	ApplicationDate  string  `json:"application_date" binding:"required"`
	Notes          *string `json:"notes" binding:"omitempty,max=1000"`
}

// UpdatePayload es el cuerpo aceptado al actualizar una aplicación de vacuna.
// Solo permite actualizar application_date y notes.
// health_record_id no es actualizable para preservar la integridad.
type UpdatePayload struct {
	ApplicationDate string  `json:"application_date" binding:"omitempty"`
	Notes           *string `json:"notes" binding:"omitempty,max=1000"`
}
