package models

import "time"

// VaccineApplication representa la aplicación de una dosis de vacuna o desparasitación.
// Está vinculado a un health_record que contiene la información general del tratamiento.
// El campo application_date indica cuándo se aplicó la dosis.
// Notes permite almacenar observaciones sobre la aplicación (lote, veterinario, etc.).
type VaccineApplication struct {
	ID             string     `json:"id"             gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	HealthRecordID string     `json:"health_record_id" gorm:"type:uuid;not null;index"`
	ApplicationDate time.Time  `json:"application_date" gorm:"type:date;not null"`
	Notes          *string    `json:"notes"          gorm:"type:text;default:null"`
	CreatedAt      time.Time  `json:"created_at"     gorm:"type:timestamptz;not null;default:now()"`

	// Relación con HealthRecord — CASCADE: si se elimina el registro, se eliminan las aplicaciones.
	HealthRecord HealthRecord `json:"health_record" gorm:"foreignKey:HealthRecordID;constraint:OnDelete:CASCADE"`
}
