package services

import (
	"Reservify/dto"
	"Reservify/repositories"
	"Reservify/utils"
	"errors"
)

type UserService struct {
	userRepo *repositories.UserRepository
	authRepo *repositories.AuthRepository
}

func NewUserService(userRepo *repositories.UserRepository, authRepo *repositories.AuthRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
		authRepo: authRepo,
	}
}

// Obtiene todos los usuarios con paginación
func (s *UserService) GetAllUsers(params utils.PaginationParams) ([]dto.UserResponse, int64, error) {
	users, total, err := s.userRepo.FindAll(params)
	if err != nil {
		return nil, 0, err
	}

	// Convertir a DTOs
	var usersResponse []dto.UserResponse
	for _, user := range users {
		usersResponse = append(usersResponse, dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FullName:  user.FullName,
			Phone:     user.Phone,
			Role:      string(user.Role),
			CreatedAt: user.CreatedAt,
		})
	}

	return usersResponse, total, nil
}

// Obtiene un usuario por ID
func (s *UserService) GetUserByID(id uint) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	response := &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		Phone:     user.Phone,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
	}

	return response, nil
}

// Actualiza la información de un usuario
func (s *UserService) UpdateUser(id uint, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	// Buscar el usuario
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Actualizar campos
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	user.Phone = req.Phone

	// Guardar cambios
	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("error al actualizar usuario")
	}

	response := &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		Phone:     user.Phone,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
	}

	return response, nil
}

// Elimina un usuario (soft delete)
func (s *UserService) DeleteUser(id uint) error {
	return s.userRepo.Delete(id)
}

// Cambia la contraseña de un usuario
func (s *UserService) ChangePassword(userID uint, req *dto.ChangePasswordRequest) error {
	// Buscar el usuario
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	// Verificar la contraseña anterior
	if !utils.CheckPassword(req.OldPassword, user.PasswordHash) {
		return errors.New("contraseña anterior incorrecta")
	}

	// Hash de la nueva contraseña
	newPasswordHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return errors.New("error al procesar la nueva contraseña")
	}

	// Actualizar la contraseña
	if err := s.userRepo.UpdatePassword(userID, newPasswordHash); err != nil {
		return errors.New("error al actualizar la contraseña")
	}

	return nil
}

// Obtiene estadísticas de usuarios (para admin)
func (s *UserService) GetUserStats() (map[string]interface{}, error) {
	total, err := s.userRepo.Count()
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_users": total,
	}

	return stats, nil
}
