package controllers

import (
	"net/http"

	"Reservify/dto"
	"Reservify/services"
	"Reservify/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// Register maneja el registro de nuevos usuarios
// POST /api/auth/register
func (ctrl *AuthController) Register(c *gin.Context) {
	var req dto.RegisterRequest

	// Validar el request
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos", err)
		return
	}

	// Llamar al servicio
	response, err := ctrl.authService.Register(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Usuario registrado exitosamente", response)
}

// Login maneja el inicio de sesión
// POST /api/auth/login
func (ctrl *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest

	// Validar el request
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos", err)
		return
	}

	// Llamar al servicio
	response, err := ctrl.authService.Login(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login exitoso", response)
}

// GetMe obtiene el perfil del usuario autenticado
// GET /api/auth/me
func (ctrl *AuthController) GetMe(c *gin.Context) {
	// Obtener el user_id del contexto (lo pone el middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "No autenticado", nil)
		return
	}

	// Obtener el usuario
	user, err := ctrl.authService.GetUserByID(userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Usuario no encontrado", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Perfil obtenido", user)
}
