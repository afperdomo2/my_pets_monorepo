package exam

import "github.com/gin-gonic/gin"

// RegisterRoutes monta todos los endpoints de exam sobre el RouterGroup recibido.
//
// Rutas bajo /exams:
//   - GET    /exams                    → todos los exámenes del usuario (paginado)
//   - POST   /exams                    → crear examen
//   - GET    /exams/pets/:pet_id       → exámenes de una mascota (paginado)
//   - GET    /exams/:id                → obtener examen con resultados
//   - PUT    /exams/:id                → actualizar examen
//   - PATCH  /exams/:id/schedule       → programar/reprogramar examen
//   - PATCH  /exams/:id/complete       → completar examen
//   - DELETE /exams/:id                → eliminar examen
func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	exams := rg.Group("/exams")
	{
		// Listar todos los exámenes del usuario y crear
		exams.GET("", h.GetAllExams)
		exams.POST("", h.CreateExam)

		// Exámenes de una mascota específica (más específico, va primero)
		pets := exams.Group("/pets")
		{
			pets.GET("/:pet_id", h.GetExamsByPet)
		}

		// Obtener, actualizar, completar y eliminar examen por ID
		exams.GET("/:id", h.GetExamByID)
		exams.PUT("/:id", h.UpdateExam)
		exams.PATCH("/:id/schedule", h.ScheduleExam)
		exams.PATCH("/:id/complete", h.CompleteExam)
		exams.DELETE("/:id", h.DeleteExam)
	}
}
