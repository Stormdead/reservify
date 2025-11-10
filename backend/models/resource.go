package models

import (
	"time"

	"gorm.io/gorm"
)

type Resource struct {
	ID                uint               `gorm:"primaryKey" json:"id"`
	Name              string             `gorm:"not null;size:255" json:"name"`
	Description       string             `gorm:"type:text" json:"description"`
	Capacity          int                `gorm:"not null" json:"capacity"`
	PricePerHour      float64            `gorm:"type:decimal(10,2)" json:"price_per_hour"`
	Category          string             `gorm:"size:100" json:"category"`
	ImageURL          string             `gorm:"size:500" json:"image_url"`
	IsActive          bool               `gorm:"default:true" json:"is_active"`
	CreatedAt         time.Time          `json:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at"`
	DeletedAt         gorm.DeletedAt     `gorm:"index" json:"-"`
	AvailabilitySlots []AvailabilitySlot `gorm:"foreignKey:ResourceID;constraint:OnDelete:CASCADE" json:"availability_slots,omitempty"`
	Bookings          []Booking          `gorm:"foreignKey:ResourceID" json:"bookings,omitempty"`
}

func (Resource) TableName() string {
	return "resources"
}
