package utils_test

import (
	"Reservify/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "mySecurePassword123"

	// Test: Hash de contraseña
	hash, err := utils.HashPassword(password)
	assert.NoError(t, err, "No debería haber error al generar hash")
	assert.NotEmpty(t, hash, "El hash no debería estar vacío")
	assert.NotEqual(t, password, hash, "El hash no debería ser igual a la contraseña")
}

func TestCheckPassword(t *testing.T) {
	password := "mySecurePassword123"
	wrongPassword := "wrongPassword"

	// Generar hash
	hash, _ := utils.HashPassword(password)

	// Test: Verificar contraseña correcta
	t.Run("Contraseña correcta", func(t *testing.T) {
		result := utils.CheckPassword(password, hash)
		assert.True(t, result, "Debería validar la contraseña correcta")
	})

	// Test: Verificar contraseña incorrecta
	t.Run("Contraseña incorrecta", func(t *testing.T) {
		result := utils.CheckPassword(wrongPassword, hash)
		assert.False(t, result, "No debería validar contraseña incorrecta")
	})
}

func TestHashPasswordConsistency(t *testing.T) {
	password := "testPassword"

	// Test: Generar múltiples hashes de la misma contraseña
	hash1, _ := utils.HashPassword(password)
	hash2, _ := utils.HashPassword(password)

	// Los hashes deberían ser diferentes (bcrypt usa salt aleatorio)
	assert.NotEqual(t, hash1, hash2, "Los hashes deberían ser diferentes debido al salt aleatorio")

	// Pero ambos deberían validar la misma contraseña
	assert.True(t, utils.CheckPassword(password, hash1))
	assert.True(t, utils.CheckPassword(password, hash2))
}
