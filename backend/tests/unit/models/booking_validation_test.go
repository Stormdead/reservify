package models_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBookingDateValidation(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		startTime   time.Time
		endTime     time.Time
		shouldError bool
	}{
		{
			name:        "Fechas válidas",
			startTime:   now,
			endTime:     now.Add(2 * time.Hour),
			shouldError: false,
		},
		{
			name:        "Fecha fin antes de inicio",
			startTime:   now.Add(2 * time.Hour),
			endTime:     now,
			shouldError: true,
		},
		{
			name:        "Fechas iguales",
			startTime:   now,
			endTime:     now,
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.startTime.Before(tt.endTime)

			if tt.shouldError {
				assert.False(t, isValid, "Debería ser inválido")
			} else {
				assert.True(t, isValid, "Debería ser válido")
			}
		})
	}
}

func TestBookingOverlapDetection(t *testing.T) {
	// Reserva existente: 10:00 - 12:00
	existingStart := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)
	existingEnd := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		newStart      time.Time
		newEnd        time.Time
		shouldOverlap bool
	}{
		{
			name:          "Antes (sin solapamiento)",
			newStart:      time.Date(2025, 1, 15, 8, 0, 0, 0, time.UTC),
			newEnd:        time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC),
			shouldOverlap: false,
		},
		{
			name:          "Después (sin solapamiento)",
			newStart:      time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
			newEnd:        time.Date(2025, 1, 15, 14, 0, 0, 0, time.UTC),
			shouldOverlap: false,
		},
		{
			name:          "Solapamiento parcial (inicio)",
			newStart:      time.Date(2025, 1, 15, 9, 0, 0, 0, time.UTC),
			newEnd:        time.Date(2025, 1, 15, 11, 0, 0, 0, time.UTC),
			shouldOverlap: true,
		},
		{
			name:          "Solapamiento parcial (fin)",
			newStart:      time.Date(2025, 1, 15, 11, 0, 0, 0, time.UTC),
			newEnd:        time.Date(2025, 1, 15, 13, 0, 0, 0, time.UTC),
			shouldOverlap: true,
		},
		{
			name:          "Contenido dentro",
			newStart:      time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC),
			newEnd:        time.Date(2025, 1, 15, 11, 30, 0, 0, time.UTC),
			shouldOverlap: true,
		},
		{
			name:          "Envolvente",
			newStart:      time.Date(2025, 1, 15, 9, 0, 0, 0, time.UTC),
			newEnd:        time.Date(2025, 1, 15, 13, 0, 0, 0, time.UTC),
			shouldOverlap: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Lógica de detección de solapamiento
			overlaps := (tt.newStart.Before(existingEnd) && tt.newEnd.After(existingStart))

			assert.Equal(t, tt.shouldOverlap, overlaps)
		})
	}
}
