package models

import (
	"time"

	"github.com/lib/pq"
)

// VaccineCatalog es el catálogo de vacunas disponibles para las mascotas.
// Solo usuarios sistema pueden crear/actualizar/eliminar registros.
// Todos los usuarios autenticados pueden leer el catálogo.
type VaccineCatalog struct {
	ID              string         `json:"id"               gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name            string         `json:"name"             gorm:"type:varchar(100);not null"`
	Species         pq.StringArray `json:"species"          gorm:"type:text[]"`
	FrequencyMonths *int           `json:"frequency_months" gorm:"default:null"`
	IsMandatory     bool           `json:"is_mandatory"     gorm:"not null;default:false"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}
