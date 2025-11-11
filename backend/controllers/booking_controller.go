package controllers

import (
	"Reservify/dto"
	"Reservify/services"
	"Reservify/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookingController struct {
	bookingService *services.BookingService
}

func NewBookingController(bookingService *services.BookingService) *BookingController {
	return &BookingController{bookingService: bookingService}
}

// GetAllBookings obtiene todas las reservas (solo admin)
// GET /api/admin/bookings?page=1&page_size=10
func (ctrl *BookingController) GetAllBookings(c *gin.Context) {
	params := utils.GetPaginationParams(c)

	bookings, total, err := ctrl.bookingService.GetAllBookings(params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener reservas", err)
		return
	}

	utils.PaginatedSuccessResponse(c, http.StatusOK, "Reservas obtenidas exitosamente", bookings, total, params)
}

// GetMyBookings obtiene las reservas del usuario autenticado
// GET /api/bookings/my?page=1&page_size=10
func (ctrl *BookingController) GetMyBookings(c *gin.Context) {
	userID, _ := c.Get("user_id")
	params := utils.GetPaginationParams(c)

	bookings, total, err := ctrl.bookingService.GetMyBookings(userID.(uint), params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener reservas", err)
		return
	}

	utils.PaginatedSuccessResponse(c, http.StatusOK, "Mis reservas obtenidas exitosamente", bookings, total, params)
}

// GetUpcomingBookings obtiene las próximas reservas del usuario
// GET /api/bookings/upcoming
func (ctrl *BookingController) GetUpcomingBookings(c *gin.Context) {
	userID, _ := c.Get("user_id")

	bookings, err := ctrl.bookingService.GetUpcomingBookings(userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener próximas reservas", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Próximas reservas obtenidas exitosamente", bookings)
}

// GetBookingByID obtiene una reserva por ID
// GET /api/bookings/:id
func (ctrl *BookingController) GetBookingByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", err)
		return
	}

	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("user_role")
	isAdmin := userRole == "admin"

	booking, err := ctrl.bookingService.GetBookingByID(uint(id), userID.(uint), isAdmin)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Reserva obtenida exitosamente", booking)
}

// CreateBooking crea una nueva reserva
// POST /api/bookings
func (ctrl *BookingController) CreateBooking(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req dto.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos", err)
		return
	}

	booking, err := ctrl.bookingService.CreateBooking(userID.(uint), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Reserva creada exitosamente", booking)
}

// UpdateBooking actualiza una reserva
// PUT /api/bookings/:id
func (ctrl *BookingController) UpdateBooking(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", err)
		return
	}

	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("user_role")
	isAdmin := userRole == "admin"

	var req dto.UpdateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos", err)
		return
	}

	booking, err := ctrl.bookingService.UpdateBooking(uint(id), userID.(uint), &req, isAdmin)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Reserva actualizada exitosamente", booking)
}

// CancelBooking cancela una reserva
// DELETE /api/bookings/:id
func (ctrl *BookingController) CancelBooking(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", err)
		return
	}

	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("user_role")
	isAdmin := userRole == "admin"

	if err := ctrl.bookingService.CancelBooking(uint(id), userID.(uint), isAdmin); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Reserva cancelada exitosamente", nil)
}

// ChangeBookingStatus cambia el estado de una reserva (solo admin)
// PATCH /api/admin/bookings/:id/status
func (ctrl *BookingController) ChangeBookingStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", err)
		return
	}

	var req dto.ChangeBookingStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos", err)
		return
	}

	booking, err := ctrl.bookingService.ChangeBookingStatus(uint(id), req.Status)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Estado de reserva actualizado exitosamente", booking)
}

// GetBookingStats obtiene estadísticas de reservas (solo admin)
// GET /api/admin/bookings/stats
func (ctrl *BookingController) GetBookingStats(c *gin.Context) {
	stats, err := ctrl.bookingService.GetBookingStats()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener estadísticas", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Estadísticas obtenidas exitosamente", stats)
}
