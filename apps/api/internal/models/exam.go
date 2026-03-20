package models

import "time"

// Exam representa un examen médico veterinario programado o completado para una mascota.
// Contiene información general sobre el examen (nombre, motivo, estado, fechas).
// Los resultados detallados del examen se almacenan en la tabla exam_results.
//
// Estados posibles:
//   - "scheduled": El examen está programado, pendiente de realizar
//   - "completed": El examen ha sido completado y tiene resultados
type Exam struct {
	ID            string     `json:"id"             gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	PetID         string     `json:"pet_id"         gorm:"type:uuid;not null;index"`
	UserID        string     `json:"user_id"        gorm:"type:uuid;not null;index"`
	Name          string     `json:"name"           gorm:"type:varchar(100);not null"`
	Reason        *string    `json:"reason"         gorm:"type:text;default:null"`
	Status        string     `json:"status"         gorm:"type:varchar(20);not null;default:'scheduled';index"`
	ScheduledDate *time.Time `json:"scheduled_date" gorm:"type:date;default:null"`
	CompletedDate *time.Time `json:"completed_date" gorm:"type:date;default:null"`
	Notes         *string    `json:"notes"          gorm:"type:text;default:null"`
	CreatedAt     time.Time  `json:"created_at"     gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt     time.Time  `json:"updated_at"     gorm:"type:timestamptz;not null;default:now()"`

	// Relación con Pet — CASCADE: si se elimina la mascota, se eliminan sus exámenes.
	Pet Pet `json:"pet" gorm:"foreignKey:PetID;constraint:OnDelete:CASCADE"`
}
