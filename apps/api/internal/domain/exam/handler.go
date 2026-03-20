package exam

import (
	"errors"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/my-pets/api/internal/models"
	"github.com/my-pets/api/internal/validation"
)

// Handler contiene las dependencias para los handlers HTTP de exam.
type Handler struct {
	repo Repository
}

// NewHandler construye un Handler con el repositorio dado.
func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

// ownerID extrae el ID del usuario autenticado del contexto Gin (seteado por el middleware JWT).
func ownerID(c *gin.Context) string {
	return c.GetString("userID")
}

// parseExamID extrae el parámetro de ruta :id del examen.
func parseExamID(c *gin.Context) (string, bool) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return "", false
	}
	return id, true
}

// parsePetIDParam extrae el parámetro de ruta :pet_id.
func parsePetIDParam(c *gin.Context) (string, bool) {
	id := c.Param("pet_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pet_id inválido"})
		return "", false
	}
	return id, true
}

// GetAllExams maneja GET /api/v1/exams
// Lista todos los exámenes del usuario autenticado con paginación.
//
//	@Summary	Listar todos los exámenes del usuario (paginado)
//	@Tags		exams
//	@Produce	json
//	@Security	CookieAuth
//	@Param		page		query		int						false	"Número de página (por defecto 1)"
//	@Param		per_page	query		int						false	"Elementos por página (por defecto 10)"
//	@Success	200	{object}	map[string]interface{}	"data: []Exam, total: int, page: int, per_page: int, total_pages: int"
//	@Failure	401	{object}	map[string]string		"autenticación requerida"
//	@Failure	500	{object}	map[string]string		"mensaje de error"
//	@Router		/api/v1/exams [get]
func (h *Handler) GetAllExams(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	exams, total, err := h.repo.GetAllByOwner(c.Request.Context(), ownerID(c), page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener exámenes"})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	c.JSON(http.StatusOK, gin.H{
		"data":        exams,
		"total":       total,
		"page":        page,
		"per_page":    perPage,
		"total_pages": totalPages,
	})
}

// GetExamsByPet maneja GET /api/v1/exams/pets/:pet_id
// Lista los exámenes de una mascota específica del usuario autenticado.
//
//	@Summary	Listar exámenes de una mascota (paginado)
//	@Tags		exams
//	@Produce	json
//	@Security	CookieAuth
//	@Param		pet_id		path		string	true	"ID de la mascota"
//	@Param		page		query		int		false	"Número de página (por defecto 1)"
//	@Param		per_page	query		int		false	"Elementos por página (por defecto 10)"
//	@Success	200	{object}	map[string]interface{}	"data: []Exam, total: int, page: int, per_page: int, total_pages: int"
//	@Failure	400	{object}	map[string]string		"pet_id inválido"
//	@Failure	401	{object}	map[string]string		"autenticación requerida"
//	@Failure	500	{object}	map[string]string		"mensaje de error"
//	@Router		/api/v1/exams/pets/{pet_id} [get]
func (h *Handler) GetExamsByPet(c *gin.Context) {
	petID, ok := parsePetIDParam(c)
	if !ok {
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	exams, total, err := h.repo.GetByPetID(c.Request.Context(), petID, ownerID(c), page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener exámenes"})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	c.JSON(http.StatusOK, gin.H{
		"data":        exams,
		"total":       total,
		"page":        page,
		"per_page":    perPage,
		"total_pages": totalPages,
	})
}

// GetExamByID maneja GET /api/v1/exams/:id
// Retorna un examen específico con sus resultados.
//
//	@Summary	Obtener examen por ID con resultados
//	@Tags		exams
//	@Produce	json
//	@Security	CookieAuth
//	@Param		id	path		string	true	"ID del examen"
//	@Success	200	{object}	map[string]interface{}	"data: Exam, results: []ExamResult"
//	@Failure	400	{object}	map[string]string		"id inválido"
//	@Failure	404	{object}	map[string]string		"examen no encontrado"
//	@Failure	401	{object}	map[string]string		"autenticación requerida"
//	@Failure	500	{object}	map[string]string		"mensaje de error"
//	@Router		/api/v1/exams/{id} [get]
func (h *Handler) GetExamByID(c *gin.Context) {
	id, ok := parseExamID(c)
	if !ok {
		return
	}

	exam, results, err := h.repo.GetByIDWithResults(c.Request.Context(), id, ownerID(c))
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "examen no encontrado"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener examen"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    exam,
		"results": results,
	})
}

// CreateExam maneja POST /api/v1/exams
// Crea un nuevo examen veterinario.
//
//	@Summary	Crear un examen
//	@Tags		exams
//	@Accept		json
//	@Produce	json
//	@Security	CookieAuth
//	@Param		exam	body		CreatePayload	true	"Datos del examen"
//	@Success	201	{object}	map[string]interface{}	"data: Exam"
//	@Failure	400	{object}	map[string]string		"error de validación"
//	@Failure	401	{object}	map[string]string		"autenticación requerida"
//	@Failure	404	{object}	map[string]string		"mascota no encontrada"
//	@Failure	500	{object}	map[string]string		"mensaje de error"
//	@Router		/api/v1/exams [post]
func (h *Handler) CreateExam(c *gin.Context) {
	var payload CreatePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validation.Translate(err)})
		return
	}

	exam := models.Exam{
		PetID:  payload.PetID,
		UserID: ownerID(c),
		Name:   payload.Name,
		Reason: payload.Reason,
		Status: "scheduled",
		Notes:  payload.Notes,
	}

	// Si se provee scheduled_date, parsearla
	if payload.ScheduledDate != nil {
		scheduledDate, err := parseDateString(*payload.ScheduledDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "formato de scheduled_date inválido"})
			return
		}
		exam.ScheduledDate = &scheduledDate
	}

	// Si el status es "completed" o se provee completed_date
	if payload.Status == "completed" || payload.CompletedDate != nil {
		exam.Status = "completed"
		completedDate := time.Now().Truncate(24 * time.Hour)
		if payload.CompletedDate != nil {
			var err error
			completedDate, err = parseDateString(*payload.CompletedDate)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "formato de completed_date inválido"})
				return
			}
		}
		exam.CompletedDate = &completedDate
		if exam.ScheduledDate == nil {
			exam.ScheduledDate = &completedDate
		}
	}

	created, err := h.repo.Create(c.Request.Context(), exam)
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "mascota no encontrada o no pertenece al usuario"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al crear examen"})
		return
	}

	// Si hay resultados, crearlos
	if len(payload.Results) > 0 {
		results := make([]models.ExamResult, len(payload.Results))
		for i, r := range payload.Results {
			results[i] = models.ExamResult{
				ExamID:        created.ID,
				ParameterName: r.ParameterName,
				Value:         r.Value,
			}
			if r.Unit != nil {
				results[i].Unit = r.Unit
			}
		}
		if err := h.repo.CreateResults(c.Request.Context(), results); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error al crear resultados"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"data": created})
}

// UpdateExam maneja PUT /api/v1/exams/:id
// Actualiza un examen existente.
//
//	@Summary	Actualizar un examen
//	@Tags		exams
//	@Accept		json
//	@Produce	json
//	@Security	CookieAuth
//	@Param		id	path		string		true	"ID del examen"
//	@Param		exam	body		UpdatePayload	true	"Datos del examen"
//	@Success	200	{object}	map[string]interface{}	"data: Exam"
//	@Failure	400	{object}	map[string]string		"error de validación"
//	@Failure	401	{object}	map[string]string		"autenticación requerida"
//	@Failure	404	{object}	map[string]string		"examen no encontrado"
//	@Failure	500	{object}	map[string]string		"mensaje de error"
//	@Router		/api/v1/exams/{id} [put]
func (h *Handler) UpdateExam(c *gin.Context) {
	id, ok := parseExamID(c)
	if !ok {
		return
	}

	var payload UpdatePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validation.Translate(err)})
		return
	}

	exam, err := h.repo.GetByID(c.Request.Context(), id, ownerID(c))
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "examen no encontrado"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener examen"})
		return
	}

	// Actualizar campos si se proporcionan
	if payload.Name != "" {
		exam.Name = payload.Name
	}
	if payload.Reason != nil {
		exam.Reason = payload.Reason
	}
	if payload.ScheduledDate != nil {
		scheduledDate, err := parseDateString(*payload.ScheduledDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "formato de scheduled_date inválido"})
			return
		}
		exam.ScheduledDate = &scheduledDate
	}
	if payload.Notes != nil {
		exam.Notes = payload.Notes
	}

	updated, err := h.repo.Update(c.Request.Context(), exam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al actualizar examen"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": updated})
}

// ScheduleExam maneja PATCH /api/v1/exams/:id/schedule
// Programa o reprograma un examen.
//
//	@Summary	Programar un examen
//	@Tags		exams
//	@Accept		json
//	@Produce	json
//	@Security	CookieAuth
//	@Param		id	path		string			true	"ID del examen"
//	@Param		data	body		SchedulePayload	true	"Fecha programada"
//	@Success	200	{object}	map[string]interface{}	"data: Exam"
//	@Failure	400	{object}	map[string]string		"error de validación"
//	@Failure	401	{object}	map[string]string		"autenticación requerida"
//	@Failure	404	{object}	map[string]string		"examen no encontrado"
//	@Failure	500	{object}	map[string]string		"mensaje de error"
//	@Router		/api/v1/exams/{id}/schedule [patch]
func (h *Handler) ScheduleExam(c *gin.Context) {
	id, ok := parseExamID(c)
	if !ok {
		return
	}

	var payload SchedulePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validation.Translate(err)})
		return
	}

	exam, err := h.repo.GetByID(c.Request.Context(), id, ownerID(c))
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "examen no encontrado"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener examen"})
		return
	}

	scheduledDate, err := parseDateString(payload.ScheduledDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "formato de scheduled_date inválido"})
		return
	}

	exam.ScheduledDate = &scheduledDate
	if exam.Status == "completed" {
		exam.Status = "scheduled"
		exam.CompletedDate = nil
	}

	updated, err := h.repo.Update(c.Request.Context(), exam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al programar examen"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": updated})
}

// CompleteExam maneja PATCH /api/v1/exams/:id/complete
// Marca un examen como completado y opcionalmente agrega resultados.
//
//	@Summary	Completar un examen
//	@Tags		exams
//	@Accept		json
//	@Produce	json
//	@Security	CookieAuth
//	@Param		id	path		string			true	"ID del examen"
//	@Param		data	body		CompletePayload	true	"Fecha de completado y resultados"
//	@Success	200	{object}	map[string]interface{}	"data: Exam"
//	@Failure	400	{object}	map[string]string		"error de validación"
//	@Failure	401	{object}	map[string]string		"autenticación requerida"
//	@Failure	404	{object}	map[string]string		"examen no encontrado"
//	@Failure	500	{object}	map[string]string		"mensaje de error"
//	@Router		/api/v1/exams/{id}/complete [patch]
func (h *Handler) CompleteExam(c *gin.Context) {
	id, ok := parseExamID(c)
	if !ok {
		return
	}

	var payload CompletePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validation.Translate(err)})
		return
	}

	exam, err := h.repo.GetByID(c.Request.Context(), id, ownerID(c))
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "examen no encontrado"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener examen"})
		return
	}

	completedDate, err := parseDateString(payload.CompletedDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "formato de completed_date inválido"})
		return
	}

	exam.Status = "completed"
	exam.CompletedDate = &completedDate
	if exam.ScheduledDate == nil {
		exam.ScheduledDate = &completedDate
	}

	updated, err := h.repo.Update(c.Request.Context(), exam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al completar examen"})
		return
	}

	// Si hay resultados, primero eliminar los existentes y luego crear los nuevos
	if len(payload.Results) > 0 {
		if err := h.repo.DeleteResultsByExamID(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error al eliminar resultados anteriores"})
			return
		}

		results := make([]models.ExamResult, len(payload.Results))
		for i, r := range payload.Results {
			results[i] = models.ExamResult{
				ExamID:        updated.ID,
				ParameterName: r.ParameterName,
				Value:         r.Value,
			}
			if r.Unit != nil {
				results[i].Unit = r.Unit
			}
		}
		if err := h.repo.CreateResults(c.Request.Context(), results); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error al crear resultados"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": updated})
}

// DeleteExam maneja DELETE /api/v1/exams/:id
// Elimina un examen.
//
//	@Summary	Eliminar un examen
//	@Tags		exams
//	@Produce	json
//	@Security	CookieAuth
//	@Param		id	path		string	true	"ID del examen"
//	@Success	200	{object}	map[string]string	"message: examen eliminado"
//	@Failure	400	{object}	map[string]string	"id inválido"
//	@Failure	401	{object}	map[string]string	"autenticación requerida"
//	@Failure	404	{object}	map[string]string	"examen no encontrado"
//	@Failure	500	{object}	map[string]string	"mensaje de error"
//	@Router		/api/v1/exams/{id} [delete]
func (h *Handler) DeleteExam(c *gin.Context) {
	id, ok := parseExamID(c)
	if !ok {
		return
	}

	err := h.repo.Delete(c.Request.Context(), id, ownerID(c))
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "examen no encontrado"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al eliminar examen"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "examen eliminado"})
}
