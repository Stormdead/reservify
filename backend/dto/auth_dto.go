package dto

// Representa los datos de registro
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required,min=3"`
	Phone    string `json:"phone"`
}

// Representa los datos de login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Representa la respuesta de autenticación
type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
