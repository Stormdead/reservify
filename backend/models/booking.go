package models

import (
	"time"

	"gorm.io/gorm"
)

type BookingStatus string

const (
	StatusPending   BookingStatus = "pending"
	StatusConfirmed BookingStatus = "confirmed"
	StatusCancelled BookingStatus = "cancelled"
	StatusCompleted BookingStatus = "completed"
)

type Booking struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserID        uint           `gorm:"not null;index" json:"user_id"`
	ResourceID    uint           `gorm:"not null;index" json:"resource_id"`
	StartDatetime time.Time      `gorm:"not null;index:idx_resource_datetime" json:"start_datetime"`
	EndDatetime   time.Time      `gorm:"not null;index:idx_resource_datetime" json:"end_datetime"`
	Status        BookingStatus  `gorm:"type:varchar(20);default:'pending'" json:"status"`
	TotalPrice    float64        `gorm:"type:decimal(10,2)" json:"total_price"`
	Notes         string         `gorm:"type:text" json:"notes"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// Relaciones
	User     User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Resource Resource `gorm:"foreignKey:ResourceID" json:"resource,omitempty"`
}

func (Booking) TableName() string {
	return "bookings"
}

// BeforeCreate hook - validaciones antes de crear
func (b *Booking) BeforeCreate(tx *gorm.DB) error {
	// Validar que la fecha de inicio sea antes de la fecha de fin
	if b.StartDatetime.After(b.EndDatetime) || b.StartDatetime.Equal(b.EndDatetime) {
		return gorm.ErrInvalidData
	}
	return nil
}
