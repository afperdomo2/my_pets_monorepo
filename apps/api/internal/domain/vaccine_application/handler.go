package vaccine_application

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/my-pets/api/internal/models"
	"github.com/my-pets/api/internal/validation"
)

// Handler contiene las dependencias para los handlers HTTP de vaccine_application.
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

// parseAppID extrae el parámetro de ruta :id de la aplicación de vacuna.
func parseAppID(c *gin.Context) (string, bool) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return "", false
	}
	return id, true
}

// GetApplicationsByHealthRecord maneja GET /api/v1/vaccine-applications/health-record/:id
// Retorna todas las aplicaciones de un health_record específico.
//
//	@Summary	Listar aplicaciones de un registro de salud
//	@Tags		vaccine-applications
//	@Produce	json
//	@Security	CookieAuth
//	@Param		id	path		string	true	"ID del health_record"
//	@Success	200	{object}	map[string]interface{}	"data: []VaccineApplication"
//	@Failure	400	{object}	map[string]string		"id inválido"
//	@Failure	401	{object}	map[string]string		"autenticación requerida"
//	@Failure	500	{object}	map[string]string		"mensaje de error"
//	@Router		/api/v1/vaccine-applications/health-record/{id} [get]
func (h *Handler) GetApplicationsByHealthRecord(c *gin.Context) {
	healthRecordID := c.Param("id")
	if healthRecordID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "health_record_id inválido"})
		return
	}

	apps, err := h.repo.GetByHealthRecordID(c.Request.Context(), healthRecordID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener aplicaciones"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": apps})
}

// GetApplicationByID maneja GET /api/v1/vaccine-applications/:id
// Retorna una aplicación específica por su ID.
//
//	@Summary	Obtener aplicación por ID
//	@Tags		vaccine-applications
//	@Produce	json
//	@Security	CookieAuth
//	@Param		id	path		string	true	"ID de la aplicación"
//	@Success	200	{object}	map[string]interface{}	"data: VaccineApplication"
//	@Failure	400	{object}	map[string]string		"id inválido"
//	@Failure	404	{object}	map[string]string		"aplicación no encontrada"
//	@Failure	401	{object}	map[string]string		"autenticación requerida"
//	@Failure	500	{object}	map[string]string		"mensaje de error"
//	@Router		/api/v1/vaccine-applications/{id} [get]
func (h *Handler) GetApplicationByID(c *gin.Context) {
	id, ok := parseAppID(c)
	if !ok {
		return
	}

	app, err := h.repo.GetByID(c.Request.Context(), id)
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "aplicación no encontrada"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener aplicación"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": app})
}

// CreateApplication maneja POST /api/v1/vaccine-applications
// Crea una nueva aplicación de vacuna/desparasitación.
//
//	@Summary	Crear una aplicación de vacuna
//	@Tags		vaccine-applications
//	@Accept		json
//	@Produce	json
//	@Security	CookieAuth
//	@Param		app	body		CreatePayload	true	"Datos de la aplicación"
//	@Success	201	{object}	map[string]interface{}	"data: VaccineApplication"
//	@Failure	400	{object}	map[string]string		"error de validación"
//	@Failure	401	{object}	map[string]string		"autenticación requerida"
//	@Failure	500	{object}	map[string]string		"mensaje de error"
//	@Router		/api/v1/vaccine-applications [post]
func (h *Handler) CreateApplication(c *gin.Context) {
	var payload CreatePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validation.Translate(err)})
		return
	}

	appDate, err := parseDateString(payload.ApplicationDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "formato de application_date inválido"})
		return
	}

	app := models.VaccineApplication{
		HealthRecordID: payload.HealthRecordID,
		ApplicationDate: appDate,
		Notes:          payload.Notes,
	}

	created, err := h.repo.Create(c.Request.Context(), app)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al crear aplicación"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": created})
}

// UpdateApplication maneja PUT /api/v1/vaccine-applications/:id
// Actualiza una aplicación existente.
//
//	@Summary	Actualizar una aplicación de vacuna
//	@Tags		vaccine-applications
//	@Accept		json
//	@Produce	json
//	@Security	CookieAuth
//	@Param		id	path		string		true	"ID de la aplicación"
//	@Param		app	body		UpdatePayload	true	"Datos de la aplicación"
//	@Success	200	{object}	map[string]interface{}	"data: VaccineApplication"
//	@Failure	400	{object}	map[string]string		"error de validación"
//	@Failure	401	{object}	map[string]string		"autenticación requerida"
//	@Failure	404	{object}	map[string]string		"aplicación no encontrada"
//	@Failure	500	{object}	map[string]string		"mensaje de error"
//	@Router		/api/v1/vaccine-applications/{id} [put]
func (h *Handler) UpdateApplication(c *gin.Context) {
	id, ok := parseAppID(c)
	if !ok {
		return
	}

	var payload UpdatePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validation.Translate(err)})
		return
	}

	app, err := h.repo.GetByID(c.Request.Context(), id)
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "aplicación no encontrada"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener aplicación"})
		return
	}

	// Actualizar campos si se proporcionan
	if payload.ApplicationDate != "" {
		appDate, err := parseDateString(payload.ApplicationDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "formato de application_date inválido"})
			return
		}
		app.ApplicationDate = appDate
	}
	if payload.Notes != nil {
		app.Notes = payload.Notes
	}

	updated, err := h.repo.Update(c.Request.Context(), app)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al actualizar aplicación"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": updated})
}

// DeleteApplication maneja DELETE /api/v1/vaccine-applications/:id
// Elimina una aplicación de vacuna.
//
//	@Summary	Eliminar una aplicación de vacuna
//	@Tags		vaccine-applications
//	@Produce	json
//	@Security	CookieAuth
//	@Param		id	path		string	true	"ID de la aplicación"
//	@Success	200	{object}	map[string]string	"message: aplicación eliminada"
//	@Failure	400	{object}	map[string]string	"id inválido"
//	@Failure	401	{object}	map[string]string	"autenticación requerida"
//	@Failure	404	{object}	map[string]string	"aplicación no encontrada"
//	@Failure	500	{object}	map[string]string	"mensaje de error"
//	@Router		/api/v1/vaccine-applications/{id} [delete]
func (h *Handler) DeleteApplication(c *gin.Context) {
	id, ok := parseAppID(c)
	if !ok {
		return
	}

	err := h.repo.Delete(c.Request.Context(), id)
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "aplicación no encontrada"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al eliminar aplicación"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "aplicación eliminada"})
}
