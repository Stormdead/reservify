package dto

import "time"

// CreateResourceRequest representa los datos para crear un recurso
type CreateResourceRequest struct {
	Name         string  `json:"name" binding:"required,min=3"`
	Description  string  `json:"description"`
	Capacity     int     `json:"capacity" binding:"required,min=1"`
	PricePerHour float64 `json:"price_per_hour" binding:"required,min=0"`
	Category     string  `json:"category"`
	ImageURL     string  `json:"image_url"`
}

// UpdateResourceRequest representa los datos para actualizar un recurso
type UpdateResourceRequest struct {
	Name         string  `json:"name" binding:"omitempty,min=3"`
	Description  string  `json:"description"`
	Capacity     int     `json:"capacity" binding:"omitempty,min=1"`
	PricePerHour float64 `json:"price_per_hour" binding:"omitempty,min=0"`
	Category     string  `json:"category"`
	ImageURL     string  `json:"image_url"`
	IsActive     *bool   `json:"is_active"` // Pointer para permitir false
}

// ResourceResponse representa la respuesta de un recurso
type ResourceResponse struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Capacity     int       `json:"capacity"`
	PricePerHour float64   `json:"price_per_hour"`
	Category     string    `json:"category"`
	ImageURL     string    `json:"image_url"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ResourceListResponse representa un recurso en la lista (más ligero)
type ResourceListResponse struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name"`
	Capacity     int     `json:"capacity"`
	PricePerHour float64 `json:"price_per_hour"`
	Category     string  `json:"category"`
	ImageURL     string  `json:"image_url"`
	IsActive     bool    `json:"is_active"`
}
