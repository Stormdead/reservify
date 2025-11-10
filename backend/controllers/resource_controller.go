package controllers

import (
	"Reservify/dto"
	"Reservify/services"
	"Reservify/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ResourceController struct {
	resourceService *services.ResourceService
}

func NewResourceController(resourceService *services.ResourceService) *ResourceController {
	return &ResourceController{resourceService: resourceService}
}

// GetAllResources obtiene todos los recursos (público - solo activos, admin - todos)
// GET /api/resources?page=1&page_size=10&search=sala
func (ctrl *ResourceController) GetAllResources(c *gin.Context) {
	params := utils.GetPaginationParams(c)

	// Verificar si es admin
	userRole, exists := c.Get("user_role")
	activeOnly := true
	if exists && userRole == "admin" {
		// Admin puede ver todos (activos e inactivos)
		activeOnly = false
	}

	resources, total, err := ctrl.resourceService.GetAllResources(params, activeOnly)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener recursos", err)
		return
	}

	utils.PaginatedSuccessResponse(c, http.StatusOK, "Recursos obtenidos exitosamente", resources, total, params)
}

// GetResourcesByCategory obtiene recursos por categoría
// GET /api/resources/category/:category?page=1&page_size=10
func (ctrl *ResourceController) GetResourcesByCategory(c *gin.Context) {
	category := c.Param("category")
	params := utils.GetPaginationParams(c)

	resources, total, err := ctrl.resourceService.GetResourcesByCategory(category, params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener recursos", err)
		return
	}

	utils.PaginatedSuccessResponse(c, http.StatusOK, "Recursos obtenidos exitosamente", resources, total, params)
}

// GetResourceByID obtiene un recurso por ID
// GET /api/resources/:id
func (ctrl *ResourceController) GetResourceByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", err)
		return
	}

	resource, err := ctrl.resourceService.GetResourceByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Recurso obtenido exitosamente", resource)
}

// CreateResource crea un nuevo recurso (solo admin)
// POST /api/admin/resources
func (ctrl *ResourceController) CreateResource(c *gin.Context) {
	var req dto.CreateResourceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos", err)
		return
	}

	resource, err := ctrl.resourceService.CreateResource(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Recurso creado exitosamente", resource)
}

// UpdateResource actualiza un recurso (solo admin)
// PUT /api/admin/resources/:id
func (ctrl *ResourceController) UpdateResource(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", err)
		return
	}

	var req dto.UpdateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos", err)
		return
	}

	resource, err := ctrl.resourceService.UpdateResource(uint(id), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Recurso actualizado exitosamente", resource)
}

// DeleteResource elimina un recurso (solo admin)
// DELETE /api/admin/resources/:id
func (ctrl *ResourceController) DeleteResource(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", err)
		return
	}

	if err := ctrl.resourceService.DeleteResource(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Recurso eliminado exitosamente", nil)
}

// GetCategories obtiene todas las categorías
// GET /api/resources/categories
func (ctrl *ResourceController) GetCategories(c *gin.Context) {
	categories, err := ctrl.resourceService.GetCategories()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener categorías", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Categorías obtenidas exitosamente", gin.H{
		"categories": categories,
	})
}

// GetResourceStats obtiene estadísticas de recursos (solo admin)
// GET /api/admin/resources/stats
func (ctrl *ResourceController) GetResourceStats(c *gin.Context) {
	stats, err := ctrl.resourceService.GetResourceStats()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener estadísticas", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Estadísticas obtenidas exitosamente", stats)
}
