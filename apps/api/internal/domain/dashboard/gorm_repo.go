package dashboard

import (
	"context"

	"gorm.io/gorm"
)

// gormRepo implementa Repository usando GORM.
type gormRepo struct {
	db *gorm.DB
}

// NewGormRepo crea una instancia de gormRepo con la conexión dada.
func NewGormRepo(db *gorm.DB) Repository {
	return &gormRepo{db: db}
}

// GetSummary implementa Repository.GetSummary.
// Ejecuta 4 consultas optimizadas usando el campo user_id para evitar JOINs innecesarios.
// Nota: Usa next_dose_date y application_date en lugar de status y due_date (eliminados).
func (r *gormRepo) GetSummary(ctx context.Context, ownerID string) (DashboardSummary, error) {
	var summary DashboardSummary

	// 1. Total de pets del usuario
	var totalPets int64
	if err := r.db.Table("pets").Where("user_id = ?", ownerID).Count(&totalPets).Error; err != nil {
		return summary, err
	}
	summary.TotalPets = totalPets

	// 2. Mascotas saludables: pets que NO tienen próxima dosis pendiente
	// Usamos el campo user_id de health_records para evitar JOINs
	// Un pet está saludable si no tiene next_dose_date programada o ya está aplicada
	var healthyPets int64
	query := `
		SELECT COUNT(*) FROM pets p
		WHERE p.user_id = ?
		AND NOT EXISTS (
			SELECT 1 FROM health_records hr
			WHERE hr.pet_id = p.id
			AND hr.user_id = ?
			AND hr.next_dose_date IS NOT NULL
			AND hr.application_date IS NULL
		)
	`
	if err := r.db.Raw(query, ownerID, ownerID).Scan(&healthyPets).Error; err != nil {
		return summary, err
	}
	summary.HealthyPets = healthyPets

	// 3. Total de tareas pendientes (next_dose_date IS NOT NULL AND application_date IS NULL)
	// Filtramos por user_id directamente sin necesidad de JOIN
	var pendingTasks int64
	if err := r.db.Table("health_records").
		Where("user_id = ? AND next_dose_date IS NOT NULL AND application_date IS NULL", ownerID).
		Count(&pendingTasks).Error; err != nil {
		return summary, err
	}
	summary.PendingTasks = pendingTasks

	// 4. Total de tareas vencidas (next_dose_date < hoy y application_date IS NULL)
	// Filtramos por user_id directamente sin necesidad de JOIN
	var overdueTasks int64
	if err := r.db.Table("health_records").
		Where("user_id = ? AND next_dose_date < CURRENT_DATE AND application_date IS NULL", ownerID).
		Count(&overdueTasks).Error; err != nil {
		return summary, err
	}
	summary.OverdueTasks = overdueTasks

	return summary, nil
}
