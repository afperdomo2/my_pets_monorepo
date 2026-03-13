package models

import (
	"time"

	"github.com/lib/pq"
)

// HealthCatalog es el catálogo de registros de salud disponibles para las mascotas.
// Agrupa vacunas, desparasitaciones y exámenes programáticos.
// Solo usuarios sistema pueden crear/actualizar/eliminar registros.
// Todos los usuarios autenticados pueden leer el catálogo.
type HealthCatalog struct {
	ID              string         `json:"id"               gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name            string         `json:"name"             gorm:"type:varchar(100);not null"`
	Category        string         `json:"category"         gorm:"type:varchar(20);not null;default:'vaccine'"`
	Description     string         `json:"description"      gorm:"type:text;not null;default:''"`
	Species         pq.StringArray `json:"species"          gorm:"type:text[]"`
	FrequencyMonths *int           `json:"frequency_months" gorm:"default:null"`
	IsMandatory     bool           `json:"is_mandatory"     gorm:"not null;default:false"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}
