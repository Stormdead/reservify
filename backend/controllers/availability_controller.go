package controllers

import (
	"Reservify/dto"
	"Reservify/services"
	"Reservify/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AvailabilityController struct {
	availabilityService *services.AvailabilityService
}

func NewAvailabilityController(availabilityService *services.AvailabilityService) *AvailabilityController {
	return &AvailabilityController{availabilityService: availabilityService}
}

// GetAvailabilityByResource obtiene la disponibilidad de un recurso
// GET /api/resources/:id/availability
func (ctrl *AvailabilityController) GetAvailabilityByResource(c *gin.Context) {
	idParam := c.Param("id")
	resourceID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", err)
		return
	}

	availability, err := ctrl.availabilityService.GetAvailabilityByResource(uint(resourceID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Disponibilidad obtenida exitosamente", availability)
}

// CreateAvailability crea un horario de disponibilidad (solo admin)
// POST /api/admin/resources/:id/availability
func (ctrl *AvailabilityController) CreateAvailability(c *gin.Context) {
	idParam := c.Param("id")
	resourceID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", err)
		return
	}

	var req dto.CreateAvailabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos", err)
		return
	}

	availability, err := ctrl.availabilityService.CreateAvailability(uint(resourceID), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Horario creado exitosamente", availability)
}

// UpdateAvailability actualiza un horario (solo admin)
// PUT /api/admin/availability/:id
func (ctrl *AvailabilityController) UpdateAvailability(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", err)
		return
	}

	var req dto.UpdateAvailabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos", err)
		return
	}

	availability, err := ctrl.availabilityService.UpdateAvailability(uint(id), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Horario actualizado exitosamente", availability)
}

// DeleteAvailability elimina un horario (solo admin)
// DELETE /api/admin/availability/:id
func (ctrl *AvailabilityController) DeleteAvailability(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", err)
		return
	}

	if err := ctrl.availabilityService.DeleteAvailability(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Horario eliminado exitosamente", nil)
}
