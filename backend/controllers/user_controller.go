package controllers

import (
	"Reservify/dto"
	"Reservify/services"
	"Reservify/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{userService: userService}
}

// Obtiene todos los usuarios (solo admin)
// GET /api/users?page=1&page_size=10&search=nombre
func (ctrl *UserController) GetAllUsers(c *gin.Context) {
	// Obtener parámetros de paginación
	params := utils.GetPaginationParams(c)

	// Obtener usuarios
	users, total, err := ctrl.userService.GetAllUsers(params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener usuarios", err)
		return
	}

	utils.PaginatedSuccessResponse(c, http.StatusOK, "Usuarios obtenidos exitosamente", users, total, params)
}

// Obtiene un usuario por ID
// GET /api/users/:id
func (ctrl *UserController) GetUserByID(c *gin.Context) {
	// Obtener ID del parámetro
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", err)
		return
	}

	// Obtener usuario
	user, err := ctrl.userService.GetUserByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Usuario obtenido exitosamente", user)
}

// Actualiza un usuario
// PUT /api/users/:id
func (ctrl *UserController) UpdateUser(c *gin.Context) {
	// Obtener ID del parámetro
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", err)
		return
	}

	// Obtener user_id del contexto (del token)
	currentUserID, _ := c.Get("user_id")
	currentUserRole, _ := c.Get("user_role")

	// Verificar permisos: solo puede editar su propio perfil o ser admin
	if currentUserID.(uint) != uint(id) && currentUserRole.(string) != "admin" {
		utils.ErrorResponse(c, http.StatusForbidden, "No tienes permisos para editar este usuario", nil)
		return
	}

	// Validar request
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos", err)
		return
	}

	// Actualizar usuario
	user, err := ctrl.userService.UpdateUser(uint(id), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Usuario actualizado exitosamente", user)
}

// Elimina un usuario (soft delete)
// DELETE /api/users/:id
func (ctrl *UserController) DeleteUser(c *gin.Context) {
	// Obtener ID del parámetro
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", err)
		return
	}

	// Verificar que no se elimine a sí mismo
	currentUserID, _ := c.Get("user_id")
	if currentUserID.(uint) == uint(id) {
		utils.ErrorResponse(c, http.StatusBadRequest, "No puedes eliminarte a ti mismo", nil)
		return
	}

	// Eliminar usuario
	if err := ctrl.userService.DeleteUser(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Usuario eliminado exitosamente", nil)
}

// Cambia la contraseña del usuario autenticado
// PUT /api/users/me/password
func (ctrl *UserController) ChangePassword(c *gin.Context) {
	// Obtener user_id del contexto
	userID, _ := c.Get("user_id")

	// Validar request
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos", err)
		return
	}

	// Cambiar contraseña
	if err := ctrl.userService.ChangePassword(userID.(uint), &req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Contraseña cambiada exitosamente", nil)
}

// Obtiene estadísticas de usuarios (solo admin)
// GET /api/users/stats
func (ctrl *UserController) GetUserStats(c *gin.Context) {
	stats, err := ctrl.userService.GetUserStats()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener estadísticas", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Estadísticas obtenidas exitosamente", stats)
}
