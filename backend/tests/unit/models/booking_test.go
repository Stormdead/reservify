package models_test

import (
	"Reservify/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBookingTableName(t *testing.T) {
	booking := models.Booking{}
	tableName := booking.TableName()
	assert.Equal(t, "bookings", tableName, "El nombre de la tabla debería ser 'bookings'")
}

func TestBookingStatus(t *testing.T) {
	// Test: Status constants
	assert.Equal(t, models.BookingStatus("pending"), models.StatusPending)
	assert.Equal(t, models.BookingStatus("confirmed"), models.StatusConfirmed)
	assert.Equal(t, models.BookingStatus("cancelled"), models.StatusCancelled)
	assert.Equal(t, models.BookingStatus("completed"), models.StatusCompleted)
}

func TestBookingModel(t *testing.T) {
	now := time.Now()
	later := now.Add(2 * time.Hour)

	booking := models.Booking{
		UserID:        1,
		ResourceID:    1,
		StartDatetime: now,
		EndDatetime:   later,
		Status:        models.StatusPending,
		TotalPrice:    100.00,
	}

	assert.Equal(t, uint(1), booking.UserID)
	assert.Equal(t, uint(1), booking.ResourceID)
	assert.True(t, booking.StartDatetime.Before(booking.EndDatetime))
	assert.Equal(t, models.StatusPending, booking.Status)
}
