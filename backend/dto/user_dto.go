package dto

import "time"

// Representa la respuesta de usuario (sin password)
type UserResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Phone     string    `json:"phone"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// Representa los datos para actualizar usuario
type UpdateUserRequest struct {
	FullName string `json:"full_name" binding:"omitempty,min=3"`
	Phone    string `json:"phone"`
}

// Representa los datos para cambiar contraseña
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}
