package models

import (
	"time"

	"gorm.io/gorm"
)

type DayOfWeek string

const (
	Monday    DayOfWeek = "monday"
	Tuesday   DayOfWeek = "tuesday"
	Wednesday DayOfWeek = "wednesday"
	Thursday  DayOfWeek = "thursday"
	Friday    DayOfWeek = "friday"
	Saturday  DayOfWeek = "saturday"
	Sunday    DayOfWeek = "sunday"
)

type AvailabilitySlot struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	ResourceID uint           `gorm:"not null;index" json:"resource_id"`
	DayOfWeek  DayOfWeek      `gorm:"type:varchar(20);not null" json:"day_of_week"`
	StartTime  string         `gorm:"type:time;not null" json:"start_time"` // Formato: "09:00:00"
	EndTime    string         `gorm:"type:time;not null" json:"end_time"`   // Formato: "18:00:00"
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// Relaciones
	Resource Resource `gorm:"foreignKey:ResourceID" json:"resource,omitempty"`
}

func (AvailabilitySlot) TableName() string {
	return "availability_slots"
}
