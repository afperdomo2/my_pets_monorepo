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

// parsePetIDParam extrae el parámetro de ruta :pet_id de las rutas /health-records/pets/:pet_id.
func parsePetIDParam(c *gin.Context) (string, bool) {
	id := c.Param("pet_id")
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

// GetAllHealthRecords maneja GET /api/v1/health-records
// Lista todos los registros de salud de todas las mascotas del usuario autenticado.
//
//	@Summary	Listar todos los registros de salud del usuario (paginado)
//	@Tags		health-records
//	@Produce	json
//	@Security	CookieAuth
//	@Param		page		query		int						false	"Número de página (por defecto 1)"
//	@Param		per_page	query		int						false	"Elementos por página (por defecto 10)"
//	@Success	200	{object}	map[string]interface{}	"data: []HealthRecord, total: int, page: int, per_page: int, total_pages: int"
//	@Failure	401	{object}	map[string]string		"autenticación requerida"
//	@Failure	500	{object}	map[string]string		"mensaje de error"
//	@Router		/api/v1/health-records [get]
func (h *Handler) GetAllHealthRecords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	records, total, err := h.repo.GetAllByOwner(c.Request.Context(), ownerID(c), page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener registros de salud"})
		return
	}

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

// GetHealthRecordsByPet maneja GET /api/v1/health-records/pets/:pet_id
// Lista todos los registros de salud de una mascota específica del usuario autenticado.
//
//	@Summary	Listar registros de salud de una mascota (paginado)
//	@Tags		health-records
//	@Produce	json
//	@Security	CookieAuth
//	@Param		pet_id		path		string	true	"ID de la mascota"
//	@Param		page		query		int		false	"Número de página (por defecto 1)"
//	@Param		per_page	query		int		false	"Elementos por página (por defecto 10)"
//	@Success	200	{object}	map[string]interface{}	"data: []HealthRecord, total: int, page: int, per_page: int, total_pages: int"
//	@Failure	400	{object}	map[string]string		"pet_id inválido"
//	@Failure	401	{object}	map[string]string		"autenticación requerida"
//	@Failure	500	{object}	map[string]string		"mensaje de error"
//	@Router		/api/v1/health-records/pets/{pet_id} [get]
func (h *Handler) GetHealthRecordsByPet(c *gin.Context) {
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

	records, total, err := h.repo.GetByPetID(c.Request.Context(), petID, ownerID(c), page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener registros de salud"})
		return
	}

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

// GetHealthRecordsByPetAndCategory maneja GET /api/v1/health-records/pets/:pet_id/category/:category
// Lista los registros de salud de una mascota filtrados por categoría (vaccine, deworming, exam).
//
//	@Summary	Listar registros de salud de una mascota por categoría (paginado)
//	@Tags		health-records
//	@Produce	json
//	@Security	CookieAuth
//	@Param		pet_id		path		string	true	"ID de la mascota"
//	@Param		category	path		string	true	"Categoría (vaccine, deworming, exam)"
//	@Param		page		query		int		false	"Número de página (por defecto 1)"
//	@Param		per_page	query		int		false	"Elementos por página (por defecto 10)"
//	@Success	200	{object}	map[string]interface{}	"data: []HealthRecord, total: int, page: int, per_page: int, total_pages: int"
//	@Failure	400	{object}	map[string]string		"pet_id o categoría inválida"
//	@Failure	401	{object}	map[string]string		"autenticación requerida"
//	@Failure	500	{object}	map[string]string		"mensaje de error"
//	@Router		/api/v1/health-records/pets/{pet_id}/category/{category} [get]
func (h *Handler) GetHealthRecordsByPetAndCategory(c *gin.Context) {
	petID, ok := parsePetIDParam(c)
	if !ok {
		return
	}

	category := c.Param("category")
	validCategories := map[string]bool{"vaccine": true, "deworming": true, "exam": true}
	if !validCategories[category] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "categoría inválida: debe ser vaccine, deworming o exam"})
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

	records, total, err := h.repo.GetByPetIDAndCategory(c.Request.Context(), petID, category, ownerID(c), page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener registros de salud"})
		return
	}

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

// CreateHealthRecord maneja POST /api/v1/health-records
// El pet_id se recibe en el body. Si se provee health_catalog_id,
// copia name y category desde el catálogo automáticamente.
//
//	@Summary	Crear un registro de salud
//	@Tags		health-records
//	@Accept		json
//	@Produce	json
//	@Security	CookieAuth
//	@Param		record	body		CreateHealthRecordPayload	true	"Datos del registro (incluye pet_id)"
//	@Success	201	{object}	map[string]interface{}	"data: HealthRecord"
//	@Failure	400	{object}	map[string]string		"error de validación"
//	@Failure	401	{object}	map[string]string		"autenticación requerida"
//	@Failure	404	{object}	map[string]string		"mascota o catálogo no encontrado"
//	@Failure	500	{object}	map[string]string		"mensaje de error"
//	@Router		/api/v1/health-records [post]
func (h *Handler) CreateHealthRecord(c *gin.Context) {
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

	created, err := h.repo.Create(c.Request.Context(), ownerID(c), payload)
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "mascota no encontrada o no pertenece al usuario"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al crear registro de salud"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": resolveOverdue(created)})
}

// UpdateHealthRecord maneja PUT /api/v1/health-records/:record_id
// No permite cambiar el pet_id ni el health_catalog_id del registro.
//
//	@Summary	Actualizar un registro de salud
//	@Tags		health-records
//	@Accept		json
//	@Produce	json
//	@Security	CookieAuth
//	@Param		record_id	path		string						true	"ID del registro de salud"
//	@Param		record		body		UpdateHealthRecordPayload	true	"Datos del registro"
//	@Success	200	{object}	map[string]interface{}	"data: HealthRecord"
//	@Failure	400	{object}	map[string]string		"error de validación"
//	@Failure	401	{object}	map[string]string		"autenticación requerida"
//	@Failure	404	{object}	map[string]string		"registro no encontrado"
//	@Failure	500	{object}	map[string]string		"mensaje de error"
//	@Router		/api/v1/health-records/{record_id} [put]
func (h *Handler) UpdateHealthRecord(c *gin.Context) {
	id, ok := parseRecordID(c)
	if !ok {
		return
	}

	var payload UpdateHealthRecordPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validation.Translate(err)})
		return
	}

	updated, err := h.repo.Update(c.Request.Context(), id, ownerID(c), payload)
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

// UpdateHealthRecordStatus maneja PATCH /api/v1/health-records/:record_id/status
//
//	@Summary	Actualizar el status de un registro de salud
//	@Tags		health-records
//	@Accept		json
//	@Produce	json
//	@Security	CookieAuth
//	@Param		record_id	path		string				true	"ID del registro de salud"
//	@Param		status		body		UpdateStatusPayload	true	"Nuevo status"
//	@Success	200	{object}	map[string]interface{}	"data: HealthRecord"
//	@Failure	400	{object}	map[string]string		"error de validación"
//	@Failure	401	{object}	map[string]string		"autenticación requerida"
//	@Failure	404	{object}	map[string]string		"registro no encontrado"
//	@Failure	500	{object}	map[string]string		"mensaje de error"
//	@Router		/api/v1/health-records/{record_id}/status [patch]
func (h *Handler) UpdateHealthRecordStatus(c *gin.Context) {
	id, ok := parseRecordID(c)
	if !ok {
		return
	}

	var payload UpdateStatusPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validation.Translate(err)})
		return
	}

	updated, err := h.repo.UpdateStatus(c.Request.Context(), id, ownerID(c), payload)
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

// DeleteHealthRecord maneja DELETE /api/v1/health-records/:record_id
//
//	@Summary	Eliminar un registro de salud
//	@Tags		health-records
//	@Produce	json
//	@Security	CookieAuth
//	@Param		record_id	path		string	true	"ID del registro de salud"
//	@Success	200	{object}	map[string]string	"message: registro eliminado"
//	@Failure	400	{object}	map[string]string	"id inválido"
//	@Failure	401	{object}	map[string]string	"autenticación requerida"
//	@Failure	404	{object}	map[string]string	"registro no encontrado"
//	@Failure	500	{object}	map[string]string	"mensaje de error"
//	@Router		/api/v1/health-records/{record_id} [delete]
func (h *Handler) DeleteHealthRecord(c *gin.Context) {
	id, ok := parseRecordID(c)
	if !ok {
		return
	}

	err := h.repo.Delete(c.Request.Context(), id, ownerID(c))
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

// GetUpcomingRecords maneja GET /api/v1/health-records/upcoming
// Lista los próximos registros pendientes de aplicación con paginación estándar.
// Los registros se ordenan por due_date ASC (los más próximos primero).
//
//	@Summary	Listar próximos registros pendientes
//	@Tags		health-records
//	@Produce	json
//	@Security	CookieAuth
//	@Param		page		query		int						false	"Número de página (default: 1)"
//	@Param		per_page	query		int						false	"Registros por página (default: 10, max: 50)"
//	@Param		category	query		string					false	"Filtrar por categoría (vaccine, deworming, exam)"
//	@Success	200		{object}	map[string]interface{}	"data, total, page, per_page, total_pages"
//	@Failure	400		{object}	map[string]string		"categoría inválida"
//	@Failure	401		{object}	map[string]string		"autenticación requerida"
//	@Failure	500		{object}	map[string]string		"mensaje de error"
//	@Router		/api/v1/health-records/upcoming [get]
func (h *Handler) GetUpcomingRecords(c *gin.Context) {
	// Parsear y validar page
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	// Parsear y validar per_page
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if perPage < 1 {
		perPage = 10
	}
	if perPage > 50 {
		perPage = 50
	}

	// Obtener categoría (opcional)
	category := c.Query("category")
	if category != "" {
		validCategories := map[string]bool{"vaccine": true, "deworming": true, "exam": true}
		if !validCategories[category] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "categoría inválida: debe ser vaccine, deworming o exam"})
			return
		}
	}

	// Obtener próximos registros pendientes del usuario con paginación
	records, total, err := h.repo.GetUpcomingRecords(c.Request.Context(), ownerID(c), category, page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener próximos registros"})
		return
	}

	// Resolver estado overdue para cada registro
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
