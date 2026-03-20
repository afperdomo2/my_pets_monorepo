package models

import "time"

// ExamResult representa un resultado individual de un examen médico veterinario.
// Cada examen puede tener múltiples resultados (parámetros medidos).
// Por ejemplo, un examen de sangre puede tener resultados para: glucosa, colesterol, hemoglobina, etc.
//
// El campo parameter_name es el nombre del parámetro medido (ej: "Glucosa").
// El campo value es el valor obtenido (ej: "90").
// El campo unit es la unidad de medida opcional (ej: "mg/dL").
type ExamResult struct {
	ID            string    `json:"id"             gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ExamID        string    `json:"exam_id"        gorm:"type:uuid;not null;index"`
	ParameterName string    `json:"parameter_name" gorm:"type:varchar(100);not null"`
	Value         string    `json:"value"          gorm:"type:varchar(100);not null"`
	Unit          *string   `json:"unit"           gorm:"type:varchar(50);default:null"`
	CreatedAt     time.Time `json:"created_at"     gorm:"type:timestamptz;not null;default:now()"`

	// Relación con Exam — CASCADE: si se elimina el examen, se eliminan sus resultados.
	Exam Exam `json:"exam" gorm:"foreignKey:ExamID;constraint:OnDelete:CASCADE"`
}
