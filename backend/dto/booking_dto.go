package dto

import "time"

// CreateBookingRequest representa los datos para crear una reserva
type CreateBookingRequest struct {
	ResourceID    uint      `json:"resource_id" binding:"required"`
	StartDatetime time.Time `json:"start_datetime" binding:"required"`
	EndDatetime   time.Time `json:"end_datetime" binding:"required"`
	Notes         string    `json:"notes"`
}

// UpdateBookingRequest representa los datos para actualizar una reserva
type UpdateBookingRequest struct {
	StartDatetime time.Time `json:"start_datetime" binding:"required"`
	EndDatetime   time.Time `json:"end_datetime" binding:"required"`
	Notes         string    `json:"notes"`
}

// BookingResponse representa la respuesta de una reserva
type BookingResponse struct {
	ID            uint             `json:"id"`
	UserID        uint             `json:"user_id"`
	User          UserResponse     `json:"user"`
	ResourceID    uint             `json:"resource_id"`
	Resource      ResourceResponse `json:"resource"`
	StartDatetime time.Time        `json:"start_datetime"`
	EndDatetime   time.Time        `json:"end_datetime"`
	Status        string           `json:"status"`
	TotalPrice    float64          `json:"total_price"`
	Notes         string           `json:"notes"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
}

// BookingListResponse representa una reserva en la lista (más ligero)
type BookingListResponse struct {
	ID            uint      `json:"id"`
	UserID        uint      `json:"user_id"`
	UserName      string    `json:"user_name"`
	ResourceID    uint      `json:"resource_id"`
	ResourceName  string    `json:"resource_name"`
	StartDatetime time.Time `json:"start_datetime"`
	EndDatetime   time.Time `json:"end_datetime"`
	Status        string    `json:"status"`
	TotalPrice    float64   `json:"total_price"`
	CreatedAt     time.Time `json:"created_at"`
}

// ChangeBookingStatusRequest representa el cambio de estado de una reserva
type ChangeBookingStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending confirmed cancelled completed"`
}

// BookingStatsResponse representa estadísticas de reservas
type BookingStatsResponse struct {
	TotalBookings     int64   `json:"total_bookings"`
	PendingBookings   int64   `json:"pending_bookings"`
	ConfirmedBookings int64   `json:"confirmed_bookings"`
	CancelledBookings int64   `json:"cancelled_bookings"`
	CompletedBookings int64   `json:"completed_bookings"`
	TotalRevenue      float64 `json:"total_revenue"`
	PendingRevenue    float64 `json:"pending_revenue"`
	ConfirmedRevenue  float64 `json:"confirmed_revenue"`
}
