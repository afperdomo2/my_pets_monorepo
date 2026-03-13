package health_record

// CreateHealthRecordPayload es el cuerpo aceptado al crear un registro de salud.
// Si se provee health_catalog_id, el handler copiará name y category desde el catálogo.
// Si no se provee health_catalog_id, name y category son obligatorios en el payload.
type CreateHealthRecordPayload struct {
	HealthCatalogID *string `json:"health_catalog_id" binding:"omitempty,uuid"`
	Category        string  `json:"category"          binding:"omitempty,oneof=vaccine deworming exam"`
	Name            string  `json:"name"              binding:"omitempty,max=100"`
	Status          string  `json:"status"            binding:"omitempty,oneof=pending applied"`
	ApplicationDate *string `json:"application_date"  binding:"omitempty"`
	DueDate         string  `json:"due_date"          binding:"required"`
	Notes           *string `json:"notes"             binding:"omitempty,max=1000"`
}

// UpdateHealthRecordPayload es el cuerpo aceptado al actualizar un registro de salud completo.
// health_catalog_id no es actualizable para preservar la integridad del historial.
type UpdateHealthRecordPayload struct {
	Category        string  `json:"category"         binding:"required,oneof=vaccine deworming exam"`
	Name            string  `json:"name"             binding:"required,min=1,max=100"`
	Status          string  `json:"status"           binding:"omitempty,oneof=pending applied"`
	ApplicationDate *string `json:"application_date" binding:"omitempty"`
	DueDate         string  `json:"due_date"         binding:"required"`
	Notes           *string `json:"notes"            binding:"omitempty,max=1000"`
}

// UpdateStatusPayload es el cuerpo aceptado al actualizar solo el status de un registro.
// El status 'overdue' no se persiste en BD; es calculado en runtime por el handler al leer.
// Solo se permiten 'pending' y 'applied' en este endpoint.
type UpdateStatusPayload struct {
	Status          string  `json:"status"           binding:"required,oneof=pending applied"`
	ApplicationDate *string `json:"application_date" binding:"omitempty"`
}
