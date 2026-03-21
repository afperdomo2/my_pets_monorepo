package models

import "time"

// HealthRecord representa un evento de salud registrado para una mascota.
// Puede estar vinculado al catálogo de salud (health_catalog_id) o ser un ingreso manual.
// Los campos category y name se copian del catálogo al momento de la inserción
// para preservar el historial aunque el catálogo cambie (foto histórica).
// UserID se almacena directamente para permitir filtrado eficiente por usuario sin joins.
//
// HealthRecord solo maneja vacunas (vaccine) y desparasitaciones (deworming).
// Los exámenes se manejan en la tabla exams.
//
// El campo NextDoseDate indica cuándo se debe aplicar la próxima dosis del tratamiento.
// Las aplicaciones de dosis se registran en la tabla vaccine_applications.
type HealthRecord struct {
	ID              string     `json:"id"                gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	PetID           string     `json:"pet_id"            gorm:"type:uuid;not null;index"`
	UserID          string     `json:"user_id"           gorm:"type:uuid;not null;index;constraint:OnDelete:CASCADE"`
	HealthCatalogID *string    `json:"health_catalog_id" gorm:"type:uuid;default:null;constraint:OnDelete:SET NULL"`
	Category        string     `json:"category"          gorm:"type:varchar(20);not null;index"`
	Name            string     `json:"name"              gorm:"type:varchar(100);not null"`
	ApplicationDate time.Time  `json:"application_date"  gorm:"type:date;not null"`
	NextDoseDate    *time.Time `json:"next_dose_date"    gorm:"type:date;default:null;index"`
	Notes           *string    `json:"notes"             gorm:"type:text;default:null"`
	CreatedAt       time.Time  `json:"created_at"        gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt       time.Time  `json:"updated_at"        gorm:"type:timestamptz;not null;default:now()"`

	// Relación con Pet — CASCADE: si se elimina la mascota, se eliminan sus registros.
	Pet Pet `json:"pet" gorm:"foreignKey:PetID;constraint:OnDelete:CASCADE"`

	// Relación con VaccineApplications — un health_record puede tener múltiples aplicaciones.
	VaccineApplications []VaccineApplication `json:"vaccine_applications" gorm:"foreignKey:HealthRecordID"`
}
