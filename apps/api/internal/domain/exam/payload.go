package exam

// CreatePayload es el cuerpo aceptado al crear un examen.
// pet_id es obligatorio — vincula el examen a la mascota.
// name es el nombre del examen (ej: "Análisis de sangre").
// reason es el motivo o razón del examen (opcional).
// scheduled_date es la fecha programada para el examen.
// status puede ser "scheduled" o "completed".
// results es un array opcional de resultados (solo si status es "completed").
type CreatePayload struct {
	PetID         string           `json:"pet_id" binding:"required,uuid"`
	Name          string           `json:"name" binding:"required,min=1,max=100"`
	Reason        *string          `json:"reason" binding:"omitempty,max=1000"`
	ScheduledDate *string          `json:"scheduled_date" binding:"omitempty"`
	Status        string           `json:"status" binding:"omitempty,oneof=scheduled completed"`
	CompletedDate *string          `json:"completed_date" binding:"omitempty"`
	Notes         *string          `json:"notes" binding:"omitempty,max=1000"`
	Results       []ResultPayload  `json:"results" binding:"omitempty"`
}

// ResultPayload representa un resultado individual de un examen.
type ResultPayload struct {
	ParameterName string  `json:"parameter_name" binding:"required,min=1,max=100"`
	Value         string  `json:"value" binding:"required,min=1,max=100"`
	Unit          *string `json:"unit" binding:"omitempty,max=50"`
}

// UpdatePayload es el cuerpo aceptado al actualizar un examen.
// Solo permite actualizar campos editables.
// pet_id y status no son actualizables directamente (usar endpoints específicos).
type UpdatePayload struct {
	Name          string  `json:"name" binding:"omitempty,min=1,max=100"`
	Reason        *string `json:"reason" binding:"omitempty,max=1000"`
	ScheduledDate *string `json:"scheduled_date" binding:"omitempty"`
	Notes         *string `json:"notes" binding:"omitempty,max=1000"`
}

// SchedulePayload es el cuerpo aceptado al programar un examen.
type SchedulePayload struct {
	ScheduledDate string `json:"scheduled_date" binding:"required"`
}

// CompletePayload es el cuerpo aceptado al completar un examen.
// completed_date es obligatoria.
// results es opcional — los resultados del examen.
type CompletePayload struct {
	CompletedDate string          `json:"completed_date" binding:"required"`
	Results       []ResultPayload `json:"results" binding:"omitempty"`
}
