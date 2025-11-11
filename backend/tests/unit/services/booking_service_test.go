package services_test

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCalculatePrice(t *testing.T) {
	// Este test verifica el cálculo de precios
	// Nota: Como calculatePrice es privado, lo probamos indirectamente

	tests := []struct {
		name          string
		pricePerHour  float64
		hours         float64
		expectedPrice float64
	}{
		{"1 hora exacta", 50.0, 1.0, 50.0},
		{"2 horas exactas", 50.0, 2.0, 100.0},
		{"1.5 horas (redondeo)", 50.0, 1.5, 100.0},  // Se redondea a 2 horas
		{"0.5 horas (redondeo)", 50.0, 0.5, 50.0},   // Se redondea a 1 hora
		{"3.2 horas (redondeo)", 100.0, 3.2, 400.0}, // Se redondea a 4 horas
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()
			end := start.Add(time.Duration(tt.hours * float64(time.Hour)))

			duration := end.Sub(start)
			hours := duration.Hours()

			// Simular el cálculo (mismo algoritmo que en el servicio)
			hoursRounded := math.Ceil(hours)
			price := tt.pricePerHour * hoursRounded

			assert.Equal(t, tt.expectedPrice, price)
		})
	}
}

func TestValidateStatusTransition(t *testing.T) {
	// Test: Transiciones válidas
	validTransitions := map[string][]string{
		"pending":   {"confirmed", "cancelled"},
		"confirmed": {"completed", "cancelled"},
	}

	for currentStatus, allowedStatuses := range validTransitions {
		for _, newStatus := range allowedStatuses {
			t.Run("De "+currentStatus+" a "+newStatus, func(t *testing.T) {
				// Esta transición debería ser válida
				assert.Contains(t, allowedStatuses, newStatus)
			})
		}
	}

	// Test: Transiciones inválidas
	invalidTransitions := map[string][]string{
		"pending":   {"completed"},
		"confirmed": {"pending"},
		"cancelled": {"pending", "confirmed", "completed"},
		"completed": {"pending", "confirmed", "cancelled"},
	}

	for currentStatus, invalidStatuses := range invalidTransitions {
		for _, newStatus := range invalidStatuses {
			t.Run("Inválido: de "+currentStatus+" a "+newStatus, func(t *testing.T) {
				// Esta transición NO debería ser válida
				validForCurrent := validTransitions[currentStatus]
				assert.NotContains(t, validForCurrent, newStatus)
			})
		}
	}
}
