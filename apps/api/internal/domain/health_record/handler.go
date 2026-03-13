package health_record

import (
	"errors"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	healthCatalog "github.com/my-pets/api/internal/domain/health_catalog"
	"github.com/my-pets/api/internal/models"
	"github.com/my-pets/api/internal/validation"
)

// Handler contiene las dependencias para los handlers HTTP de health_record.
type Handler struct {
	repo              Repository
	healthCatalogRepo healthCatalog.Repository
}

// NewHandler construye un Handler con el repositorio dado.
// healthCatalogRepo es opcional: se usa para copiar name y category al crear desde catálogo.
func NewHandler(repo Repository, healthCatalogRepo healthCatalog.Repository) *Handler {
	return &Handler{repo: repo, healthCatalogRepo: healthCatalogRepo}
}

// ownerID extrae el ID del usuario autenticado del contexto Gin (seteado por el middleware JWT).
func ownerID(c *gin.Context) string {
	return c.GetString("userID")
}

// petID extrae el parámetro de ruta :id que representa el pet_id.
// Se usa :id (en lugar de :pet_id) para evitar conflictos con las rutas /pets/:id del dominio pet.
func parsePetID(c *gin.Context) (string, bool) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pet_id inválido"})
		return "", false
	}
	return id, true
}

// parseRecordID extrae el parámetro de ruta :record_id del registro de salud.
func parseRecordID(c *gin.Context) (string, bool) {
	id := c.Param("record_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return "", false
	}
	return id, true
}

// resolveOverdue calcula el status efectivo de un registro en runtime.
// Si el status almacenado es 'pending' y due_date ya pasó, retorna 'overdue'.
// El valor 'overdue' nunca se persiste en la base de datos.
func resolveOverdue(rec models.HealthRecord) models.HealthRecord {
	if rec.Status == "pending" && rec.DueDate.Before(time.Now().Truncate(24*time.Hour)) {
		rec.Status = "overdue"
	}
	return rec
}

// GetHealthRecords maneja GET /api/v1/pets/:pet_id/health-records
//
//	@Summary	Listar registros de salud de una mascota (paginado)
//	@Tags		health-records
//	@Produce	json
//	@Security	CookieAuth
//	@Param		pet_id		path		string					true	"ID de la mascota"
//	@Param		page		query		int						false	"Número de página (por defecto 1)"
//	@Param		per_page	query		int						false	"Elementos por página (por defecto 10)"
//	@Success	200	{object}	map[string]interface{}	"data: []HealthRecord, total: int, page: int, per_page: int, total_pages: int"
//	@Failure	400	{object}	map[string]string		"pet_id inválido"
//	@Failure	401	{object}	map[string]string		"autenticación requerida"
//	@Failure	404	{object}	map[string]string		"mascota no encontrada"
//	@Failure	500	{object}	map[string]string		"mensaje de error"
//	@Router		/api/v1/pets/{pet_id}/health-records [get]
func (h *Handler) GetHealthRecords(c *gin.Context) {
	petID, ok := parsePetID(c)
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

	records, total, err := h.repo.GetPaginated(c.Request.Context(), petID, ownerID(c), page, perPage)
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "mascota no encontrada"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener registros de salud"})
		return
	}

	// Calcular overdue en runtime para cada registro.
	for i, rec := range records {
		records[i] = resolveOverdue(rec)
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	c.JSON(http.StatusOK, gin.H{
		"data":        records,
		"total":       total,
		"page":        page,
		"per_page":    perPage,
		"total_pages": totalPages,
	})
}

// GetHealthRecord maneja GET /api/v1/pets/:pet_id/health-records/:id
//
//	@Summary	Obtener un registro de salud por ID
//	@Tags		health-records
//	@Produce	json
//	@Security	CookieAuth
//	@Param		pet_id	path		string					true	"ID de la mascota"
//	@Param		id		path		string					true	"ID del registro de salud"
//	@Success	200	{object}	map[string]interface{}	"data: HealthRecord"
//	@Failure	400	{object}	map[string]string		"id inválido"
//	@Failure	401	{object}	map[string]string		"autenticación requerida"
//	@Failure	404	{object}	map[string]string		"registro no encontrado"
//	@Failure	500	{object}	map[string]string		"mensaje de error"
//	@Router		/api/v1/pets/{pet_id}/health-records/{id} [get]
func (h *Handler) GetHealthRecord(c *gin.Context) {
	petID, ok := parsePetID(c)
	if !ok {
		return
	}
	id, ok := parseRecordID(c)
	if !ok {
		return
	}

	rec, err := h.repo.GetByID(c.Request.Context(), id, petID, ownerID(c))
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "registro de salud no encontrado"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener registro de salud"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resolveOverdue(rec)})
}

// CreateHealthRecord maneja POST /api/v1/pets/:pet_id/health-records
// Si se provee health_catalog_id, copia name y category desde el catálogo automáticamente.
//
//	@Summary	Crear un registro de salud para una mascota
//	@Tags		health-records
//	@Accept		json
//	@Produce	json
//	@Security	CookieAuth
//	@Param		pet_id	path		string						true	"ID de la mascota"
//	@Param		record	body		CreateHealthRecordPayload	true	"Datos del registro"
//	@Success	201	{object}	map[string]interface{}	"data: HealthRecord"
//	@Failure	400	{object}	map[string]string		"error de validación"
//	@Failure	401	{object}	map[string]string		"autenticación requerida"
//	@Failure	404	{object}	map[string]string		"mascota o catálogo no encontrado"
//	@Failure	500	{object}	map[string]string		"mensaje de error"
//	@Router		/api/v1/pets/{pet_id}/health-records [post]
func (h *Handler) CreateHealthRecord(c *gin.Context) {
	petID, ok := parsePetID(c)
	if !ok {
		return
	}

	var payload CreateHealthRecordPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validation.Translate(err)})
		return
	}

	// Si se provee health_catalog_id, copiar name y category desde el catálogo.
	if payload.HealthCatalogID != nil {
		catalog, err := h.healthCatalogRepo.GetByID(c.Request.Context(), *payload.HealthCatalogID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "catálogo de salud no encontrado"})
			return
		}
		payload.Name = catalog.Name
		payload.Category = catalog.Category
	}

	// Validar que name y category estén presentes (sea desde catálogo o manual).
	if payload.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "el campo name es obligatorio cuando no se provee health_catalog_id"})
		return
	}
	if payload.Category == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "el campo category es obligatorio cuando no se provee health_catalog_id"})
		return
	}

	created, err := h.repo.Create(c.Request.Context(), petID, ownerID(c), payload)
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "mascota no encontrada"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al crear registro de salud"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": resolveOverdue(created)})
}

// UpdateHealthRecord maneja PUT /api/v1/pets/:pet_id/health-records/:id
//
//	@Summary	Actualizar un registro de salud completo
//	@Tags		health-records
//	@Accept		json
//	@Produce	json
//	@Security	CookieAuth
//	@Param		pet_id	path		string						true	"ID de la mascota"
//	@Param		id		path		string						true	"ID del registro de salud"
//	@Param		record	body		UpdateHealthRecordPayload	true	"Datos del registro"
//	@Success	200	{object}	map[string]interface{}	"data: HealthRecord"
//	@Failure	400	{object}	map[string]string		"error de validación"
//	@Failure	401	{object}	map[string]string		"autenticación requerida"
//	@Failure	404	{object}	map[string]string		"registro no encontrado"
//	@Failure	500	{object}	map[string]string		"mensaje de error"
//	@Router		/api/v1/pets/{pet_id}/health-records/{id} [put]
func (h *Handler) UpdateHealthRecord(c *gin.Context) {
	petID, ok := parsePetID(c)
	if !ok {
		return
	}
	id, ok := parseRecordID(c)
	if !ok {
		return
	}

	var payload UpdateHealthRecordPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validation.Translate(err)})
		return
	}

	updated, err := h.repo.Update(c.Request.Context(), id, petID, ownerID(c), payload)
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "registro de salud no encontrado"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al actualizar registro de salud"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resolveOverdue(updated)})
}

// UpdateHealthRecordStatus maneja PATCH /api/v1/pets/:pet_id/health-records/:id/status
//
//	@Summary	Actualizar el status de un registro de salud
//	@Tags		health-records
//	@Accept		json
//	@Produce	json
//	@Security	CookieAuth
//	@Param		pet_id	path		string				true	"ID de la mascota"
//	@Param		id		path		string				true	"ID del registro de salud"
//	@Param		status	body		UpdateStatusPayload	true	"Nuevo status"
//	@Success	200	{object}	map[string]interface{}	"data: HealthRecord"
//	@Failure	400	{object}	map[string]string		"error de validación"
//	@Failure	401	{object}	map[string]string		"autenticación requerida"
//	@Failure	404	{object}	map[string]string		"registro no encontrado"
//	@Failure	500	{object}	map[string]string		"mensaje de error"
//	@Router		/api/v1/pets/{pet_id}/health-records/{id}/status [patch]
func (h *Handler) UpdateHealthRecordStatus(c *gin.Context) {
	petID, ok := parsePetID(c)
	if !ok {
		return
	}
	id, ok := parseRecordID(c)
	if !ok {
		return
	}

	var payload UpdateStatusPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validation.Translate(err)})
		return
	}

	updated, err := h.repo.UpdateStatus(c.Request.Context(), id, petID, ownerID(c), payload)
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "registro de salud no encontrado"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al actualizar status del registro"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resolveOverdue(updated)})
}

// DeleteHealthRecord maneja DELETE /api/v1/pets/:pet_id/health-records/:id
//
//	@Summary	Eliminar un registro de salud
//	@Tags		health-records
//	@Produce	json
//	@Security	CookieAuth
//	@Param		pet_id	path		string	true	"ID de la mascota"
//	@Param		id		path		string	true	"ID del registro de salud"
//	@Success	200	{object}	map[string]string	"message: registro eliminado"
//	@Failure	400	{object}	map[string]string	"id inválido"
//	@Failure	401	{object}	map[string]string	"autenticación requerida"
//	@Failure	404	{object}	map[string]string	"registro no encontrado"
//	@Failure	500	{object}	map[string]string	"mensaje de error"
//	@Router		/api/v1/pets/{pet_id}/health-records/{id} [delete]
func (h *Handler) DeleteHealthRecord(c *gin.Context) {
	petID, ok := parsePetID(c)
	if !ok {
		return
	}
	id, ok := parseRecordID(c)
	if !ok {
		return
	}

	err := h.repo.Delete(c.Request.Context(), id, petID, ownerID(c))
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "registro de salud no encontrado"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al eliminar registro de salud"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registro de salud eliminado"})
}
