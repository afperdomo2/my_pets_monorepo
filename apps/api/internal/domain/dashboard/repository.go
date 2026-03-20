package dashboard

import (
	"context"
)

// DashboardSummary contiene los datos principales para el dashboard del usuario.
type DashboardSummary struct {
	TotalPets     int64 `json:"total_pets"`
	HealthyPets   int64 `json:"healthy_pets"`
	PendingTasks  int64 `json:"pending_tasks"`
	OverdueTasks  int64 `json:"overdue_tasks"`
}

// Repository define el contrato de persistencia para el dominio dashboard.
// Todas las operaciones están delimitadas por ownerID (el usuario autenticado).
type Repository interface {
	// GetSummary retorna un resumen de los datos principales del dashboard para el usuario.
	// - total_pets: total de mascotas del usuario
	// - healthy_pets: mascotas sin próxima dosis pendiente de aplicación
	// - pending_tasks: total de registros con próxima dosis programada y no aplicada
	// - overdue_tasks: total de registros con próxima dosis vencida y no aplicada
	GetSummary(ctx context.Context, ownerID string) (DashboardSummary, error)
}
