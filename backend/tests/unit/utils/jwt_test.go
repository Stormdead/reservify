package utils_test

import (
	"Reservify/config"
	"Reservify/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	// Setup: Configurar config para testing
	config.AppConfig = &config.Config{
		JWTSecret: "test-secret-key-for-testing",
	}
}

func TestGenerateToken(t *testing.T) {
	userID := uint(1)
	email := "test@example.com"
	role := "user"

	// Test: Generar token
	token, err := utils.GenerateToken(userID, email, role)
	assert.NoError(t, err, "No debería haber error al generar token")
	assert.NotEmpty(t, token, "El token no debería estar vacío")
}

func TestValidateToken(t *testing.T) {
	userID := uint(1)
	email := "test@example.com"
	role := "user"

	// Generar token
	token, _ := utils.GenerateToken(userID, email, role)

	// Test: Validar token correcto
	t.Run("Token válido", func(t *testing.T) {
		claims, err := utils.ValidateToken(token)
		assert.NoError(t, err, "No debería haber error al validar token válido")
		assert.NotNil(t, claims, "Claims no deberían ser nil")
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, email, claims.Email)
		assert.Equal(t, role, claims.Role)
	})

	// Test: Validar token inválido
	t.Run("Token inválido", func(t *testing.T) {
		invalidToken := "invalid.token.here"
		claims, err := utils.ValidateToken(invalidToken)
		assert.Error(t, err, "Debería haber error con token inválido")
		assert.Nil(t, claims, "Claims deberían ser nil")
	})

	// Test: Validar token vacío
	t.Run("Token vacío", func(t *testing.T) {
		claims, err := utils.ValidateToken("")
		assert.Error(t, err, "Debería haber error con token vacío")
		assert.Nil(t, claims, "Claims deberían ser nil")
	})
}

func TestTokenExpiration(t *testing.T) {
	userID := uint(1)
	email := "test@example.com"
	role := "admin"

	token, _ := utils.GenerateToken(userID, email, role)
	claims, _ := utils.ValidateToken(token)

	// Verificar que ExpiresAt está en el futuro
	assert.True(t, claims.ExpiresAt.Time.After(time.Now()), "ExpiresAt debería estar en el futuro")
}
