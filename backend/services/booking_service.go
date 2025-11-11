package services

import (
	"Reservify/dto"
	"Reservify/models"
	"Reservify/repositories"
	"Reservify/utils"
	"errors"
	"fmt"
	"math"
	"time"
)

type BookingService struct {
	bookingRepo  *repositories.BookingRepository
	resourceRepo *repositories.ResourceRepository
	userRepo     *repositories.UserRepository
}

func NewBookingService(
	bookingRepo *repositories.BookingRepository,
	resourceRepo *repositories.ResourceRepository,
	userRepo *repositories.UserRepository,
) *BookingService {
	return &BookingService{
		bookingRepo:  bookingRepo,
		resourceRepo: resourceRepo,
		userRepo:     userRepo,
	}
}

// GetAllBookings obtiene todas las reservas (admin)
func (s *BookingService) GetAllBookings(params utils.PaginationParams) ([]dto.BookingListResponse, int64, error) {
	bookings, total, err := s.bookingRepo.FindAll(params)
	if err != nil {
		return nil, 0, err
	}

	return s.mapToListResponse(bookings), total, nil
}

// GetMyBookings obtiene las reservas del usuario autenticado
func (s *BookingService) GetMyBookings(userID uint, params utils.PaginationParams) ([]dto.BookingListResponse, int64, error) {
	bookings, total, err := s.bookingRepo.FindByUserID(userID, params)
	if err != nil {
		return nil, 0, err
	}

	return s.mapToListResponse(bookings), total, nil
}

// GetBookingByID obtiene una reserva por ID
func (s *BookingService) GetBookingByID(id uint, userID uint, isAdmin bool) (*dto.BookingResponse, error) {
	booking, err := s.bookingRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Verificar permisos: solo el dueño o admin pueden ver
	if !isAdmin && booking.UserID != userID {
		return nil, errors.New("no tienes permisos para ver esta reserva")
	}

	return s.mapToResponse(booking), nil
}

// CreateBooking crea una nueva reserva
func (s *BookingService) CreateBooking(userID uint, req *dto.CreateBookingRequest) (*dto.BookingResponse, error) {
	// Validar que start_datetime < end_datetime
	if req.StartDatetime.After(req.EndDatetime) || req.StartDatetime.Equal(req.EndDatetime) {
		return nil, errors.New("la fecha de inicio debe ser anterior a la fecha de fin")
	}

	// Validar que la reserva sea en el futuro
	if req.StartDatetime.Before(time.Now()) {
		return nil, errors.New("no se pueden crear reservas en el pasado")
	}

	// Verificar que el recurso existe y está activo
	resource, err := s.resourceRepo.FindByID(req.ResourceID)
	if err != nil {
		return nil, err
	}
	if !resource.IsActive {
		return nil, errors.New("el recurso no está disponible")
	}

	// Verificar disponibilidad (no hay solapamiento)
	overlap, err := s.bookingRepo.CheckOverlap(req.ResourceID, req.StartDatetime, req.EndDatetime, nil)
	if err != nil {
		return nil, err
	}
	if overlap {
		return nil, errors.New("el recurso no está disponible en ese horario")
	}

	// Calcular precio total
	totalPrice := s.calculatePrice(resource.PricePerHour, req.StartDatetime, req.EndDatetime)

	// Crear la reserva
	booking := &models.Booking{
		UserID:        userID,
		ResourceID:    req.ResourceID,
		StartDatetime: req.StartDatetime,
		EndDatetime:   req.EndDatetime,
		Status:        models.StatusPending,
		TotalPrice:    totalPrice,
		Notes:         req.Notes,
	}

	if err := s.bookingRepo.Create(booking); err != nil {
		return nil, errors.New("error al crear la reserva")
	}

	// Obtener la reserva completa con relaciones
	return s.GetBookingByID(booking.ID, userID, false)
}

// UpdateBooking actualiza una reserva
func (s *BookingService) UpdateBooking(id uint, userID uint, req *dto.UpdateBookingRequest, isAdmin bool) (*dto.BookingResponse, error) {
	// Buscar la reserva
	booking, err := s.bookingRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Verificar permisos
	if !isAdmin && booking.UserID != userID {
		return nil, errors.New("no tienes permisos para editar esta reserva")
	}

	// Solo se pueden editar reservas pendientes
	if booking.Status != models.StatusPending {
		return nil, errors.New("solo se pueden editar reservas pendientes")
	}

	// Validar fechas
	if req.StartDatetime.After(req.EndDatetime) || req.StartDatetime.Equal(req.EndDatetime) {
		return nil, errors.New("la fecha de inicio debe ser anterior a la fecha de fin")
	}

	if req.StartDatetime.Before(time.Now()) {
		return nil, errors.New("no se pueden crear reservas en el pasado")
	}

	// Verificar disponibilidad (excluyendo esta reserva)
	overlap, err := s.bookingRepo.CheckOverlap(booking.ResourceID, req.StartDatetime, req.EndDatetime, &id)
	if err != nil {
		return nil, err
	}
	if overlap {
		return nil, errors.New("el recurso no está disponible en ese horario")
	}

	// Obtener recurso para recalcular precio
	resource, _ := s.resourceRepo.FindByID(booking.ResourceID)

	// Actualizar campos
	booking.StartDatetime = req.StartDatetime
	booking.EndDatetime = req.EndDatetime
	booking.Notes = req.Notes
	booking.TotalPrice = s.calculatePrice(resource.PricePerHour, req.StartDatetime, req.EndDatetime)

	if err := s.bookingRepo.Update(booking); err != nil {
		return nil, errors.New("error al actualizar la reserva")
	}

	return s.GetBookingByID(booking.ID, userID, isAdmin)
}

// CancelBooking cancela una reserva
func (s *BookingService) CancelBooking(id uint, userID uint, isAdmin bool) error {
	booking, err := s.bookingRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Verificar permisos
	if !isAdmin && booking.UserID != userID {
		return errors.New("no tienes permisos para cancelar esta reserva")
	}

	// Solo se pueden cancelar reservas pending o confirmed
	if booking.Status != models.StatusPending && booking.Status != models.StatusConfirmed {
		return errors.New("solo se pueden cancelar reservas pendientes o confirmadas")
	}

	// Cambiar estado a cancelled
	booking.Status = models.StatusCancelled
	return s.bookingRepo.Update(booking)
}

// ChangeBookingStatus cambia el estado de una reserva (solo admin)
func (s *BookingService) ChangeBookingStatus(id uint, status string) (*dto.BookingResponse, error) {
	booking, err := s.bookingRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Validar transiciones de estado
	if err := s.validateStatusTransition(booking.Status, models.BookingStatus(status)); err != nil {
		return nil, err
	}

	booking.Status = models.BookingStatus(status)
	if err := s.bookingRepo.Update(booking); err != nil {
		return nil, errors.New("error al cambiar el estado")
	}

	return s.mapToResponse(booking), nil
}

// GetBookingStats obtiene estadísticas de reservas (admin)
func (s *BookingService) GetBookingStats() (*dto.BookingStatsResponse, error) {
	pending, _ := s.bookingRepo.CountByStatus("pending")
	confirmed, _ := s.bookingRepo.CountByStatus("confirmed")
	cancelled, _ := s.bookingRepo.CountByStatus("cancelled")
	completed, _ := s.bookingRepo.CountByStatus("completed")

	totalRevenue, _ := s.bookingRepo.GetTotalRevenue()
	pendingRevenue, _ := s.bookingRepo.GetRevenueByStatus("pending")
	confirmedRevenue, _ := s.bookingRepo.GetRevenueByStatus("confirmed")

	stats := &dto.BookingStatsResponse{
		TotalBookings:     pending + confirmed + cancelled + completed,
		PendingBookings:   pending,
		ConfirmedBookings: confirmed,
		CancelledBookings: cancelled,
		CompletedBookings: completed,
		TotalRevenue:      totalRevenue,
		PendingRevenue:    pendingRevenue,
		ConfirmedRevenue:  confirmedRevenue,
	}

	return stats, nil
}

// GetUpcomingBookings obtiene las próximas reservas del usuario
func (s *BookingService) GetUpcomingBookings(userID uint) ([]dto.BookingListResponse, error) {
	bookings, err := s.bookingRepo.GetUpcomingBookings(userID, 5)
	if err != nil {
		return nil, err
	}
	return s.mapToListResponse(bookings), nil
}

// ============= FUNCIONES AUXILIARES =============

// calculatePrice calcula el precio total basado en las horas
func (s *BookingService) calculatePrice(pricePerHour float64, start, end time.Time) float64 {
	duration := end.Sub(start)
	hours := duration.Hours()

	// Redondear hacia arriba (si reservas 1.5 horas, pagas 2 horas)
	hoursRounded := math.Ceil(hours)

	return pricePerHour * hoursRounded
}

// validateStatusTransition valida las transiciones de estado permitidas
func (s *BookingService) validateStatusTransition(currentStatus, newStatus models.BookingStatus) error {
	validTransitions := map[models.BookingStatus][]models.BookingStatus{
		models.StatusPending:   {models.StatusConfirmed, models.StatusCancelled},
		models.StatusConfirmed: {models.StatusCompleted, models.StatusCancelled},
		models.StatusCancelled: {}, // No se puede cambiar desde cancelled
		models.StatusCompleted: {}, // No se puede cambiar desde completed
	}

	allowed, exists := validTransitions[currentStatus]
	if !exists {
		return fmt.Errorf("estado actual inválido: %s", currentStatus)
	}

	for _, status := range allowed {
		if status == newStatus {
			return nil
		}
	}

	return fmt.Errorf("no se puede cambiar de %s a %s", currentStatus, newStatus)
}

// mapToResponse convierte un modelo a DTO completo
func (s *BookingService) mapToResponse(booking *models.Booking) *dto.BookingResponse {
	return &dto.BookingResponse{
		ID:     booking.ID,
		UserID: booking.UserID,
		User: dto.UserResponse{
			ID:        booking.User.ID,
			Email:     booking.User.Email,
			FullName:  booking.User.FullName,
			Phone:     booking.User.Phone,
			Role:      string(booking.User.Role),
			CreatedAt: booking.User.CreatedAt,
		},
		ResourceID: booking.ResourceID,
		Resource: dto.ResourceResponse{
			ID:           booking.Resource.ID,
			Name:         booking.Resource.Name,
			Description:  booking.Resource.Description,
			Capacity:     booking.Resource.Capacity,
			PricePerHour: booking.Resource.PricePerHour,
			Category:     booking.Resource.Category,
			ImageURL:     booking.Resource.ImageURL,
			IsActive:     booking.Resource.IsActive,
			CreatedAt:    booking.Resource.CreatedAt,
			UpdatedAt:    booking.Resource.UpdatedAt,
		},
		StartDatetime: booking.StartDatetime,
		EndDatetime:   booking.EndDatetime,
		Status:        string(booking.Status),
		TotalPrice:    booking.TotalPrice,
		Notes:         booking.Notes,
		CreatedAt:     booking.CreatedAt,
		UpdatedAt:     booking.UpdatedAt,
	}
}

// mapToListResponse convierte múltiples modelos a DTOs ligeros
func (s *BookingService) mapToListResponse(bookings []models.Booking) []dto.BookingListResponse {
	var response []dto.BookingListResponse
	for _, booking := range bookings {
		response = append(response, dto.BookingListResponse{
			ID:            booking.ID,
			UserID:        booking.UserID,
			UserName:      booking.User.FullName,
			ResourceID:    booking.ResourceID,
			ResourceName:  booking.Resource.Name,
			StartDatetime: booking.StartDatetime,
			EndDatetime:   booking.EndDatetime,
			Status:        string(booking.Status),
			TotalPrice:    booking.TotalPrice,
			CreatedAt:     booking.CreatedAt,
		})
	}
	return response
}
