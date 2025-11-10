package models_test

import (
	"Reservify/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserTableName(t *testing.T) {
	user := models.User{}
	tableName := user.TableName()
	assert.Equal(t, "users", tableName, "El nombre de la tabla debería ser 'users'")
}

func TestUserRole(t *testing.T) {
	// Test: Role constants
	assert.Equal(t, models.UserRole("admin"), models.RoleAdmin)
	assert.Equal(t, models.UserRole("user"), models.RoleUser)
}

func TestUserModel(t *testing.T) {
	// Test: Crear un usuario
	user := models.User{
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		FullName:     "Test User",
		Phone:        "1234567890",
		Role:         models.RoleUser,
	}

	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "Test User", user.FullName)
	assert.Equal(t, models.RoleUser, user.Role)
}
