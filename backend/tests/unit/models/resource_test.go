package models_test

import (
	"Reservify/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResourceTableName(t *testing.T) {
	resource := models.Resource{}
	tableName := resource.TableName()
	assert.Equal(t, "resources", tableName, "El nombre de la tabla debería ser 'resources'")
}

func TestResourceModel(t *testing.T) {
	// Test: Crear un recurso
	resource := models.Resource{
		Name:         "Sala A",
		Description:  "Sala de reuniones",
		Capacity:     10,
		PricePerHour: 50.00,
		Category:     "Salas",
		IsActive:     true,
	}

	assert.Equal(t, "Sala A", resource.Name)
	assert.Equal(t, 10, resource.Capacity)
	assert.Equal(t, 50.00, resource.PricePerHour)
	assert.True(t, resource.IsActive)
}
