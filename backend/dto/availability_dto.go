package dto

import "time"

// CreateAvailabilityRequest representa los datos para crear disponibilidad
type CreateAvailabilityRequest struct {
	DayOfWeek string `json:"day_of_week" binding:"required,oneof=monday tuesday wednesday thursday friday saturday sunday"`
	StartTime string `json:"start_time" binding:"required"` // Formato: "09:00:00" o "09:00"
	EndTime   string `json:"end_time" binding:"required"`   // Formato: "18:00:00" o "18:00"
}

// UpdateAvailabilityRequest representa los datos para actualizar disponibilidad
type UpdateAvailabilityRequest struct {
	DayOfWeek string `json:"day_of_week" binding:"omitempty,oneof=monday tuesday wednesday thursday friday saturday sunday"`
	StartTime string `json:"start_time" binding:"omitempty"`
	EndTime   string `json:"end_time" binding:"omitempty"`
}

// AvailabilityResponse representa la respuesta de disponibilidad
type AvailabilityResponse struct {
	ID         uint      `json:"id"`
	ResourceID uint      `json:"resource_id"`
	DayOfWeek  string    `json:"day_of_week"`
	StartTime  string    `json:"start_time"`
	EndTime    string    `json:"end_time"`
	CreatedAt  time.Time `json:"created_at"`
}
